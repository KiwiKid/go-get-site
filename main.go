package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//// Uncomment to deploy DB changes (commentted out as it improves rebuild time)
	//// (Or comment to improve build speed)
	//db, err := NewDB()
	//if err != nil {
	//	panic(err)
	//}
	//db.Migrate()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	r := mux.NewRouter()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	r.Handle("/", presentWebsite()).Methods("GET", "POST")
	r.HandleFunc("/site/{websiteId}", presentWebsite()).Methods("GET", "POST", "PUT")
	r.HandleFunc("/site/{websiteId}/{deleteOpt}", presentWebsite()).Methods("DELETE")

	r.HandleFunc("/site/{websiteId}/login", presentLogin(ctx)).Methods("GET")

	r.HandleFunc("/site/{websiteId}/pages", handlePages(ctx)).Methods("GET", "POST")

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

func presentWebsite() http.HandlerFunc {
	log.Print("presentWebsite")

	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("presentWebsite - NewDB")
		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		switch r.Method {
		case http.MethodPost:
			{
				log.Print("presentWebsite - NewDB MethodPost")
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
				loginName := r.FormValue("loginName")
				loginNameSelector := r.FormValue("loginNameSelector")
				loginPass := r.FormValue("loginPass")
				loginPassSelector := r.FormValue("loginPassSelector")
				site, err := db.InsertWebsite(Website{BaseUrl: websiteUrl, CustomQueryParam: customQueryParams, LoginName: loginName, LoginPass: loginPass, LoginNameSelector: loginNameSelector, LoginPassSelector: loginPassSelector})
				if err != nil {
					http.Error(w, "InsertWebsite is failed", http.StatusBadRequest)
					return
				}
				emptyLink := []string{}
				inserPageErr := db.InsertPage(Page{URL: site.BaseUrl, WebsiteId: site.ID, Links: emptyLink})
				if inserPageErr != nil {
					log.Fatalf("inserPageErr  %v", inserPageErr)
					http.Error(w, "UpsertPage is failed", http.StatusBadRequest)
					panic(inserPageErr)
				}

				log.Printf("INSERT DATES MethodPost %v", site)

			}
		case http.MethodPut:
			{
				log.Print("presentWebsite - NewDB MethodPut")

				// Parse the form data to retrieve the fields
				err := r.ParseForm()
				if err != nil {
					http.Error(w, "Failed to parse form", http.StatusBadRequest)
					return
				}

				websiteIdStr := r.FormValue("websiteId")
				websiteId, err := stringToUint(websiteIdStr)
				if err != nil {
					http.Error(w, "stringToUint is failed", http.StatusBadRequest)
					return
				}

				websiteUrl := r.FormValue("websiteUrl")
				if websiteUrl == "" {
					http.Error(w, "websiteUrl is required", http.StatusBadRequest)
					return
				}
				customQueryParams := r.FormValue("customQueryParams")
				loginName := r.FormValue("loginName")
				loginNameSelector := r.FormValue("loginNameSelector")
				loginPass := r.FormValue("loginPass")
				loginPassSelector := r.FormValue("loginPassSelector")
				submitButtonSelector := r.FormValue("submitButtonSelector")
				successIndicatorSelector := r.FormValue("successIndicatorSelector")

				site, err := db.UpdateWebsite(Website{
					ID:                       websiteId,
					BaseUrl:                  websiteUrl,
					CustomQueryParam:         customQueryParams,
					LoginName:                loginName,
					LoginPass:                loginPass,
					LoginNameSelector:        loginNameSelector,
					LoginPassSelector:        loginPassSelector,
					SubmitButtonSelector:     submitButtonSelector,
					SuccessIndicatorSelector: successIndicatorSelector,
				})
				if err != nil {
					http.Error(w, "UpdateWebsite is failed", http.StatusBadRequest)
					return
				}

				log.Printf("UPDATE DATES MethodPut %v", site)
			}
		case http.MethodDelete:
			{

				vars := mux.Vars(r)
				websiteIdStr := vars["websiteId"]
				websiteId, err := stringToUint(websiteIdStr)
				if err != nil {
					http.Error(w, "stringToUint for websiteId is failed", http.StatusBadRequest)
					return
				}

				deleteOpt := vars["deleteOpt"]

				log.Printf("deleteOpt %s", deleteOpt)

				switch deleteOpt {
				case "all":
					{
						deleteWebsite := db.DeleteWebsite(websiteId)
						if deleteWebsite != nil {
							http.Error(w, "InsertPage is failed", http.StatusBadRequest)
							return
						}

						r.Header.Add("HX-Redirect", "/")
					}
				case "pages":
					{
						deleteWebsite := db.DeletePages(websiteId)
						if deleteWebsite != nil {
							http.Error(w, "InsertPage is failed", http.StatusBadRequest)
							return
						}

						website, err := db.GetWebsite(websiteId)
						if err != nil {
							http.Error(w, "GetWebsite failed", http.StatusBadRequest)
							return
						}
						emptyLink := []string{}
						newEmptyPage := Page{
							URL:       website.BaseUrl,
							WebsiteId: website.ID,
							Links:     emptyLink,
						}
						insertErr := db.InsertPage(newEmptyPage)
						if insertErr != nil {
							log.Printf("Error adding seed page after pages delete %v", insertErr)
							http.Error(w, "Error adding seed page", http.StatusBadRequest)

							return
						}
					}
				default:
					{
						panic(fmt.Sprintf("Invalid delete option: '%s'", deleteOpt))
					}
				}

			}
		}

		websites, pageErr := db.ListWebsites()
		if err != nil {
			panic(pageErr)
		}
		homeComp := home(websites)

		templ.Handler(homeComp).ServeHTTP(w, r)
	}
}

func GetPageDoneCacheKey(websiteId uint, url string) string {
	return fmt.Sprintf("%d:%s", websiteId, url)
}

func presentLogin(ctx context.Context) http.HandlerFunc {
	log.Printf("presentLogin")

	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]
		websiteId, webIderr := stringToUint(websiteIdStr)
		if webIderr != nil {
			http.Error(w, fmt.Sprintf("Failed to stringToUint websiteId %s", websiteIdStr), http.StatusInternalServerError)
			panic(webIderr)
		}

		website, getWebErr := db.GetWebsite(websiteId)
		if webIderr != nil {
			http.Error(w, fmt.Sprintf("Failed to get getWebErr %v", getWebErr), http.StatusInternalServerError)
			loginRes := loginResult(*website, "", "", "Failed to get getWebErr", webIderr.Error())

			templ.Handler(loginRes).ServeHTTP(w, r)
		}
		log.Printf("presentLogin:getLoginTasks")

		tasks := getLoginTasks(*website)

		var content string
		var title string
		tasks = append(tasks,
			chromedp.Evaluate(`document.title`, &title),
			chromedp.Evaluate(`document.body.innerText`, &content),
		)
		log.Printf("presentLogin:Run")

		chromeRunErr := chromedp.Run(ctx, tasks...)
		if chromeRunErr != nil {
			log.Printf("Error logging in for websiteIdStr %s: %v", websiteIdStr, chromeRunErr)
			loginRes := loginResult(*website, title, content, fmt.Sprintf("Error logging in for websiteIdStr %s", websiteIdStr), chromeRunErr.Error())

			templ.Handler(loginRes).ServeHTTP(w, r)

			//http.Error(w, fmt.Sprintf("Failed to get getWebErr %v", chromeRunErr), http.StatusInternalServerError)
			// return newPages, err
		}
		log.Printf("presentLogin:Run:Complete")

		loginRes := loginResult(*website, title, content, "", "")

		templ.Handler(loginRes).ServeHTTP(w, r)

	}
}

func handlePages(ctx context.Context) http.HandlerFunc {
	log.Print("handlePages")
	addedPagesSet := make(map[string]struct{})
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]

		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to stringToUint websiteId %s", websiteIdStr), http.StatusInternalServerError)
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

			processAll := r.FormValue("processAll") == "on"
			processPageSizeStr := r.FormValue("processPageSize")
			pageSize, err := strconv.Atoi(processPageSizeStr)
			if err != nil || pageSize <= 0 {
				pageSize = 5
			}

			log.Printf("handlePages - processWebsite processAll:%v", processAll)
			processErr := processWebsite(ctx, *db, *website, processAll, 1, pageSize, addedPagesSet)
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

		prog := Progress{
			Total: count.TotalLinks,
			Done:  count.LinksHavePages,
		}

		thisPageUrl := fmt.Sprintf("/site/%d/pages?page=%d&pageSize=%d", website.ID, page, pageSize)
		prevPageUrl := fmt.Sprintf("/site/%d/pages?page=%d&pageSize=%d", website.ID, page, pageSize)

		var nextPageUrl string
		if count.TotalLinks > (page * pageSize) {
			nextPageUrl = fmt.Sprintf("/site/%d/pages?page=%d&pageSize=%d", website.ID, page, pageSize)
		} else {
			nextPageUrl = ""
		}

		percentage := fmt.Sprintf("%f", (float64(prog.Done) / float64(prog.Total) * 100.0))
		log.Printf("percentage percentage percentage %s", percentage)

		pagesComp := pages(pagesList, *website, *count, thisPageUrl, prevPageUrl, nextPageUrl, addedPagesSet, percentage)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}

func processWebsite(ctx context.Context, db DB, website Website, processAll bool, page int, pageSize int, addedPagesSet map[string]struct{}) error {
	log.Print("StartProcessingSite ", website.BaseUrl)
	var pageProcessedAfter time.Time
	if processAll {
		pageProcessedAfter = time.Now().Add(-365 * 24 * time.Hour)
	} else {
		pageProcessedAfter = time.Now().Add(-7 * 24 * time.Hour)
	}

	pagesToProcess, err := db.GetPages(website.ID, page, pageSize, processAll, pageProcessedAfter)
	linksAlreadyProcessed, apErr := db.GetCompletedPageUrls(website.ID)
	for _, url := range linksAlreadyProcessed {
		addedPagesSet[GetPageDoneCacheKey(website.ID, url)] = struct{}{}
	}
	log.Printf("GetPages got %d links to process [processAll:%v] [pageProcessedAfter:%v]", len(pagesToProcess), processAll, pageProcessedAfter)
	if err != nil || apErr != nil {
		log.Printf("Error GetLink from %v", err)
		return err
	}

	if len(pagesToProcess) > 0 {
		pagesToSave, err := fetchContentFromPages(ctx, website, pagesToProcess, 5, addedPagesSet)
		if err != nil {
			panic(err)
		}
		log.Printf("Got %d pagesToSave from fetchContentFromPages", len(pagesToSave))

		for _, page := range pagesToSave {
			//if _, exists := addedPagesSet[page.URL]; !exists {
			insertErr := db.UpsertPage(page)
			if insertErr != nil {
				return insertErr
			} // else {
			//	addedPagesSet[page.URL] = struct{}{}
			//}
			//}
		}
	}

	return err
}

func SetCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		err := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)

		if err != nil {
			return fmt.Errorf("could not set cookie %s", name)
		} else {
			log.Print("Set cookie")
		}
		return nil
	})
}

func logAction(message string, id string, querySelector bool) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if querySelector {
			log.Printf("%s document.querySelectorAll(\"%s\")", message, id)
		} else {
			log.Printf("%s %s", message, id)
		}

		return nil
	})
}

func getLoginTasks(website Website) []chromedp.Action {
	tasks := []chromedp.Action{}

	url := website.websiteNavigateURL()

	if url != "" {

		tasks = append(tasks,
			logAction("Navigate-to:start", url, false),
			chromedp.Navigate(url),
			logAction("Navigate-to-success:", url, false),
		)
	} else {
		log.Print("fetchContentFromPagesERROR-  no urlWithParams")
	}

	if website.LoginName != "" {

		log.Printf("fetchContentFromPages Logging-in as '%s'", website.LoginName)

		tasks = append(tasks,
			logAction("LoginStart-looking-for:", website.LoginNameSelector, true),
			chromedp.WaitVisible(website.LoginNameSelector),
			logAction("LoginStart-LoginNameSelector-SUCCESS", website.LoginNameSelector, true),
			chromedp.SendKeys(website.LoginNameSelector, website.LoginName),
			logAction("LoginStart-LoginNameSelector-Entry-SUCCESS:", website.LoginNameSelector, true),
			chromedp.SendKeys(website.LoginPassSelector, website.LoginPass),
			logAction("LoginStart-LoginPassSelector-Entry-SUCCESS:", website.LoginPassSelector, true),
			logAction("LoginStart-looking-for SubmitButtonSelector: ", website.SubmitButtonSelector, true),
			chromedp.WaitVisible(website.SubmitButtonSelector),
			logAction("LoginStart-SubmitButtonSelector-found: ", website.SubmitButtonSelector, true),
			chromedp.Click(website.SubmitButtonSelector),
			logAction("LoginStart:SubmitButtonSelector-clicked. Correct Login & Pass? Looking for SuccessIndicatorSelector: ", website.SuccessIndicatorSelector, true),
			chromedp.WaitVisible(website.SuccessIndicatorSelector),
			logAction("LoginStart:SUCCESS", "", false),
		)
	}

	return tasks

}

func fetchContentFromPages(ctx context.Context, website Website, pages []Page, remainingToProcess int, addedPagesSet map[string]struct{}) ([]Page, error) {
	log.Print("fetchContentFromPages:start")

	var newPages []Page

	allLinksJS := `Array.from(document.querySelectorAll("*[href]")).map((i) => i.href)`

	/*autoUrl, addQueryErr := addQueryParam(urlStr, "1c8ca3a202b84c47961b79700b40f01a")
	if addQueryErr != nil {
		panic(addQueryErr)
	}*/

	tasks := getLoginTasks(website)

	if website.RequestCookieName != "" && website.RequestCookieValue != "" {
		log.Printf("fetchContentFromPages Setting Cookies")

		tasks = append(tasks, SetCookie(website.RequestCookieName, website.RequestCookieValue, website.BaseUrl, "/", false, false))
	}

	for _, page := range pages {

		var title string
		var content string
		var links []string

		log.Printf("fetchContentFromPages Starting Page %s", page.URL)

		// Add the rest of the tasks
		tasks = append(tasks,
			logAction("fetchContentFromPages:Navigate-to:", page.URL, false),
			chromedp.Navigate(page.URL),
			logAction("fetchContentFromPages:Navigate-to:Success", page.URL, false),
			chromedp.Evaluate(`document.title`, &title),
			chromedp.Evaluate(`document.body.innerText`, &content),
			logAction(fmt.Sprintf("fetchContentFromPages:Got Title (%d) and Content (%d)", len(title), len(content)), page.URL, false),
			chromedp.Evaluate(`document.body.innerText`, &content),
			logAction("fetchContentFromPages:GetLinks", page.URL, false),
			chromedp.Evaluate(allLinksJS, &links),
			logAction("fetchContentFromPages:GotLinks", page.URL, false),
		)

		err := chromedp.Run(ctx, tasks...)
		if err != nil {
			msg := fmt.Sprintf("Error fetching content for page %s: %s", page.URL, err.Error())

			db, dbErr := NewDB()
			if dbErr != nil {
				panic(dbErr)
			}

			setWarErr := db.UpdateWarning(page.ID, msg)
			if setWarErr != nil {
				panic(setWarErr)
			}
			panic(err)
			// return newPages, err
		}

		// BUILD A "Page" object
		newPage := Page{
			URL:       page.URL,
			Title:     title,
			Content:   content,
			Links:     links,
			WebsiteId: website.ID,
		}
		log.Printf("fetchContentFromPages - new page \n%v", newPage)
		// ADD the Page object to the "pages" list
		newPages = append(newPages, newPage)
		emptyLink := []string{}

		for _, link := range links {
			for _, baseUrl := range strings.Split(website.BaseUrl, ",") {
				link, err = stripAnchors(link)
				if err != nil {
					log.Printf("fetchContentFromPages non-page link %s", link)
					panic(err)
				}
				if linkCouldBePage(link, baseUrl) {

					if _, exists := addedPagesSet[GetPageDoneCacheKey(newPage.WebsiteId, link)]; !exists {

						log.Printf("fetchContentFromPages page link %s", link)

						newEmptyPage := Page{

							URL:       link,
							WebsiteId: website.ID,
							Links:     emptyLink,
						}
						newPages = append(newPages, newEmptyPage)
						addedPagesSet[GetPageDoneCacheKey(newPage.WebsiteId, link)] = struct{}{}
						// You might want to add this newPage to some slice or process it further
						break
					} else {
						log.Printf("fetchContentFromPages already added %s", link)
					}
				} else {
					log.Printf("fetchContentFromPages not a relvant link %s", link)
				}

			}
		}
		if len(newPages) > remainingToProcess {
			log.Printf("fetchContentFromPages Early exit %d paged processed", len(newPages))
			return newPages, nil
		}
	}

	log.Print("Finished tasks")
	return newPages, nil
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

		website, webErr := db.GetWebsite(websiteId)
		if webErr != nil {
			log.Printf("Error CountLinksAndPages link %d: %v", websiteId, webErr)
			return
		}

		pagesComp := process(*count, *website)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}
