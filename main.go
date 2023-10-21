package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var processedLinks = make(map[string]bool)

// Create a channel to send data
var dataCh = make(chan string)

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
	r.HandleFunc("/site/{websiteURL}", presentPages(ctx)).Methods("GET", "POST")

	r.Handle("/", presentHome()).Methods("GET", "POST")

	r.HandleFunc("/progress", eventsHandler)
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

			insertErr := db.InsertPage(Page{Title: "", Content: "", URL: websiteUrl, IsSeedUrl: true})
			if insertErr != nil {
				log.Printf("Error saving page %s: %v", websiteUrl, insertErr)
				http.Error(w, "Error saving page", http.StatusMethodNotAllowed)
				return
			}
			log.Print("INSERT DATES MethodPost")

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
		websiteURLRaw := vars["websiteURL"]
		websiteURL, err := url.QueryUnescape(websiteURLRaw)
		if err != nil {
			http.Error(w, "Failed to parse websiteURL", http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodPost {
			log.Print("presentPages - POST")
			db, err := NewDB()
			if err != nil {
				http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
				return
			}
			log.Print("presentPages - POST 2")
			processLink(ctx, websiteURL, *db)
			log.Print("presentPages - POST 3")
			log.Printf("saved page")

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
			pageSize = 10
		}

		/*	db, err := NewDB()
			if err != nil {
				panic(err)
			}*/

		log.Printf("ListPages %s %d %d", websiteURL, page, pageSize)

		pagesList, pageErr := db.ListPages(websiteURL, page, pageSize)
		if pageErr != nil {
			panic(pageErr)
		}
		pagesComp := pages(pagesList, websiteURL)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}

func sendMessage(code string, messageStr string) {
	data := fmt.Sprintf(`{"code": "%s", "messageStr": "%s"}`, code, messageStr)
	dataCh <- data
}

func processLink(ctx context.Context, link string, db DB) {
	log.Print("StartProcessingLink", link)

	existingLink, err := db.GetExistingPage(link)
	if err != nil {
		log.Printf("Error GetLink from %s: %v", link, err)
		return
	}

	if existingLink {
		log.Printf("Already processed %d link.", len(processedLinks))
		return
	}

	log.Print("StartProcessingLink 1", link)

	links, err := fetchLinksFromPage(ctx, link)
	if err != nil {
		log.Printf("Error fetching links from %s: %v", link, err)
		return
	}

	title, content, err := fetchContentFromPage(ctx, link)
	if err != nil {
		log.Printf("Error fetching content from %s: %v", link, err)
		return
	}

	fmt.Println("Content from:", link)
	fmt.Println(links)

	insertErr := db.InsertPage(Page{Title: title, Content: content, URL: link})
	if insertErr != nil {
		log.Printf("Error saving page %s: %v", link, err)
		return
	}

	message := fmt.Sprintf(`processed %s`, link)
	fmt.Println("FinishedProcessingLink", message)
	// data := fmt.Sprintf(`{"code": "%s", "message": "%s"}`, "ProcessedLink", message)
	processedLinks[link] = true

	// Process child links
	for _, l := range links {

		insertErr := db.InsertLink(Link{sourceURL: link, URL: l, DateCreated: time.Now(), LastProcessed: time.Now()})
		if insertErr != nil {
			log.Printf("Error saving link %s: %v", link, err)
			return
		}
		processLink(ctx, l, db)
	}
}

func fetchLinksFromPage(ctx context.Context, urlStr string) ([]string, error) {
	var links []string
	url := os.Getenv("WEBSITE_NAME")
	js := fmt.Sprintf(`Array.from(document.querySelectorAll("a[href^='%s']")).map(a => a.href)`, url)

	err := chromedp.Run(ctx,
		network.SetExtraHTTPHeaders(network.Headers{
			"User-Agent": "MyCustomUserAgent",
		}),
		chromedp.Navigate(urlStr),
		chromedp.Evaluate(js, &links),
	)

	return links, err
}

func fetchContentFromPage(ctx context.Context, urlStr string) (string, string, error) {
	var content string
	var title string

	err := chromedp.Run(ctx,
		network.SetExtraHTTPHeaders(network.Headers{
			"User-Agent": "MyCustomUserAgent",
		}),
		chromedp.Navigate(urlStr),
		chromedp.Evaluate(`document.title`, &title),
		chromedp.Evaluate(`document.body.innerText`, &content),
	)

	return title, content, err
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Listen for client disconnection to stop sending events
	ctx := r.Context()
	// Send data to the client
	go func() {
		for {
			select {
			case data := <-dataCh:
				fmt.Fprintf(w, "data: %s\n\n", data)
				w.(http.Flusher).Flush()
			case <-ctx.Done():
				return

			}
		}
	}()
}
