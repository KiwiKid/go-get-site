package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	db.Migrate()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	r := mux.NewRouter()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	r.HandleFunc("/site/{websiteId}", handlePages(ctx)).Methods("GET", "POST")

	r.Handle("/", presentHome()).Methods("GET", "POST")

	r.Handle("/search", presentChat()).Methods("GET", "POST")
	r.Handle("/search/{threadId}", presentChat()).Methods("GET", "POST")

	r.Handle("/progress", presentLinkCount())
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Catch-all route triggered")
		fmt.Fprint(w, "Catch-all route")
	})

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Printf("error listening: %v", err)
	}

	// startURL := "https://www.hearwell.co.nz/" // Replace with your domain
	// 	connStr := os.Getenv("DB_CONN_STR")

	//url := os.Getenv("WEBSITE_NAME")
	//
	//db, err := NewDB(connStr)
	//if err != nil {
	//	panic(err)
	//}
	// processLink(ctx, url, *db)
}

func presentChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("presentChat - NewDB")
		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		vars := mux.Vars(r)
		threadIdStr := vars["threadId"]
		log.Printf("presentChat - threadIdStr %v", threadIdStr)

		if threadIdStr == "" {
			log.Printf("presentChat - EMPTY threadIdStr %v", threadIdStr)

			chatThreads, pageErr := db.ListChatThreads()
			if err != nil {
				panic(pageErr)
			}
			rand.Seed(time.Now().UnixNano())
			randomValue := uint(rand.Uint32())
			randomValueStr := strconv.FormatUint(uint64(randomValue), 10)
			newThreadURL := "/search/" + randomValueStr

			websites, pageErr := db.ListWebsites()
			if err != nil {
				panic(pageErr)
			}

			chatComp := threads(chatThreads, newThreadURL, websites)

			templ.Handler(chatComp).ServeHTTP(w, r)
			return
		}

		threadId, err := stringToUint(threadIdStr)
		if threadIdStr == "" || err != nil {
			http.Error(w, "Failed to ThreadId", http.StatusInternalServerError)
			return
		}

		log.Printf("chat http methd: %v", r.Method)
		if r.Method == http.MethodPost {
			websiteIdStr := r.FormValue("websiteId")
			websiteId, err := stringToUint(websiteIdStr)
			if websiteIdStr == "" || err != nil {
				log.Printf("Failed based on websiteId %v", err)
				http.Error(w, "Failed based on websiteId", http.StatusInternalServerError)
				return
			}
			query := r.FormValue("query")
			insErr := db.InsertChat(Chat{ThreadId: threadId, Message: query, WebsiteId: websiteId})
			if insErr != nil {
				http.Error(w, "InsertWebsite is failed", http.StatusBadRequest)
				return
			}

			queryRes, queryErr := db.QueryWebsite(query, websiteId)
			if queryErr != nil {
				http.Error(w, "QueryWebsite queryErr\n\n"+queryErr.Error(), http.StatusBadRequest)
				return
			}

			queryComp := queryResult(queryRes, websiteIdStr, query)

			templ.Handler(queryComp).ServeHTTP(w, r)

			//			insAIErr := db.InsertChat(Chat{ThreadId: threadId, Message: queryRes[0].String(), WebsiteId: websiteId})

			/*if insAIErr != nil {
				http.Error(w, "insAIErr failed", http.StatusBadRequest)
				return
			}*/
		}

		chats, pageErr := db.ListChats(threadId)
		if err != nil {
			panic(pageErr)
		}

		websiteIdStr := strconv.Itoa(int(chats[0].WebsiteId))

		newChatUrl := chats[0].ChatURL()

		chatComp := chat(threadIdStr, websiteIdStr, newChatUrl, chats)

		templ.Handler(chatComp).ServeHTTP(w, r)
	}
}

func presentHome() http.HandlerFunc {
	log.Print("PresentHome")
	var newSiteId uint

	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("PresentHome - NewDB")
		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		switch r.Method {
		case http.MethodPost:
			{
				log.Print("PresentHome - NewDB MethodPost")
				// Parse the form data to retrieve 'websiteUrl'
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Failed to parse form", http.StatusBadRequest)
					return
				}

				websiteUrl := r.FormValue("websiteUrl")
				if websiteUrl == "" {
					http.Error(w, "websiteUrl is required", http.StatusBadRequest)
					return
				}
				customQueryParams := r.FormValue("customQueryParams")

				site, err := db.InsertWebsite(Website{BaseUrl: websiteUrl, CustomQueryParam: customQueryParams})
				if err != nil {
					http.Error(w, "InsertWebsite is failed", http.StatusBadRequest)
					return
				}

				newSiteId = site.ID

				insertPageErr := db.InsertPage(Page{URL: site.BaseUrl, WebsiteId: site.ID})
				if insertPageErr != nil {
					http.Error(w, "InsertPage is failed", http.StatusBadRequest)
					return
				}

				/*insertErr := db.InsertLink(Link{URL: websiteUrl, SourceURL: websiteUrl})
				if insertErr != nil {
					log.Printf("Error InsertLink 1 %s: %v", websiteUrl, insertErr)
					//http.Error(w, "Error saving page", http.StatusMethodNotAllowed)
					//		return
				}*/

				/*insertErr := db.InsertPage(Page{Title: "", Content: "", URL: websiteUrl, IsSeedUrl: true, WebsiteId: site.ID})
				 */
				log.Printf("INSERT DATES MethodPost %v", site)

			}
		case http.MethodDelete:
			{
				websiteIdStr := r.FormValue("websiteId")
				websiteId, err := stringToUint(websiteIdStr)
				if err != nil {
					http.Error(w, "stringToUint is failed", http.StatusBadRequest)
					return
				}
				insertPageErr := db.DeleteWebsite(websiteId)
				if insertPageErr != nil {
					http.Error(w, "InsertPage is failed", http.StatusBadRequest)
					return
				}
			}
		}

		websites, pageErr := db.ListWebsites()
		if err != nil {
			panic(pageErr)
		}
		homeComp := home(websites, newSiteId)

		templ.Handler(homeComp).ServeHTTP(w, r)
	}
}

/*
func handleProcessing(ctx context.Context) http.HandlerFunc {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]

		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			http.Error(w, "Failed to stringToUint", http.StatusInternalServerError)
			return
		}
		website, err := db.GetWebsite(websiteId)
		if err != nil {
			http.Error(w, "Failed to GetWebsite", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodPost {

			processErr := processWebsite(ctx, *db, *website)
			if processErr != nil {
				http.Error(w, "Failed to processWebsite", http.StatusInternalServerError)
				return
			}

		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		pageSizeStr := r.URL.Query().Get("pageSize")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 100
		}

			db, err := NewDB()
			if err != nil {
				panic(err)
			}

		//log.Printf("ListPages %s %d %d", websiteURL, page, pageSize)

		pagesList, pageErr := db.ListPages(website.ID, page, pageSize)
		if pageErr != nil {
			panic(pageErr)
		}

		count, countErr := db.CountLinksAndPages(website.ID)
		if countErr != nil {
			log.Printf("Error CountLinksAndPages link %s: %v", websiteId, err)
			return
		}

		pagesComp := pages(pagesList, *website, *count)

		templ.Handler(pagesComp).ServeHTTP(w, r)

	}
}*/

func handlePages(ctx context.Context) http.HandlerFunc {
	log.Print("handlePages")
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]

		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			http.Error(w, "Failed to stringToUint", http.StatusInternalServerError)
			return
		}
		website, err := db.GetWebsite(websiteId)
		if err != nil {
			http.Error(w, "Failed to GetWebsite", http.StatusInternalServerError)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		pageSizeStr := r.URL.Query().Get("pageSize")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 100
		}

		if r.Method == http.MethodPost {
			log.Print("handlePages - POST")
			db, err := NewDB()
			if err != nil {
				http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
				return
			}

			//websiteURL := vars["websiteUrl"]
			/*log.Printf("Link: %v %v", website, websiteId)
			insertErr := db.InsertLink(Link{URL: website.BaseUrl, DateCreated: time.Now(), LastProcessed: time.Now(), WebsiteId: websiteId})
			if insertErr != nil {
				log.Printf("Error saving base link %s: %v", website.BaseUrl, insertErr)
				// return
			}
			unProcess := db.unProcessLink(website.BaseUrl)
			if unProcess != nil {
				log.Printf("Error unProcessLink link %s: %v", website.BaseUrl, unProcess)
				// return
			}*/
			processAll := r.FormValue("processAll") == "on"

			log.Printf("handlePages - processWebsite processAll:%v", processAll)
			processErr := processWebsite(ctx, *db, *website, processAll)
			if processErr != nil {
				http.Error(w, "Failed to processWebsite", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated) // 201 Created status
		}

		/*	db, err := NewDB()
			if err != nil {
				panic(err)
			}*/

		//log.Printf("ListPages %s %d %d", websiteURL, page, pageSize)

		pagesList, pageErr := db.ListPages(website.ID, page, pageSize)
		if pageErr != nil {
			panic(pageErr)
		}

		count, countErr := db.CountLinksAndPages(website.ID)
		if countErr != nil {
			log.Printf("Error CountLinksAndPages link %d: %v", websiteId, err)
			return
		}

		thisPageUrl := fmt.Sprintf("/site/%d?page=%d&pageSize=%d", website.ID, page, pageSize)
		prevPageUrl := fmt.Sprintf("/site/%d?page=%d&pageSize=%d", website.ID, page, pageSize)
		nextPageUrl := fmt.Sprintf("/site/%d?page=%d&pageSize=%d", website.ID, page, pageSize)

		pagesComp := pages(pagesList, *website, *count, thisPageUrl, prevPageUrl, nextPageUrl)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}

func processWebsite(ctx context.Context, db DB, website Website, processAll bool) error {
	log.Print("StartProcessingSite ", website.BaseUrl)
	var pageUpdatedAfter time.Time
	if processAll {
		pageUpdatedAfter = time.Now().Add(-365 * 24 * time.Hour)
	} else {
		pageUpdatedAfter = time.Now().Add(-7 * 24 * time.Hour)
	}

	pagesToProcess, err := db.GetPages(website.ID, 1, 5, processAll, pageUpdatedAfter)
	log.Printf("linksToProcess: %d", len(pagesToProcess))
	if err != nil {
		log.Printf("Error GetLink from %v", err)
		return err
	}

	if len(pagesToProcess) > 0 {
		pagesToSave, err := fetchContentFromPages(ctx, website, pagesToProcess, 5)
		log.Printf("Got %d pagesToSave from fetchContentFromPages %v", len(pagesToSave), err)

		for _, page := range pagesToSave {
			insertErr := db.InsertPage(page)
			if insertErr != nil {
				updateErr := db.UpdatePage(page)
				if updateErr != nil {
					log.Printf("Error saving page 2 %v: %v", page.ID, updateErr)
					return updateErr
				}
			}
		}
	}

	return err
}

func SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		/*expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		success := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)

		if success.Error != nil {
			return fmt.Errorf("could not set cookie %s", name)
		} else {
			log.Print("Set cookie")
		}*/
		return nil
	})
}

func fetchContentFromPages(ctx context.Context, website Website, pages []Page, remainingToProcess int) ([]Page, error) {
	var content string
	var title string
	var links []string
	var newPages []Page

	js := fmt.Sprintf(`Array.from(document.querySelectorAll("a[href^='%s']")).map(a => a.href)`, website.BaseUrl)

	/*autoUrl, addQueryErr := addQueryParam(urlStr, "1c8ca3a202b84c47961b79700b40f01a")
	if addQueryErr != nil {
		panic(addQueryErr)
	}*/

	urlWithParams := fmt.Sprintf("%s?%s", website.StartUrl, website.CustomQueryParam)
	log.Printf("url to go: %s", urlWithParams)

	tasks := []chromedp.Action{
		chromedp.Navigate(urlWithParams),
	}

	if website.LoginName != "" {
		tasks = append(tasks,
			chromedp.WaitVisible(`#txtUserName`, chromedp.ByID),
			chromedp.SendKeys(`#txtUserName`, website.LoginName),
			chromedp.SendKeys(`#txtPassword`, website.LoginPass),
			chromedp.Click(`#btnLogin`),
			chromedp.WaitVisible(`#elementAfterLogin`, chromedp.ByID),
		)
	}

	for _, page := range pages {
		// Add the rest of the tasks
		tasks = append(tasks,
			chromedp.Navigate(page.URL),
			chromedp.Evaluate(`document.title`, &title),
			chromedp.Evaluate(`document.body.innerText`, &content),
			chromedp.Evaluate(js, &links),
		)

		// BUILD A "Page" object
		newPage := Page{
			URL:     page.URL,
			Title:   title,
			Content: content,
			Links:   links,
		}
		// ADD the Page object to the "pages" list
		newPages = append(newPages, newPage)
	}

	log.Print("Starting tasks")
	err := chromedp.Run(ctx, tasks...)
	log.Print("Finished tasks")
	return newPages, err
}

func presentLinkCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]
		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			http.Error(w, "Failed to stringToUint", http.StatusInternalServerError)
			return
		}

		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		count, countErr := db.CountLinksAndPages(websiteId)
		if countErr != nil {
			log.Printf("Error CountLinksAndPages link %d: %v", websiteId, err)
			return
		}

		pagesComp := process(*count, websiteIdStr)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}
