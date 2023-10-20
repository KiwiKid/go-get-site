package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	//defer cancel()

	r := mux.NewRouter()

	//ctx, cancel = chromedp.NewContext(ctx)
	//defer cancel()

	http.Handle("/", presentHome())

	r.HandleFunc("/site/{websiteURL}", presentPages())

	http.HandleFunc("/site/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

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

		taskComp := task(websiteUrl)

		// Return the templ.Handler(taskComp)
		templ.Handler(taskComp).ServeHTTP(w, r)

	})

	http.HandleFunc("/progress", eventsHandler)

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
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
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := NewDB()
		if err != nil {
			panic(err)
		}
		websites, pageErr := db.ListWebsites()
		if err != nil {
			panic(pageErr)
		}
		homeComp := home(websites)

		templ.Handler(homeComp).ServeHTTP(w, r)
	}
}

func presentPages() http.HandlerFunc {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteURL := vars["websiteURL"]

		if r.Method == http.MethodPost {
			// Logic for POST request to add a website
			// This is where you would read the request body, parse any required data, and add the website to your database
			// For example:
			// 1. Parse the request body.
			// 2. Validate the data.
			// 3. Insert the website into the database.
			// 4. Send a response to the client.

			db, err := NewDB()
			if err != nil {
				http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
				return
			}

			insertErr := db.InsertPage(Page{Title: "", Content: "", URL: websiteURL})
			if insertErr != nil {
				log.Printf("Error saving page %s: %v", websiteURL, err)
				return
			}

			w.WriteHeader(http.StatusCreated) // 201 Created status
			fmt.Fprintf(w, "Website %s added successfully", websiteURL)
		}

		// The previous logic for GET request
		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1 // default or fallback value
		}

		pageSizeStr := r.URL.Query().Get("pageSize")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize <= 0 {
			pageSize = 10 // default or fallback value
		}

		fmt.Fprintf(w, "Website URL: %s", websiteURL)

		/*	db, err := NewDB()
			if err != nil {
				panic(err)
			}*/

		pagesList, pageErr := db.ListPages(websiteURL, page, pageSize)
		if err != nil {
			panic(pageErr)
		}
		pagesComp := pages(pagesList)

		templ.Handler(pagesComp).ServeHTTP(w, r)

	}
}

func sendMessage(code string, messageStr string) {
	data := fmt.Sprintf(`{"code": "%s", "messageStr": "%s"}`, code, messageStr)
	dataCh <- data
}

func processLink(ctx context.Context, link string, db DB) {
	if processedLinks[link] {
		//msg := fmt.Sprintf("Already processed %d link.", len(processedLinks))
		//sendMessage("Finished", msg)
		return
	}
	sendMessage("StartProcessingLink", link)

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
	fmt.Println(content)

	insertErr := db.InsertPage(Page{Title: title, Content: content, URL: link})
	if insertErr != nil {
		log.Printf("Error saving page %s: %v", link, err)
		return
	}

	// Mark the link as processed
	insertUrlErr := db.InsertLink(Link{URL: link, DateCreated: time.Now(), LastProcessed: time.Now()})
	if insertUrlErr != nil {
		log.Printf("Error saving link %s: %v", link, err)
		return
	}
	message := fmt.Sprintf(`processed %s`, link)
	sendMessage("FinishedProcessingLink", message)
	// data := fmt.Sprintf(`{"code": "%s", "message": "%s"}`, "ProcessedLink", message)
	processedLinks[link] = true

	// Process child links
	for _, l := range links {
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
