package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

var processedLinks = make(map[string]bool)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	
	connStr := os.Getenv("DB_CONN_STR")

	url := os.Getenv("WEBSITE_NAME")

	db, err := NewDB(connStr)
	if err != nil {
		panic(err)
	}

	// startURL := "https://www.hearwell.co.nz/" // Replace with your domain
	processLink(ctx, url, *db)
}

func processLink(ctx context.Context, link string, db DB) {
	if processedLinks[link] {
		return
	}
	links, err := fetchLinksFromPage(ctx, link)
	if err != nil {
		log.Printf("Error fetching links from %s: %v", link, err)
		return
	}

	title, content , err := fetchContentFromPage(ctx, link)
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
	insertUrlErr := db.InsertLink(Link{URL: link, DateCreated: time.Now(), LastProcessed: time.Now() })
	if insertUrlErr != nil {
		log.Printf("Error saving link %s: %v", link, err)
		return
	}
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

	return title, content , err
}
