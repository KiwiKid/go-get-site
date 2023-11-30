package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	log.Print("Start")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//// Uncomment to deploy DB changes (commentted out as it improves rebuild time)
	//// (Or comment to improve build speed)

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Print("Start:Migrate()")
		db, err := NewDB()
		if err != nil {
			panic(err)
		}
		log.Print("Start:Migrate()")
		db.Migrate()

	} else {
		log.Print("Skipping migration")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	r := mux.NewRouter()

	//ctx, cancel = chromedp.NewContext(ctx)
	// UNCOMMENT to add browser debugging
	debugBrowserLogsMode := os.Getenv("DEBUG_BROWSER_LOGS") == "true"

	if debugBrowserLogsMode {
		ctx, cancel = chromedp.NewContext(ctx, chromedp.WithDebugf(log.Printf))
	} else {
		ctx, cancel = chromedp.NewContext(ctx)
	}
	defer cancel()

	r.Handle("/", presentWebsite()).Methods("GET", "POST")
	r.HandleFunc("/sites/{websiteId}", presentWebsite()).Methods("GET", "POST", "PUT")
	r.HandleFunc("/sites/{websiteId}/{deleteOpt}", presentWebsite()).Methods("DELETE")
	r.HandleFunc("/sites/{websiteId}/{deleteOpt}/{pageId}", presentWebsite()).Methods("DELETE")

	r.HandleFunc("/sites/{websiteId}/blocks", presentWebsitePageBlocks()).Methods("GET", "DELETE", "POST")
	r.HandleFunc("/sites/{websiteId}/pages/{pageId}/blocks", presentPageBlocks()).Methods("GET", "DELETE", "POST")
	r.HandleFunc("/sites/{websiteId}/pages/{pageId}/blocks/{pageBlockId}", presentPageBlocks()).Methods("GET", "POST")

	r.HandleFunc("/sites/{websiteId}/login", presentLogin(ctx)).Methods("GET")

	r.HandleFunc("/sites/{websiteId}/pages", handlePages(ctx)).Methods("GET", "POST")

	r.HandleFunc("/sites/{websiteId}/result/{attributeSetId}", presentAttributeSetResult()).Methods("GET", "POST")

	r.Handle("/search", presentQuery()).Methods("GET", "POST")
	r.Handle("/search/{queryId}", presentQuery()).Methods("GET", "POST")

	//r.Handle("/attributes/models", presentAttributeModels())
	r.Handle("/aset", presentAttributeSet()).Methods("GET", "POST")
	r.Handle("/aset/{attributeSetId}", presentAttributeSet()).Methods("GET", "POST")

	r.Handle("/attributes", presentAttribute()).Methods("GET", "POST")
	//r.Handle("/attributes/{attributeSetId}/attribute/{attributeId}", presentAttribute())

	r.Handle("/sites/{websiteId}/pages/{pageId}/blocks/{pageBlockId}/questions", presentQuestion()).Methods("GET", "POST")

	r.Handle("/sites/{websiteId}/pages/{pageId}/blocks/{pageBlockId}/questions/{questionId}/improved", presentImprovedQuestions()).Methods("GET", "POST")

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
func presentQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]
		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			panic(err)
		}

		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		website, err := db.GetWebsite(websiteId)
		if err != nil {
			panic(err)
		}
		pageIdStr := vars["pageId"]
		pageId, err := stringToUint(pageIdStr)
		if err != nil {
			panic(err)
		}
		page, err := db.GetPage(websiteId, pageId)
		if err != nil {
			panic(err)
		}

		pageBlockIdStr := vars["pageBlockId"]
		pageBlockId, err := stringToUint(pageBlockIdStr)
		if err != nil {
			panic(err)
		}

		if r.Method == http.MethodPost {

			/*oai, err := NewOpenAI()
			if err != nil {
				failedQuestionRes := questionsFailedResult(*website, pageId, err.Error())
				templ.Handler(failedQuestionRes).ServeHTTP(w, r)
			}*/

			/*questions, err := oai.generateQuestions(page.Content, startIndex, endIndex)
			if err != nil {
				failedQuestionRes := questionsFailedResult(*website, pageId, err.Error())
				templ.Handler(failedQuestionRes).ServeHTTP(w, r)
			}

			if questions == nil {
				failedQuestionRes := questionsFailedResult(*website, pageId, "Failed to get questions")
				templ.Handler(failedQuestionRes).ServeHTTP(w, r)
				return
			}*/
			relevantContent := r.FormValue("relevantContent")

			if len(relevantContent) == 0 {
				failedQuestionRes := questionsFailedResult(website.ID, pageId, pageBlockId, "relevantContent is required")
				templ.Handler(failedQuestionRes).ServeHTTP(w, r)
			}
			id := uuid.New()
			inserQuestionErr := db.BatchInsertQuestions([]Question{{
				WebsiteID:       website.ID,
				PageID:          page.ID,
				BatchID:         id,
				PageBlockID:     pageBlockId,
				RelevantContent: relevantContent,
			}})
			if inserQuestionErr != nil {
				log.Fatalf("BatchInsertQuestions 1 has failed  %v", inserQuestionErr)
				http.Error(w, "BatchInsertQuestions 1 has failed", http.StatusBadRequest)
			}
		}

		questions, getQuestionErr := db.GetQuestions(website.ID, page.ID, pageBlockId, 1, 10000)
		if getQuestionErr != nil {
			http.Error(w, "GetQuestions has failed", http.StatusBadRequest)
			log.Fatalf("GetQuestions has failed  %v", getQuestionErr)
		}

		pageBlock, getPageErr := db.GetPageBlock(pageBlockId)
		if getPageErr != nil {
			log.Printf("getPage - failed %v", err)
			http.Error(w, "GetPage failed", http.StatusBadRequest)
			return
		}

		chatComp := questionResult(website.ID, page.ID, pageBlock.ID, questions, pageBlock.Content)
		templ.Handler(chatComp).ServeHTTP(w, r)
	}
}

func handleQuestionError(w http.ResponseWriter, r *http.Request, website *Website, pageId uint, err error) {
	// Assuming questionsFailedResult can handle a nil website

}

func presentQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("presentQuery - NewDB")
		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		vars := mux.Vars(r)
		queryIdStr := vars["queryId"]
		log.Printf("presentQuery - threadIdStr %v", queryIdStr)

		if queryIdStr == "" {
			log.Printf("queryHome - EMPTY threadIdStr %v", queryIdStr)

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

			queryContainerComp := queryContainer(chatThreads, newThreadURL, websites)

			templ.Handler(queryContainerComp).ServeHTTP(w, r)
			return
		}

		queryId, err := stringToUint(queryIdStr)
		if queryIdStr == "" || err != nil {

			chats, pageErr := db.ListChats(queryId)
			if err != nil {
				panic(pageErr)
			}

			websiteIdStr := strconv.Itoa(int(chats[0].WebsiteID))
			websiteId, err := stringToUint(websiteIdStr)
			if websiteIdStr == "" || err != nil {
				log.Printf("Failed based on websiteId %v", err)
				http.Error(w, "Failed based on websiteId", http.StatusInternalServerError)
				return
			}
			newChatUrl := chats[0].ChatURL()
			log.Printf("queryContainer - EMPTY threadIdStr %v", queryIdStr)

			chatComp := query(queryId, websiteId, newChatUrl, chats)

			templ.Handler(chatComp).ServeHTTP(w, r)

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
			_, insChatErr := db.InsertChat(Chat{ThreadId: queryId, Message: query, WebsiteID: websiteId})
			if insChatErr != nil {
				log.Printf("Failed to insert chat %v", insChatErr)
				http.Error(w, "InsertChat has failed", http.StatusBadRequest)
				return
			}
			var pageIds []uint
			pageIdsStr := r.FormValue("pageIds")
			if len(pageIdsStr) == 0 {
				log.Print("Inserted Chat.")
			} else {
				idStrs := strings.Split(pageIdsStr, ",")

				log.Printf("Inserted Chat. limit to these pages %s", pageIdsStr)
				for _, idStr := range idStrs {
					id, err := strconv.ParseUint(idStr, 10, 32)
					if err != nil {
						// Handle the error if the conversion fails
						http.Error(w, "Invalid page ID: "+idStr, http.StatusBadRequest)
						return
					}
					pageIds = append(pageIds, uint(id))
				}
			}

			// Split the string by commas

			// Initialize a slice to hold the uint values

			// Convert each string ID to uint and add to the slice

			log.Print("QueryWebsiteForPages. Started")
			queryRes, queryErr := db.QueryWebsiteForPages(query, websiteId, pageIds)
			if queryErr != nil {
				log.Printf("QueryWebsiteForPages. failed - %v", queryErr.Error())

				http.Error(w, "QueryWebsiteForPages failed queryErr\n\n"+queryErr.Error(), http.StatusBadRequest)
				return
			}

			log.Printf("QueryWebsiteForPages. End - %s", queryRes)

			queryComp := queryResult(queryRes, websiteId, query)

			templ.Handler(queryComp).ServeHTTP(w, r)

			//			insAIErr := db.InsertChat(Chat{ThreadId: threadId, Message: queryRes[0].String(), WebsiteID: websiteId})

			/*if insAIErr != nil {
				http.Error(w, "insAIErr failed", http.StatusBadRequest)
				return
			}*/
		}

	}
}

func presentImprovedQuestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("presentImprovedQuestions - NewDB")
		db, err := NewDB()
		if err != nil {
			log.Printf("NewDB failed %v", err)
			panic(err)
		}

		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]
		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			log.Printf("websiteId - failed %v", err)
			http.Error(w, "websiteId stringToUint is failed", http.StatusBadRequest)
			return
		}
		pageIdStr := vars["pageId"]
		pageId, pageIdErr := stringToUint(pageIdStr)
		if pageIdErr != nil {
			log.Printf("pageId - failed %v", pageIdErr)
			http.Error(w, "pageIdstringToUint for presentPageBlocks failed", http.StatusBadRequest)
			return
		}

		pageBlockIdStr := vars["pageBlockId"]
		pageBlockId, pageBlockIdErr := stringToUint(pageBlockIdStr)
		if pageBlockIdErr != nil {
			log.Printf("pageId - failed %v", pageBlockIdErr)
			http.Error(w, "stringToUint for pageBlockId failed", http.StatusBadRequest)
			return
		}

		questionIdStr := vars["questionId"]
		questionId, questionIdStrErr := stringToUint(questionIdStr)
		if questionIdStrErr != nil {
			log.Printf("questionIdStrErr - failed %v", questionIdStrErr)
			http.Error(w, "stringToUint for questionIdStrErr failed", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPost:
			mode := r.URL.Query().Get("mode")
			switch mode {
			case "new":
				{
					rawQuestionText := r.FormValue("rawQuestionText")
					rawQuestionAnswer := r.FormValue("rawQuestionAnswer")
					insErr := db.InsertImprovedQuestion(ImprovedQuestion{
						WebsiteID:          websiteId,
						PageID:             pageId,
						PageBlockID:        pageBlockId,
						QuestionID:         questionId,
						QuestionText:       rawQuestionText,
						CorrectAnswerText:  rawQuestionAnswer,
						IncorrectAnswerOne: "",
						IncorrectAnswerTwo: ""})
					if insErr != nil {
						log.Printf("Failed to insert chat %v", insErr)
						http.Error(w, "InsertChat has failed", http.StatusBadRequest)
						return
					}
				}
			case "update":
				{

				}
			case "gen-answers":
				{

				}
			default:
				{
					panic("no implemented, set a mode")
				}
			}
		}

		impQuestions, impQErr := db.GetImprovedQuestions(websiteId, pageId, pageBlockId)
		if impQErr != nil {
			log.Printf("GetImprovedQuestions - failed %v", impQErr)
			panic(impQErr)
		}
		log.Print("GetImprovedQuestions:end")

		impQComp := improvedQuestions(websiteId, pageId, pageBlockId, impQuestions)

		templ.Handler(impQComp).ServeHTTP(w, r)
	}
}

func presentPageBlocks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("presentPageBlocks - NewDB")
		db, err := NewDB()
		if err != nil {
			log.Printf("NewDB failed %v", err)
			panic(err)
		}

		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]
		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			log.Printf("websiteId - failed %v", err)
			http.Error(w, "websiteId stringToUint is failed", http.StatusBadRequest)
			return
		}
		pageIdStr := vars["pageId"]
		pageId, pageIdErr := stringToUint(pageIdStr)
		if pageIdErr != nil {
			log.Printf("pageId - failed %v", pageIdErr)
			http.Error(w, "pageIdstringToUint for presentPageBlocks failed", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPost:
			{

				page, getPageErr := db.GetPage(websiteId, pageId)
				if getPageErr != nil {
					log.Printf("getPage - failed %v", err)
					http.Error(w, "GetPage failed", http.StatusBadRequest)
					return
				}

				log.Printf("presentPageBlock:got-page %v length: %d", page.ID, len(page.Content))

				blocks := splitIntoBlocks(page.Content)
				var pageBlocks []PageBlock
				for _, blockContent := range blocks {
					pageBlock := PageBlock{
						WebsiteID: websiteId,
						PageID:    pageId,
						Content:   blockContent,
					}
					pageBlocks = append(pageBlocks, pageBlock)
				}

				log.Printf("presentPageBlock:got %d blocks", len(pageBlocks))
				if len(pageBlocks) > 0 {
					// Batch insert PageBlocks into the database
					blocks, err := db.BatchInsertPageBlocks(pageBlocks)

					if err != nil {
						log.Printf("BatchInsertPageBlocks - failed %v", err)
						http.Error(w, "BatchInsertPageBlocks failed", http.StatusInternalServerError)
						return
					}

					for _, newBlock := range blocks {

						id := uuid.New()

						genQuestions := r.URL.Query().Get("genQuestionss")
						if genQuestions == "on" {
							content := processBlockContent(page.TidyTitle, page.Content)
							inserQuestionErr := db.BatchInsertQuestions([]Question{{
								WebsiteID:       websiteId,
								PageID:          pageId,
								BatchID:         id,
								PageBlockID:     newBlock.ID,
								RelevantContent: content,
							}})
							if inserQuestionErr != nil {
								log.Printf("Error: BatchInsertQuestions 2 has failed  %v  with content: %s", inserQuestionErr, content)

								http.Error(w, "BatchInsertQuestions 2  has failed: %s", http.StatusBadRequest)
							}
						}
					}

				} else {
					log.Print("BatchInsertPageBlocks:none-to-add for \n%s", page.Content)
				}

				log.Print("BatchInsertPageBlocks:done:%d", len(pageBlocks))

				// Split the page.Content into block that are good for sentence embedding and batch insert PageBlocks into the database using db.BatchInsertPageBlocks()
				// Include a gorm type def for a PageBlock object, its should include a reference to the id of the 'beforeBlock' and 'AfterBlock'.

			}

		case http.MethodDelete:
			{
				delPageErr := db.DeletePageBlocks(pageId)
				if delPageErr != nil {
					log.Printf("DeletePageBlocks - failed %v", delPageErr)
					http.Error(w, "DeletePageBlocks failed", http.StatusInternalServerError)
				}
			}
		case http.MethodGet:
			{

			}
		}
		log.Print("ListPageBlocks:start")

		pageBlockList, pageErr := db.ListPageBlocks(pageId)
		if pageErr != nil {
			log.Printf("pageBlockList - failed %v", pageErr)
			panic(pageErr)
		}
		log.Print("ListPageBlocks:end")

		homeComp := pageBlocks(websiteId, pageId, pageBlockList)

		templ.Handler(homeComp).ServeHTTP(w, r)
	}
}

func presentWebsitePageBlocks() http.HandlerFunc {
	log.Print("presentWebsitePageBlocks")

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		websiteIdStr := vars["websiteId"]
		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			log.Printf("websiteId - failed %v", err)
			http.Error(w, "websiteId stringToUint is failed", http.StatusBadRequest)
			return
		}

		db, err := NewDB()
		if err != nil {
			log.Printf("NewDB failed %v", err)
			panic(err)
		}

		switch r.Method {
		case http.MethodPost:
			remainingStr := r.FormValue("remaining")

			remaining, err := strconv.Atoi(remainingStr)
			if err != nil || remaining <= 0 {
				remaining = 10
			}

			blockPageStr := r.FormValue("blockPage")

			blockPage, err := strconv.Atoi(blockPageStr)
			if err != nil || remaining <= 0 {
				blockPage = 1
			}

			blockLimitStr := r.FormValue("blockLimit")

			blockLimit, err := strconv.Atoi(blockLimitStr)
			if err != nil || remaining <= 0 {
				blockLimit = 10
			}

			nextPage, getPageErr := db.GetPages(websiteId, blockPage, blockLimit, false, time.Now(), false, 0)

			if getPageErr != nil {
				message := fmt.Sprintf("Failed to get any pages for  %v", getPageErr)

				pageLoaderInit := pageBlockLoaderInit(websiteBlocksURL(websiteId), message)
				templ.Handler(pageLoaderInit).ServeHTTP(w, r)
				return
			}

			if len(nextPage) > 0 {
				pageLoader := pageBlockLoader(websitePageBlocksURL(nextPage[0].WebsiteID, nextPage[0].ID), "load, every 2s", remaining)
				templ.Handler(pageLoader).ServeHTTP(w, r)
			} else {
				pageLoaderInit := pageBlockLoaderInit(websiteBlocksURL(websiteId), "No pages found to load blocks for")
				templ.Handler(pageLoaderInit).ServeHTTP(w, r)
			}

		}

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
				customQueryParam := r.FormValue("customQueryParam")
				loginName := r.FormValue("loginName")
				loginNameSelector := r.FormValue("loginNameSelector")
				loginPass := r.FormValue("loginPass")
				loginPassSelector := r.FormValue("loginPassSelector")
				preLoadPageClickSelector := r.FormValue("preLoadPageClickSelector")
				titleReplace := r.FormValue("titleReplace")

				startUrl := r.FormValue("startUrl")
				log.Printf("Form Values - websiteUrl: %s, customQueryParam: %s, loginName: %s, loginNameSelector: %s, loginPass: %s, loginPassSelector: %s",
					websiteUrl, customQueryParam, loginName, loginNameSelector, loginPass, loginPassSelector)

				site, err := db.InsertWebsite(Website{
					BaseUrl:                  websiteUrl,
					StartUrl:                 startUrl,
					CustomQueryParam:         customQueryParam,
					LoginName:                loginName,
					LoginPass:                loginPass,
					LoginNameSelector:        loginNameSelector,
					LoginPassSelector:        loginPassSelector,
					PreLoadPageClickSelector: preLoadPageClickSelector,
					TitleReplace:             titleReplace,
				})
				if err != nil {
					http.Error(w, "InsertWebsite is failed", http.StatusBadRequest)
					log.Fatalf("InsertWebsite  %v", err)

					return
				}
				emptyLink := []string{}
				inserPageErr := db.InsertPage(Page{URL: site.BaseUrl, WebsiteID: site.ID, Links: emptyLink})
				if inserPageErr != nil {
					http.Error(w, "InsertPage has failed", http.StatusBadRequest)
					log.Fatalf("InsertPage has failed  %v", inserPageErr)

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
				customQueryParam := r.FormValue("customQueryParam")
				loginName := r.FormValue("loginName")
				loginNameSelector := r.FormValue("loginNameSelector")
				loginPass := r.FormValue("loginPass")
				loginPassSelector := r.FormValue("loginPassSelector")
				submitButtonSelector := r.FormValue("submitButtonSelector")
				successIndicatorSelector := r.FormValue("successIndicatorSelector")
				preLoadPageClickSelector := r.FormValue("preLoadPageClickSelector")

				startUrl := r.FormValue("startUrl")

				site, err := db.UpdateWebsite(Website{
					ID:                       websiteId,
					BaseUrl:                  websiteUrl,
					CustomQueryParam:         customQueryParam,
					LoginName:                loginName,
					LoginPass:                loginPass,
					LoginNameSelector:        loginNameSelector,
					LoginPassSelector:        loginPassSelector,
					SubmitButtonSelector:     submitButtonSelector,
					SuccessIndicatorSelector: successIndicatorSelector,
					StartUrl:                 startUrl,
					PreLoadPageClickSelector: preLoadPageClickSelector,
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

						var url string
						if len(website.StartUrl) > 0 {
							url = website.StartUrl
						} else {
							url = website.BaseUrl
						}

						newEmptyPage := Page{
							URL:       url,
							WebsiteID: website.ID,
							Links:     emptyLink,
						}
						insertErr := db.InsertPage(newEmptyPage)
						if insertErr != nil {
							log.Printf("Error adding seed page after pages delete %v", insertErr)
							http.Error(w, "Error adding seed page", http.StatusBadRequest)

							return
						}
					}
				case "reset-page":
					{
						pageIdStr := vars["pageId"]
						pageId, err := stringToUint(pageIdStr)
						if err != nil {
							http.Error(w, "stringToUint for reset-page pageId is failed", http.StatusBadRequest)
							return
						}
						resetPageErr := db.ResetPage(pageId)

						if resetPageErr != nil {
							http.Error(w, "ResetPage has failed", http.StatusBadRequest)
							return
						}

					}
				case "warnings-reset":
					{
						resetPages := db.ResetWarnings(websiteId)
						if resetPages != nil {
							http.Error(w, "resetPages has failed", http.StatusBadRequest)
							return
						}

						r.Header.Add("HX-Redirect", fmt.Sprintf("/sites/%d/pages", websiteId))

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
	dripLoad := false
	dripLoadCount := 0
	skipNewLinkInsert := false
	processAll := false
	ignoreWarnings := false
	viewPageSize := 300
	processPageSize := 5
	dripLoadFreqMin := 5

	if err != nil {
		panic(err)
	}
	var message string
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
			http.Error(w, fmt.Sprintf("Failed to GetWebsite %v", err), http.StatusInternalServerError)
			return
		}

		pageStr := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}

		selectedAttributeSetIdStr := r.URL.Query().Get("selectedAttributeSetId")
		selectedAttributeSetId, selectedAttributeSetIdErr := stringToUint(selectedAttributeSetIdStr)

		viewPageSizeStr := r.URL.Query().Get("viewPageSize")
		viewPageSizeInt, err := strconv.Atoi(viewPageSizeStr)
		if err != nil || viewPageSizeInt <= 0 {
			log.Printf("no viewPageSize %v", err)
			viewPageSize = 300
		} else {
			viewPageSize = viewPageSizeInt
		}

		if r.Method == http.MethodPost {
			log.Print("handlePages - POST")
			db, err := NewDB()
			if err != nil {
				log.Printf("Failed to connect to the database %v", err)
				http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
				return
			}

			processAll = r.FormValue("processAll") == "on"
			dripLoad = r.FormValue("dripLoad") == "on"

			if dripLoad {

				dripLoadFreqMinStr := r.FormValue("dripLoadFreqMin")
				dripLoadFreqMinInt, err := strconv.Atoi(dripLoadFreqMinStr)
				if err != nil {
					log.Printf("Failed dripLoadFreqMinStr %s %v", dripLoadFreqMinStr, err)

					http.Error(w, "Failed dripLoadFreqMinStr", http.StatusInternalServerError)
					return
				}
				dripLoadFreqMin = dripLoadFreqMinInt

				dripLoadCountStr := r.FormValue("dripLoadCount")
				provideddripLoadCount, err := strconv.Atoi(dripLoadCountStr)
				if err != nil {
					log.Printf("Failed provideddripLoadCount %s %v", dripLoadCountStr, err)

					http.Error(w, "Failed provideddripLoadCount", http.StatusInternalServerError)
					return
				}
				dripLoadCount = provideddripLoadCount + 1

			}

			skipNewLinkInsert = r.FormValue("skipNewLinkInsert") == "on"

			ignoreWarnings = r.FormValue("ignoreWarnings") == "on"

			processPageSizeStr := r.FormValue("processPageSize")
			processPageSizeInt, err := strconv.Atoi(processPageSizeStr)
			if err != nil || processPageSizeInt <= 0 {
				processPageSize = 5
			} else {
				processPageSize = processPageSizeInt
			}

			log.Printf("handlePages - processWebsite processAll:%v", processAll)
			pagesProcessed, processErr := processWebsite(ctx, *db, *website, processAll, 1, processPageSize, addedPagesSet, skipNewLinkInsert, uint(dripLoadCount), ignoreWarnings)
			if processErr != nil {
				http.Error(w, "Failed to processWebsite", http.StatusInternalServerError)
				return
			}

			if pagesProcessed == nil && processErr == nil {
				log.Printf("handlePages:No more pages to process:%v", processAll)
				dripLoad = false
				message = "All Page Processing Complete"
			}

			w.WriteHeader(http.StatusCreated) // 201 Created status
		}

		/*	db, err := NewDB()
			if err != nil {
				panic(err)
			}*/

		log.Printf("ListPages %s %d %d", website.BaseUrl, page, viewPageSize)

		pagesList, pageErr := db.ListPages(website.ID, page, viewPageSize)
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

		thisPageUrl := fmt.Sprintf("/sites/%d/pages?page=%d&pageSize=%d", website.ID, page, viewPageSize)
		prevPageUrl := fmt.Sprintf("/sites/%d/pages?page=%d&pageSize=%d", website.ID, page-1, viewPageSize)

		var nextPageUrl string
		if count.TotalLinks > (page * viewPageSize) {
			nextPageUrl = fmt.Sprintf("/sites/%d/pages?page=%d&pageSize=%d&", website.ID, page, viewPageSize)
		} else {
			nextPageUrl = ""
		}

		if len(selectedAttributeSetIdStr) != 0 {
			if selectedAttributeSetIdErr != nil {
				panic(selectedAttributeSetIdErr)
			}
			thisPageUrl = fmt.Sprintf("%sselectedAttributeSetId=%b", thisPageUrl, selectedAttributeSetId)
			prevPageUrl = fmt.Sprintf("%sselectedAttributeSetId=%b", prevPageUrl, selectedAttributeSetId)
			nextPageUrl = fmt.Sprintf("%sselectedAttributeSetId=%b", nextPageUrl, selectedAttributeSetId)
		}

		percentage := fmt.Sprintf("%.0f", (float64(prog.Done) / float64(prog.Total) * 100.0))
		dripLoadStr := fmt.Sprintf("every %dm", dripLoadFreqMin)
		log.Printf("pagesList length: %d dripLoad %v dripLoadCount %d selectedAttributeSetId: %d", len(pagesList), dripLoad, dripLoadCount, selectedAttributeSetId)

		attributeSets, modelsErr := db.ListAllAttributeSets()
		if modelsErr != nil {
			log.Printf("Error ListAllAttributeSets link: %v", modelsErr)
			return
		}

		pagesComp := pages(
			pagesList,
			*website,
			*count,
			thisPageUrl,
			prevPageUrl,
			nextPageUrl,
			addedPagesSet,
			percentage,
			viewPageSize,
			processPageSize,
			dripLoad,
			dripLoadCount,
			dripLoadFreqMin,
			dripLoadStr,
			processAll,
			skipNewLinkInsert,
			ignoreWarnings,
			attributeSets,
			selectedAttributeSetId,
			message,
		)

		templ.Handler(pagesComp).ServeHTTP(w, r)
	}
}

func processWebsite(ctx context.Context, db DB, website Website, processAll bool, page int, pageSize int, addedPagesSet map[string]struct{}, skipNewLinkInsert bool, dripLoadCount uint, ignoreWarnings bool) ([]Page, error) {
	log.Print("StartProcessingSite ", website.BaseUrl)
	var pageProcessedAfter time.Time
	if processAll {
		pageProcessedAfter = time.Now().Add(-365 * 24 * time.Hour)
	} else {
		pageProcessedAfter = time.Now().Add(-7 * 24 * time.Hour)
	}

	pagesToProcess, err := db.GetPages(website.ID, page, 1000, processAll, pageProcessedAfter, ignoreWarnings, 0)
	linksAlreadyProcessed, apErr := db.GetCompletedPageUrls(website.ID)
	for _, url := range linksAlreadyProcessed {
		addedPagesSet[GetPageDoneCacheKey(website.ID, url)] = struct{}{}
	}
	log.Printf("GetPages got %d links to process [linksAlreadyProcessed:%d] [processAll:%v] [pageProcessedAfter:%v]", len(pagesToProcess), len(linksAlreadyProcessed), processAll, pageProcessedAfter)
	if err != nil || apErr != nil {
		log.Printf("Error GetLink from %v", err)
		return nil, err
	}

	if len(pagesToProcess) > 0 {
		// Just focus on content every second run
		// skipNewLinkInsert = skipNewLinkInsert || dripLoadCount%2 == 0
		pagesToSave, err := fetchContentFromPages(ctx, website, pagesToProcess, pageSize, addedPagesSet, skipNewLinkInsert)
		if err != nil {
			log.Printf("Error fetchContentFromPages %v", err)
			panic(err)
		}
		log.Printf("Got %d pagesToSave from fetchContentFromPages", len(pagesToSave))
		return pagesToSave, err
		/*for _, page := range pagesToSave {
			//if _, exists := addedPagesSet[page.URL]; !exists {
			log.Printf("UpsertPage: %s (T:%d, C:%d)", page.URL, len(page.Title), len(page.Content))
			insertErr := db.UpsertPage(page)
			if insertErr != nil {
				panic(insertErr)
			} // else {
			//	addedPagesSet[page.URL] = struct{}{}
			//}
			//}
		}*/
	} else {
		return nil, nil
	}
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
			logAction("getLoginTasks:Navigate-to-start-url:start", url, false),
			chromedp.Navigate(url),
			logAction("getLoginTasks:Navigate-to-start-url:success. Looking for successindicator:", website.SuccessIndicatorSelector, true),
			chromedp.WaitVisible(website.SuccessIndicatorSelector),
			logAction("getLoginTasks:Navigate-to-start-url:SUCCESS", "", false),
		)
	} else {
		log.Print("getLoginTasks:fetchContentFromPagesERROR-  no urlWithParams")
	}

	shouldLogin := website.LoginName != "" && website.LoginName != "dont-login"

	if shouldLogin {

		log.Printf("fetchContentFromPages Logging-in as '%s'", website.LoginName)

		tasks = append(tasks,
			logAction("getLoginTasks:LoginStart-looking-for:", website.LoginNameSelector, true),
			chromedp.WaitVisible(website.LoginNameSelector),
			logAction("getLoginTasks:LoginStart-LoginNameSelector-SUCCESS", website.LoginNameSelector, true),
			chromedp.SendKeys(website.LoginNameSelector, website.LoginName),
			logAction("getLoginTasks:LoginStart-LoginNameSelector-Entry-SUCCESS:", website.LoginNameSelector, true),
			chromedp.SendKeys(website.LoginPassSelector, website.LoginPass),
			logAction("getLoginTasks:LoginStart-LoginPassSelector-Entry-SUCCESS:", website.LoginPassSelector, true),
			logAction("getLoginTasks:LoginStart-looking-for SubmitButtonSelector: ", website.SubmitButtonSelector, true),
			chromedp.WaitVisible(website.SubmitButtonSelector),
			logAction("getLoginTasks:LoginStart-SubmitButtonSelector-found: ", website.SubmitButtonSelector, true),
			chromedp.Click(website.SubmitButtonSelector),
			logAction("getLoginTasks:LoginStart:SubmitButtonSelector-clicked. Correct Login & Pass? Looking for SuccessIndicatorSelector: ", website.SuccessIndicatorSelector, true),
			chromedp.WaitVisible(website.SuccessIndicatorSelector),
			logAction("getLoginTasks:LoginStart:SUCCESS", "", false),
		)
	}

	if website.RequestCookieName != "" && website.RequestCookieValue != "" {
		log.Printf("fetchContentFromPages Setting Cookies")

		tasks = append(tasks, SetCookie(website.RequestCookieName, website.RequestCookieValue, website.BaseUrl, "/", false, false))
	}

	return tasks

}

func fetchContentFromPages(ctx context.Context, website Website, pages []Page, remainingToProcess int, addedPagesSet map[string]struct{}, skipNewLinkInsert bool) ([]Page, error) {
	log.Print("fetchContentFromPages:start")

	db, dbErr := NewDB()
	if dbErr != nil {
		panic(dbErr)
	}

	var newPages []Page

	allLinksJS := `Array.from(document.querySelectorAll("*[href]")).map((i) => i.href)`

	/*autoUrl, addQueryErr := addQueryParam(urlStr, "1c8ca3a202b84c47961b79700b40f01a")
	if addQueryErr != nil {
		panic(addQueryErr)
	}*/
	tasks := getLoginTasks(website)

	for i, page := range pages {

		var title string
		var content string
		var links []string

		log.Printf("fetchContentFromPages Starting Page (%d/%d)[limit %d left] \n%s \n %s", i, len(pages), remainingToProcess, page.URL, website.StartUrl)

		// Add the rest of the tasks
		tasks = append(tasks,
			//	logAction("fetchContentFromPages: Got past auth? Looking for SuccessIndicatorSelector: ", website.SuccessIndicatorSelector, true),
			//	chromedp.WaitVisible(website.SuccessIndicatorSelector),
			logAction("fetchContentFromPages:Navigate-to:", page.URL, false),
			chromedp.Navigate(page.URL),
			logAction("fetchContentFromPages:Navigate-to:Success", page.URL, false),
			chromedp.WaitVisible(website.SuccessIndicatorSelector),
		)

		if website.PreLoadPageClickSelector != "" && website.PreLoadPageClickSelector != "NA" {
			tasks = append(tasks,
				logAction("fetchContentFromPages:PreLoadPageClickSelector:start:", website.PreLoadPageClickSelector, true),
				chromedp.ActionFunc(func(ctx context.Context) error {
					// Create a context with a timeout
					clickCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
					defer cancel()

					// Attempt to click. If the element is not found within the timeout, an error will be returned.
					err := chromedp.Click(website.PreLoadPageClickSelector, chromedp.NodeVisible).Do(clickCtx)
					if err != nil {
						// Log the error and return nil to continue with the next actions
						log.Printf("Optional click failed or timed out: %v", err)
						return nil
					}
					// If click was successful, wait for 5 seconds
					return chromedp.Sleep(5 * time.Second).Do(ctx)
				}),
				logAction("fetchContentFromPages:PreLoadPageClickSelector:end", page.URL, false),
			)
		}
		tasks = append(tasks,
			logAction("fetchContentFromPages:Navigate-to", "", false),
			chromedp.Sleep(2),
			logAction("fetchContentFromPages:eval:title- document.title", "", false),
			chromedp.Evaluate(`document.title`, &title),
			logAction("fetchContentFromPages:eval:title - document.title ", title, false),
			logAction("fetchContentFromPages:eval:content - document.body.textContent", "", false),
			chromedp.Evaluate(`document.body.textContent`, &content),
			logAction("fetchContentFromPages:eval:content", content, false),
			logAction(fmt.Sprintf("fetchContentFromPages:Got Title (%d) and Content (%d)", len(title), len(content)), page.URL, false),
			logAction("fetchContentFromPages:GetLinks", page.URL, false),
			chromedp.Evaluate(allLinksJS, &links),
			logAction("fetchContentFromPages:GotLinks", page.URL, false),
		)

		blockedUrlsStr := os.Getenv("BLOCKED_URLS")

		if blockedUrlsStr != "" {

			// Split the string back into an array
			blockedURLs := strings.Split(blockedUrlsStr, ",")
			tasks = append(tasks,
				network.Enable(),
				chromedp.ActionFunc(func(ctx context.Context) error {
					return network.SetBlockedURLS(blockedURLs).Do(ctx)
				}),
			)
		} else {
			fmt.Println("fetchContentFromPages:Environment variable BLOCKED_URLS is not set")
			panic("SET THE BLOCKED_URLS")
		}

		err := chromedp.Run(ctx,
			tasks...,
		)
		if err != nil {
			msg := fmt.Sprintf("Error fetching content for page %s: %s", page.URL, err.Error())

			setWarErr := db.UpdateWarning(page.ID, msg)
			if setWarErr != nil {
				panic(setWarErr)
			}

			if strings.Contains(err.Error(), "context deadline exceeded") {
				log.Printf("fetchContentFromPages:context deadline exceeded early exit")
				return newPages, nil
			} else {
				log.Printf("fetchContentFromPages:non-deadline error: %v", err.Error())
			}
			// panic(err)
			// return newPages, err
		} else {
			if len(content) > 0 {
				// BUILD A "Page" object
				newPage := Page{
					URL:       page.URL,
					Title:     title,
					TidyTitle: website.getTidyTitle(title),
					Content:   content,
					Links:     links,
					WebsiteID: website.ID,
				}
				log.Printf("fetchContentFromPages - content len: %d [pagecontent len: %d] ", len(content), len(newPage.Content))
				// ADD the Page object to the "pages" list
				newPages = append(newPages, newPage)

				insertErr := db.UpsertPage(newPage)
				if insertErr != nil {
					panic(insertErr)
				}
			} else {
				msg := fmt.Sprintf("Error fetching content for page (No-content) %s: %s", page.URL, err.Error())

				setWarErr := db.UpdateWarning(page.ID, msg)
				if setWarErr != nil {
					panic(setWarErr)
				}
			}

			emptyLink := []string{}

			if !skipNewLinkInsert {
				for _, link := range links {
					for _, baseUrl := range strings.Split(website.BaseUrl, ",") {
						link, err = stripAnchors(link)
						if err != nil {
							log.Printf("fetchContentFromPages:links: non-page link %s", link)
							panic(err)
						}
						if linkCouldBePage(link, baseUrl) {

							if _, exists := addedPagesSet[GetPageDoneCacheKey(website.ID, link)]; !exists {

								log.Printf("fetchContentFromPages:links: page link %s", link)

								newEmptyPage := Page{
									URL:       link,
									WebsiteID: website.ID,
									Links:     emptyLink,
								}
								newPages = append(newPages, newEmptyPage)
								linkInsertErr := db.UpsertPage(page)
								if linkInsertErr != nil {
									msg := fmt.Sprintf("fetchContentFromPages:links:Error Inserting linkCouldBePage page %s: %s", page.URL, linkInsertErr.Error())
									setWarErr := db.UpdateWarning(page.ID, msg)
									if setWarErr != nil {
										panic(setWarErr)
									}
								}
								addedPagesSet[GetPageDoneCacheKey(website.ID, link)] = struct{}{}
								// You might want to add this newPage to some slice or process it further
								break
							} else {
								log.Printf("fetchContentFromPages:links: already added %s", link)
							}
						} else {
							log.Printf("fetchContentFromPages:links: not a relvant link %s", link)
						}

					}
				}
			} else {
				log.Printf("fetchContentFromPages:links:skipNewLinkInsert - set to skip new inserts")
			}
		}
		if len(newPages) > remainingToProcess {
			log.Printf("fetchContentFromPages:Early exit %d paged processed", len(newPages))
			return newPages, nil
		} else {
			log.Printf("Continuing as we only have %d (limit %d)", len(newPages), remainingToProcess)
		}
	}

	log.Printf("Finished tasks (%d pages)", len(newPages))
	for _, pg := range newPages {
		log.Printf("- %s\n", pg.URL)
	}
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

func presentAttribute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var message string

		db, err := NewDB()
		if err != nil {
			panic(err)
		}
		switch r.Method {
		case http.MethodPost:

			vars := mux.Vars(r)
			attributeSetIdStr := vars["attributeSetId"]

			if len(attributeSetIdStr) == 0 {

				attributeModelIDStr := r.FormValue("attributeModelID")
				attributeModelID, attributeSetErr := stringToUint(attributeModelIDStr)
				if attributeSetErr != nil {
					panic(attributeSetErr)
				}

				aiSeedQuery := r.FormValue("attributeSeedQuery")
				aiTask := r.FormValue("aiTask")

				aiOptions := r.FormValue("aiOptions")

				createAttributeError := db.CreateAttribute(Attribute{
					AISeedQuery:      aiSeedQuery,
					AttributeModelID: attributeModelID,
					AITask:           aiTask,
					AIOptions:        aiOptions,
				})

				if createAttributeError != nil {
					panic(createAttributeError.Error)
				}
			} else {
				panic("update not supported")
			}

		}

		attrs, modelsErr := db.ListAttributes()
		if modelsErr != nil {
			log.Printf("Error ListAttributes link: %v", modelsErr)
		}

		attrModels, attModelsErr := db.ListAttributeModels()
		if attModelsErr != nil {
			log.Printf("Error ListAttributeModels link: %v", attModelsErr)
		}

		if len(attrs) == 0 {
			message = "(no attributes)"
		}

		if r.Method == "POST" {
			attrList := attributeList(attrs, attrModels, message)

			templ.Handler(attrList).ServeHTTP(w, r)

		} else {
			attrComp := attributeContainer(attrs, attrModels, message)

			templ.Handler(attrComp).ServeHTTP(w, r)
		}

	}

}

type AttributeResultsError struct {
	PageID         uint
	WebsiteID      uint
	AttributeSetID uint
	Message        string
}

func presentAttributeSetResult() http.HandlerFunc {

	log.Printf("presentAttributeSetResult link")
	return func(w http.ResponseWriter, r *http.Request) {

		db, err := NewDB()
		if err != nil {
			panic(err)
		}

		vars := mux.Vars(r)
		attributeSetIdStr := vars["attributeSetId"]
		attributeSetId, attributeSetIdErr := stringToUint(attributeSetIdStr)
		if attributeSetIdErr != nil {
			log.Printf("Failed to stringToUint attributeSetId, %v", err)
			http.Error(w, "Failed to stringToUint attributeSetId", http.StatusInternalServerError)
		}

		websiteIdStr := vars["websiteId"]
		websiteId, err := stringToUint(websiteIdStr)
		if err != nil {
			log.Printf("Failed to stringToUint websiteId, %v", err)
			http.Error(w, "Failed to stringToUint websiteId", http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodPost:
			{
				attributes, err := db.ListAttributesForSet(attributeSetId)
				if err != nil {
					log.Printf("Failed to ListAttributesForSet, %v", err)
					http.Error(w, "Failed ListAttributesForSet ", http.StatusInternalServerError)
				}

				if len(attributes) == 0 {
					log.Printf("No attribute in set")
					http.Error(w, "Add attrtibutes to the set ", http.StatusInternalServerError)
				}

				for _, attr := range attributes {
					pages, pageGetErr := db.GetPages(websiteId, 1, 3, true, time.Now(), true, attributeSetId)

					if pageGetErr != nil {
						log.Printf("Failed to GetPages, %v", pageGetErr)
						http.Error(w, "Failed to GetPages", http.StatusInternalServerError)
						return
					}

					if len(pages) == 0 {
						log.Printf("No pages on website?")
						http.Error(w, "Add pages to the website ", http.StatusInternalServerError)
						return
					}

					for _, page := range pages {

						if !page.GoodForSearch() {
							var message = fmt.Sprintf("page notGoodForSearch (no content) %d %s", page.ID, page.URL)

							setErrorComp := attributeSetError(AttributeResultsError{
								PageID:         page.ID,
								WebsiteID:      websiteId,
								AttributeSetID: attributeSetId,
								Message:        message,
							})

							templ.Handler(setErrorComp).ServeHTTP(w, r)
							continue
						}

						pageFilter := []uint{}
						pageFilter = append(pageFilter, page.ID)
						pageQueryResult, queryWebErr := db.GetRelatedPageBlocks(attr.AISeedQuery, "", websiteId, pageFilter)
						if queryWebErr != nil {
							log.Printf("Failed to GetRelatedPageBlocks, %v", queryWebErr)
							http.Error(w, "Failed to stringToUint websiteId", http.StatusInternalServerError)
							return
						}

						if len(pageQueryResult) == 0 {
							message := fmt.Sprintf("Failed pageQueryResult (no page blocks?) - no content matching: %s for page:%d and websiteId:%d", attr.AISeedQuery, page.ID, websiteId)

							setErrorComp := attributeSetError(AttributeResultsError{
								PageID:         page.ID,
								WebsiteID:      websiteId,
								AttributeSetID: attributeSetId,
								Message:        message,
							})

							templ.Handler(setErrorComp).ServeHTTP(w, r)

							log.Print(message)
							continue
						} else {

							result := fmt.Sprintf(`You are searching the ANSWER to the following QUESTION in the CONTENT below
						QUESTION: %s
						`, attr.AISeedQuery)

							var pageBlocksIds = ""
							for _, bp := range pageQueryResult {
								log.Printf("GOT CONTENT MATCHES: %s\n\n%v", bp.PageBlock.Content, bp)
								result += bp.PageBlock.Content + "\n"
								pageBlocksIds = fmt.Sprintf("%s,%d", pageBlocksIds, bp.PageBlock.ID)
							}

							prompt := fmt.Sprintf("CONTENT: %s", result)

							oai, err := NewOpenAI()

							if err != nil {
								http.Error(w, "Failed NewOpenAI", http.StatusInternalServerError)
							}

							res, err := oai.createChatCompletion(prompt)
							if err != nil {
								http.Error(w, "Failed createChatCompletion", http.StatusInternalServerError)
							}

							log.Printf("pageBlocksIds: %s\n", pageBlocksIds)

							createAttrErr := db.CreateAttributeResult(AttributeResult{
								PageID:          page.ID,
								WebsiteID:       websiteId,
								PageBlockIDsStr: pageBlocksIds,
								AttributeSetID:  attributeSetId,
								AttributeID:     attr.ID,
								AttributeResult: *res,
							})

							if createAttrErr != nil {
								log.Printf("Failed to CreateAttributeResult, %v", createAttrErr)
								http.Error(w, "Failed CreateAttributeResult", http.StatusInternalServerError)
							}
						}

					}

				}
			}
		case http.MethodGet:
			{

			}
		default:
			log.Printf("Unsupported method")
			http.Error(w, "Unsupported method", http.StatusInternalServerError)
			return
		}

		setResult, err := db.ListAttributeResults(attributeSetId, websiteId)
		if err != nil {
			var message = fmt.Sprintf("Failed to ListAttributeResults, %v", err)

			err := attributeSetErrorGeneral(websiteId, attributeSetId, message)
			templ.Handler(err).ServeHTTP(w, r)
			http.Error(w, "Failed to ListAttributeResults", http.StatusInternalServerError)
			return
		}

		setListContainerComp := attributeSetResult(setResult)

		templ.Handler(setListContainerComp).ServeHTTP(w, r)

	}
}

func presentAttributeSet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var message string

		db, err := NewDB()
		if err != nil {
			panic(err)
		}
		switch r.Method {
		case http.MethodPost:
			{
				vars := mux.Vars(r)
				attributeSetIdStr := vars["attributeSetId"]

				if len(attributeSetIdStr) > 0 {
					//attributeSetIdStr := r.FormValue("attributeSetId")
					attributeSetId, attributeSetIdErr := stringToUint(attributeSetIdStr)
					if attributeSetIdErr != nil {
						panic(attributeSetIdErr)
					}
					attributeIdStr := r.FormValue("attributeId")
					attributeId, attributeIdErr := stringToUint(attributeIdStr)
					if attributeIdErr != nil {
						panic(attributeIdErr)
					}

					log.Printf("CreateAttributeSetLink(%s, %s)", attributeSetIdStr, attributeIdStr)
					newSetLinkErr := db.CreateAttributeSetLink(AttributeSetLink{
						AttributeSetID: attributeSetId,
						AttributeID:    attributeId,
					})

					if newSetLinkErr != nil {
						panic(newSetLinkErr.Error)
					}

				} else {
					setName := r.FormValue("setName")
					newSet, newSetErr := db.CreateAttributeSet(AttributeSet{
						Name: setName,
					})

					if newSetErr != nil {
						panic(newSetErr.Error)
					}

					message = fmt.Sprintf("New set created: %s", newSet.Name)
				}
			}
		case http.MethodGet:
			{

			}
		case http.MethodDelete:
			{
				vars := mux.Vars(r)
				attributeSetIdStr := vars["attributeSetId"]
				attributeSetId, attributeSetIdErr := stringToUint(attributeSetIdStr)
				if attributeSetIdErr != nil {
					panic(attributeSetIdErr)
				}

				newSetErr := db.DeleteAttributeSet(attributeSetId)

				if newSetErr != nil {
					panic(newSetErr.Error)
				}

			}
		default:
			{
				panic("Not supported")
			}

		}

		sets, modelsErr := db.ListAllAttributeSets()
		if modelsErr != nil {
			log.Printf("Error ListAllAttributeSets link: %v", modelsErr)
			return
		}

		attributeModels, getModelsErr := db.ListAttributeModels()
		if getModelsErr != nil {
			log.Printf("Error ListAttributeModels link: %v", getModelsErr)
			return
		}

		attributes, getAttributesErr := db.ListAttributes()
		if getModelsErr != nil {
			log.Printf("Error ListAttributes link: %v", getAttributesErr)
			return
		}

		if r.Method == http.MethodDelete || r.Method == http.MethodPost {
			setListComp := attributeSetList(attributes, sets, attributeModels, message)

			templ.Handler(setListComp).ServeHTTP(w, r)

		}
		setListContainerComp := attributeSetContainer(attributes, sets, attributeModels, message)

		templ.Handler(setListContainerComp).ServeHTTP(w, r)
	}
}
