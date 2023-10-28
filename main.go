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
	r.HandleFunc("/site/{websiteId}", presentPages(ctx)).Methods("GET", "POST")

	r.Handle("/", presentHome()).Methods("GET", "POST")

	r.Handle("/chat", presentChat()).Methods("GET", "POST")
	r.Handle("/chat/{threadId}", presentChat()).Methods("GET", "POST")

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
			newThreadURL := "/chat/" + randomValueStr

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
			message := r.FormValue("message")
			insErr := db.InsertChat(Chat{ThreadId: threadId, Message: message, WebsiteId: websiteId})
			if insErr != nil {
				http.Error(w, "InsertWebsite is failed", http.StatusBadRequest)
				return
			}

			queryRes, queryErr := db.QueryWebsite(message)
			if queryErr != nil {
				http.Error(w, "Query failed", http.StatusBadRequest)
				return
			}

			insAIErr := db.InsertChat(Chat{ThreadId: threadId, Message: queryRes[0].String(), WebsiteId: websiteId})

			if insAIErr != nil {
				http.Error(w, "insAIErr failed", http.StatusBadRequest)
				return
			}
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

	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("PresentHome - NewDB")
		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		if r.Method == http.MethodPost {
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

			site, err := db.InsertWebsite(Website{BaseUrl: websiteUrl})
			if err != nil {
				http.Error(w, "InsertWebsite is failed", http.StatusBadRequest)
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

		websites, pageErr := db.ListWebsites()
		if err != nil {
			panic(pageErr)
		}
		homeComp := home(websites)

		templ.Handler(homeComp).ServeHTTP(w, r)
	}
}

func presentPages(ctx context.Context) http.HandlerFunc {
	log.Print("presentPages")
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
			log.Print("presentPages - POST")
			db, err := NewDB()
			if err != nil {
				http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
				return
			}

			//websiteURL := vars["websiteUrl"]
			log.Printf("Link: %v %v", website, websiteId)
			insertErr := db.InsertLink(Link{URL: website.BaseUrl, DateCreated: time.Now(), LastProcessed: time.Now(), WebsiteId: websiteId})
			if insertErr != nil {
				log.Printf("Error saving base link %s: %v", website.BaseUrl, insertErr)
				// return
			}
			unProcess := db.unProcessLink(website.BaseUrl)
			if unProcess != nil {
				log.Printf("Error unProcessLink link %s: %v", website.BaseUrl, unProcess)
				// return
			}

			log.Print("presentPages - InsertLink")
			processLink(ctx, website.BaseUrl, *db, *website)

			w.WriteHeader(http.StatusCreated) // 201 Created status
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
			log.Printf("Error CountLinksAndPages link %s: %v", websiteId, err)
			return
		}

		pagesComp := pages(pagesList, *website, *count)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}

func processLink(ctx context.Context, baseLink string, db DB, website Website) {
	log.Print("StartProcessingLink", baseLink)

	linksToProcess, err := db.GetUnprocessedLinks(baseLink, 20)
	log.Printf("linksToProcess: %d", len(linksToProcess))
	if err != nil {
		log.Printf("Error GetLink from %s: %v", baseLink, err)
		return
	}

	if len(linksToProcess) == 0 {
		log.Printf("SetLinkProcessed() %s", baseLink)
		linkDoneErr := db.SetLinkProcessed(baseLink)
		if linkDoneErr != nil {
			log.Printf("Error marking link done %s: %v", baseLink, err)
			return
		}
		return
	}

	for _, link := range linksToProcess {
		log.Printf("linksToProcess IN LOOP ===> : %s", link.URL)

		/*links, err := fetchLinksFromPage(ctx, link.URL)
		if err != nil {
			log.Printf("Error fetching links from %s: %v", link, err)
			return
		}*/

		// Process child links

		title, content, links, err := fetchContentFromPage(ctx, link.URL)
		if err != nil {
			log.Printf("Error fetching content from %v: %v", link, err)
			return
		}
		log.Printf("UpdatePage %v - %v", title, content)
		insertErr := db.InsertPage(Page{Title: title, Content: content, URL: link.URL, WebsiteId: website.ID})
		if insertErr != nil {
			updateErr := db.UpdatePage(Page{Title: title, Content: content, URL: link.URL, WebsiteId: website.ID})
			if updateErr != nil {
				log.Printf("Error saving page 2 %s: %v", link, updateErr)
				return
			}
		}
		log.Printf("links start %v", links)
		for _, l := range links {
			log.Printf("links - %v", l)
			insertErr := db.InsertLink(Link{SourceURL: baseLink, URL: l, DateCreated: time.Now(), WebsiteId: website.ID})
			if insertErr != nil {
				log.Printf("Error saving link (maybe duplicate). Continuing %v: %v", link, insertErr)
			} else {
				log.Printf("saved link %v", l)
			}
		}

		message := fmt.Sprintf(`processed %v`, link)
		fmt.Println("FinishedProcessingLink", message)
	}
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

func fetchContentFromPage(ctx context.Context, urlStr string) (string, string, []string, error) {
	var content string
	var title string
	var links []string

	js := fmt.Sprintf(`Array.from(document.querySelectorAll("a[href^='%s']")).map(a => a.href)`, urlStr)

	/*autoUrl, addQueryErr := addQueryParam(urlStr, "1c8ca3a202b84c47961b79700b40f01a")
	if addQueryErr != nil {
		panic(addQueryErr)
	}*/

	log.Printf("url to go: %s", urlStr)

	tasks := []chromedp.Action{
		chromedp.Navigate(urlStr),
	}
	/*
		pass := os.Getenv("PASS")
		// Check if the environment variable "NAME" exists
		if name := os.Getenv("NAME"); name != "" {
			// Add the login tasks if the NAME env variable exists
			tasks = append(tasks,
				chromedp.WaitVisible(`#txtUserName`, chromedp.ByID), // Wait for an element that's visible after login
				chromedp.SendKeys(`#txtUserName`, name),
				chromedp.SendKeys(`#txtPassword`, pass),
				chromedp.Click(`#btnLogin`),
				chromedp.WaitVisible(`#elementAfterLogin`, chromedp.ByID), // Wait for an element that's visible after login
			)
		}
	*/
	// Add the rest of the tasks
	tasks = append(tasks,
		chromedp.Evaluate(`document.title`, &title),
		chromedp.Evaluate(`document.body.innerText`, &content),
		chromedp.Evaluate(js, &links),
	)

	log.Print("Starting tasks")
	err := chromedp.Run(ctx, tasks...)
	log.Print("Finished tasks")
	return title, content, links, err
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
