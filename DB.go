package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrArray []string

type Page struct {
	ID            uint            `gorm:"primary_key"`
	Title         string          `gorm:"size:512"`
	TidyTitle     string          `gorm:"size:512;default:''"`
	Keywords      string          `gorm:"size:512"`
	Content       string          `gorm:"type:text"`
	Embedding     []byte          `gorm:"-"`
	URL           string          `gorm:"size:255"`
	Links         JSONStringArray `gorm:"type:jsonb"`
	Warning       string          `gorm:"size:1024"`
	IsSeedUrl     bool            `gorm:"type:boolean;default:false;not null"`
	DateCreated   time.Time       `gorm:"type:timestamp"`
	DateUpdated   time.Time       `gorm:"type:timestamp"`
	DateProcessed time.Time       `gorm:"type:timestamp"`
	WebsiteId     uint            `gorm:"index;not null"`

	_ struct{} `gorm:"unique_index:idx_website_url;column:website_id;column:url"`
}

type PageBlock struct {
	ID        uint   `gorm:"primary_key"`
	PageID    uint   `gorm:"index;not null"`
	WebsiteId uint   `gorm:"index;not null"`
	Content   string `gorm:"type:text"`
	Embedding []byte `gorm:"-"`
}

func (db *DB) ListPageBlocks(pageID uint) ([]PageBlock, error) {
	var pageBlocks []PageBlock
	result := db.conn.Where("page_id = ?", pageID).Find(&pageBlocks)
	return pageBlocks, result.Error
}

func (db *DB) BatchInsertPageBlocks(blocks []PageBlock) ([]PageBlock, error) {
	log.Printf("BatchInsertPageBlocks: Inserting %d blocks", len(blocks))

	if err := db.conn.Omit("Summary").Create(&blocks).Error; err != nil {
		return nil, err
	}

	return blocks, nil
}

func (db *DB) DeletePageBlocks(pageId uint) error {
	result := db.conn.Where("page_id = ?", pageId).Delete(&PageBlock{})
	return result.Error
}

func (Page) TableName() string {
	return "pgml.page"
}

type Question struct {
	ID              uint      `gorm:"primary_key"`
	PageID          uint      `gorm:"index"`           // Foreign key to Page URL
	PageBlockID     uint      `gorm:"index"`           // Foreign key to Page URL
	WebsiteID       uint      `gorm:"index"`           // Foreign key to Website ID
	BatchID         uuid.UUID `gorm:"type:uuid;index"` // To identify the batch of questions
	QuestionText    string    `gorm:"-:migration"`
	QuestionTextTwo string    `gorm:"-:migration"`
	RelevantContent string    `gorm:"type:text"`
	DateCreated     time.Time `gorm:"type:timestamp"`
	Status          string    `gorm:"size:255"`
}

type ImprovedQuestion struct {
	QuestionID         uint      `gorm:"index"`
	PageID             uint      `gorm:"index"` // Foreign key to Page URL
	PageBlockID        uint      `gorm:"index"` // Foreign key to Page URL
	WebsiteID          uint      `gorm:"index"` // Foreign key to Website ID
	QuestionText       string    `gorm:"type:text"`
	CorrectAnswerText  string    `gorm:"type:text"`
	IncorrectAnswerOne string    `gorm:"type:text"`
	IncorrectAnswerTwo string    `gorm:"type:text"`
	DateCreated        time.Time `gorm:"type:timestamp"`
	DateModified       time.Time `gorm:"type:timestamp"`
}

func (p Page) ToProcess() bool {
	return len(p.Content) == 0 || p.DateUpdated.Before(time.Now().Add(-7*24*time.Hour))
}

const format = "2006-01-02 15:04:05"

func (p Page) PageStatus() string {

	return fmt.Sprintf("[P:%s]", p.DateProcessed.Format(format))
}

type Website struct {
	ID                       uint      `gorm:"primary_key"`
	CustomQueryParam         string    `gorm:"size:1024"`
	BaseUrl                  string    `gorm:"size:1024"`
	StartUrl                 string    `gorm:"size:1024"`
	PreLoadPageClickSelector string    `gorm:"size:512"`
	RemoveAnchorTags         bool      `gorm:"type:boolean;default:true"`
	TitleReplace             string    `gorm:"size:512"`
	DateCreated              time.Time `gorm:"type:timestamp"`
	LoginName                string    `gorm:"size:512"`
	LoginNameSelector        string    `gorm:"size:512"`
	LoginPass                string    `gorm:"size:512"`
	LoginPassSelector        string    `gorm:"size:512"`
	SubmitButtonSelector     string    `gorm:"size:512"`
	SuccessIndicatorSelector string    `gorm:"size:512"`
	RequestCookieName        string    `gorm:"size:512"`
	RequestCookieValue       string    `gorm:"size:512"`
}

func (Website) TableName() string {
	return "pgml.website"
}

type AttributeModel struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:text"`
}

type AttributeSet struct {
	ID          uint `gorm:"primary_key"`
	Name        string
	DateCreated time.Time   `gorm:"type:timestamp"`
	Attributes  []Attribute `gorm:"foreignkey:AttributeSetID"`
}

type Attribute struct {
	ID                 uint   `gorm:"primary_key"`
	AttributeSeedQuery string `gorm:"type:text"`
	AttributeModelID   uint   `gorm:"index"`
	AttributeSetID     uint   `gorm:"index"`
}

type AttributeResult struct {
	ID              uint      `gorm:"primary_key"`
	PageID          uint      `gorm:"index"`
	PageBlockID     uint      `gorm:"index"`
	WebsiteID       uint      `gorm:"index"`
	AttributeID     uint      `gorm:"index"`
	AttributeResult string    `gorm:"type:text"`
	DateCreated     time.Time `gorm:"type:timestamp"`
}

/*
type Embedding struct {
	ID     uint   `gorm:"primary_key"`
	PageId uint   `gorm"index;not null"`
	Model  string `gorm:"size:255`

	DateCreated time.Time `gorm:"type:timestamp"`
}*/

type Chat struct {
	ID          uint      `gorm:"primary_key"`
	ThreadId    uint      `gorm:"primary_key"`
	WebsiteId   uint      `gorm:"primary_key"`
	Message     string    `gorm:"size:255"`
	DateCreated time.Time `gorm:"type:timestamp"`
}

func (Chat) TableName() string {
	return "pgml.chat"
}

func (w *Website) websiteURL() string {
	return fmt.Sprintf("/sites/%d", w.ID)
}

func (w *Website) getProcessURL() string {
	return fmt.Sprintf("/process/%d", w.ID)
}

func (w *Website) websiteURLWithPostFix(postfix string) string {
	return fmt.Sprintf("/sites/%d/%s", w.ID, postfix)
}

func websitePageBlocksURL(websiteId uint, pageId uint) string {
	return fmt.Sprintf("/sites/%d/pages/%d/blocks", websiteId, pageId)
}

func websitePageBlockQuestionURL(websiteId uint, pageId uint, pageBlockId uint) string {
	return fmt.Sprintf("/sites/%d/pages/%d/blocks/%d/questions", websiteId, pageId, pageBlockId)
}

func questionImprovement(websiteId uint, pageId uint, pageBlockId uint, questionId uint, mode string) string {
	if len(mode) > 0 {
		return fmt.Sprintf("/sites/%d/pages/%d/blocks/%d/questions/%d/improved?mode=%s", websiteId, pageId, pageBlockId, questionId, mode)
	} else {
		return fmt.Sprintf("/sites/%d/pages/%d/blocks/%d/questions/%d/improved", websiteId, pageId, pageBlockId, questionId)
	}

}

func websitePageBlockURL(websiteId uint, pageId uint, pageBlockId uint) string {
	return fmt.Sprintf("/sites/%d/pages/%d/blocks/%d", websiteId, pageId, pageBlockId)
}

func (w *Website) websiteURLWithPagesId(pageId uint) string {
	return fmt.Sprintf("/sites/%d/pages/%d", w.ID, pageId)
}

func (w *Website) websiteURLWithPagesIdAndPostFix(pageId uint, postfix string) string {
	return fmt.Sprintf("/sites/%d/pages/%d/%s", w.ID, pageId, postfix)
}

func (w *Website) websitePagesURL() string {
	return fmt.Sprintf("/sites/%d/pages", w.ID)
}

func attributeSetURL(attributeSetId uint) string {
	return fmt.Sprintf("/aset/%d", attributeSetId)
}

func (w *Website) websiteNavigateURL() string {
	if len(w.StartUrl) == 0 {
		return fmt.Sprintf("%s?%s", w.BaseUrl, w.CustomQueryParam)
	} else {
		return fmt.Sprintf("%s", w.StartUrl)
	}
}

func (w *Website) websiteLoginURL() string {
	return fmt.Sprintf("/sites/%d/login", w.ID)
}

func (w *ChatThread) ChatThreadURL() string {
	return fmt.Sprintf("/search/%d", w.ThreadId)
}

func (w *Chat) ChatURL() string {
	return fmt.Sprintf("/search/%d", w.ThreadId)
}

type JSONStringArray []string

// Implement the sql.Scanner interface
func (a *JSONStringArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), a)
}

// Implement the driver.Valuer interface
func (a JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type DB struct {
	conn *gorm.DB
}

// NewDB creates a new DB instance with GORM
func NewDB() (*DB, error) {
	log.Print("Init DB...")
	connStr := os.Getenv("DB_CONN_STR")
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Printf("Init DB..,failed %v...", err)
		panic(err)
	}
	log.Print("Init DB...Success")

	return &DB{conn: db}, nil
}

func (db *DB) Migrate() {
	// AutoMigrate will ONLY create tables, missing columns and missing indexes
	log.Print("Migrate ALL START")
	db.conn.AutoMigrate(&Page{})
	pageEmbeddingErr := db.conn.Exec(`
	-- Add the new generated embedding column if it doesn't exist
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='page' AND column_name='embedding') THEN
			ALTER TABLE pgml.page 
			ADD COLUMN embedding vector(384) GENERATED ALWAYS AS 
			(
				CASE
					WHEN content IS NOT NULL AND CONTENT <> '' THEN 
						pgml.embed(transformer => 'sentence-transformers/multi-qa-MiniLM-L6-cos-v1'::text, 
							text => content, 
							kwargs => '{"device": "cpu"}'::jsonb
						)
					ELSE NULL
				END
			) STORED;
		END IF;
	END $$;
	`).Error

	if pageEmbeddingErr != nil {
		log.Fatalf("Failed to load embeddingErr %v", pageEmbeddingErr)
	}
	log.Print("Migrate Start - Website")
	db.conn.AutoMigrate(&Website{})
	log.Print("Migrate Start - Chat")
	db.conn.AutoMigrate(&Chat{})
	log.Print("Migrate Start - Question")
	db.conn.AutoMigrate(&Question{})

	var columnExists int64
	db.conn.Raw(`
        SELECT COUNT(*) 
        FROM information_schema.columns
        WHERE table_name = 'questions'
		AND column_name='question_text'
	`).Scan(&columnExists)

	if columnExists == 0 {
		err := db.conn.Exec(`
			ALTER TABLE questions
			ADD COLUMN question_text text GENERATED ALWAYS AS 
			(
				CASE
					WHEN relevant_content IS NOT NULL AND relevant_content <> '' THEN 
						(pgml.transform(
							task => '{
								"task": "text2text-generation",
								"model": "potsawee/t5-large-generation-squad-QuestionAnswer"
							}'::JSONB, 
							inputs => ARRAY[ relevant_content ]
							,     args => '{
								"max_length" : 200
							}'::JSONB
						)::JSONB -> 0 ->> 'generated_text')
					ELSE NULL
				END
			) STORED;
			`).Error

		if err != nil {
			log.Fatalf("Failed to load questions.question_text column %v", err)
		}
	}

	var columnQtwoExists int64
	db.conn.Raw(`
        SELECT COUNT(*) 
        FROM information_schema.columns
        WHERE table_name = 'questions'
		AND column_name='question_text_two'
	`).Scan(&columnQtwoExists)

	if columnQtwoExists == 0 {
		qtextErr := db.conn.Exec(`
		ALTER TABLE questions
		ADD COLUMN question_text_two text GENERATED ALWAYS AS 
		(
			CASE
				WHEN relevant_content IS NOT NULL AND relevant_content <> '' THEN 
					(pgml.transform(
						task => '{
							"task": "text2text-generation"
						}'::JSONB, 
						inputs => ARRAY[ relevant_content ]
						, args => '{
							"max_length" : 100
						}'::JSONB
					)::JSONB -> 0 ->> 'generated_text')
				ELSE NULL
			END
		) STORED;
		`).Error

		if qtextErr != nil {
			log.Fatalf("Failed to load questions.question_text_two column %v", qtextErr)
		}
	}
	log.Print("Migrate Start - PageBlock")
	db.conn.AutoMigrate(&PageBlock{})
	/*
		pageBlockSummaryTriggerErr := db.conn.Exec(`

		CREATE OR REPLACE FUNCTION update_page_block_summary() RETURNS TRIGGER AS $$
		BEGIN
			IF NEW.content IS NOT NULL AND NEW.content <> '' THEN
				NEW.summary := pgml.transform(
					task => '{
						"task": "summarization",
						"model": "sshleifer/distilbart-cnn-12-6"
					}'::JSONB,
					inputs => ARRAY[ NEW.content ]
				);
			ELSE
				NEW.summary := NULL;
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
		`).Error

		if pageBlockSummaryTriggerErr != nil {
			log.Fatalf("Failed CREATE OR REPLACE FUNCTION update_page_block_summary() %v", pageBlockSummaryTriggerErr)
		}

		pageBlockSummaryTriggerAddErr := db.conn.Exec(`
				CREATE OR REPLACE TRIGGER update_page_block_summary
				BEFORE INSERT OR UPDATE ON page_blocks
				FOR EACH ROW EXECUTE FUNCTION update_page_block_summary();
			`).Error

		if pageBlockSummaryTriggerAddErr != nil {
			log.Fatalf("CREATE OR REPLACE TRIGGER update_page_block_summary %v", pageBlockSummaryTriggerAddErr)
		}*/

	/*
		,
								args => '{
									"min_length" : 3,
									"max_length" : 15
								}'::JSONB
	*/

	log.Print("Migrate Start - ImprovedQuestion")
	db.conn.AutoMigrate(&ImprovedQuestion{})

	log.Print("Migrate Start - Attribute")
	db.conn.AutoMigrate(&AttributeSet{}, &AttributeModel{})
	db.conn.AutoMigrate(&Attribute{}, &AttributeResult{})
	log.Print("Migrate Start - Question")
	db.conn.AutoMigrate(&Question{})

	log.Print("Migrate ALL END")

	// Check if the index exists
	var count int64
	db.conn.Raw(`
        SELECT COUNT(*) 
        FROM pg_indexes 
        WHERE indexname = ?
	`, "pgml.page", "website_url").Scan(&count)

	if count == 0 {
		db.conn.Exec("CREATE INDEX CONCURRENTLY website_url ON pgml.page (url)")
	}

	/*
			DO $$

		--	BEGIN
		--		IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='page' AND column_name='embedding') THEN
		--			ALTER TABLE pgml.page DROP COLUMN embedding;
		--		END IF;
		--	END $$;


	*/

	pageBlockEmbeddingErr := db.conn.Exec(`
	-- Add the new generated embedding column if it doesn't exist
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='page_blocks' AND column_name='embedding') THEN
			ALTER TABLE page_blocks 
			ADD COLUMN embedding vector(384) GENERATED ALWAYS AS 
			(
				CASE
					WHEN content IS NOT NULL AND CONTENT <> '' THEN 
						pgml.embed(transformer => 'sentence-transformers/multi-qa-MiniLM-L6-cos-v1'::text, 
							text => content, 
							kwargs => '{"device": "cpu"}'::jsonb
						)
					ELSE NULL
				END
			) STORED;
		END IF;
	END $$;
	`).Error

	if pageBlockEmbeddingErr != nil {
		log.Fatalf("Failed to load pageBlockEmbeddingErr %v", pageBlockEmbeddingErr)
	}

	log.Print("Migrate End - Index")

}

func (db *DB) InsertWebsite(website Website) (Website, error) {
	website.DateCreated = time.Now()
	result := db.conn.Create(&website).Clauses(clause.OnConflict{UpdateAll: true})
	if result.Error != nil {
		log.Print(result.Error)
		return Website{}, result.Error
	}
	return website, nil
}

func (db *DB) UpdateWebsite(website Website) (Website, error) {
	// Fetch the existing record to ensure it exists
	var existing Website
	if err := db.conn.First(&existing, website.ID).Error; err != nil {
		log.Print(err)
		return Website{}, err
	}

	// Update the record
	result := db.conn.Model(&existing).Updates(website)
	if result.Error != nil {
		log.Print(result.Error)
		return Website{}, result.Error
	}

	return website, nil
}

func (db *DB) DeleteWebsite(websiteId uint) error {
	chatDelErr := db.conn.Delete(Chat{}, "website_id = ?", websiteId)
	if chatDelErr.Error != nil {
		log.Print(chatDelErr.Error)
		return chatDelErr.Error
	}
	pageDelErr := db.conn.Delete(Page{}, "website_id = ?", websiteId)
	if pageDelErr.Error != nil {
		log.Print(pageDelErr.Error)
		return pageDelErr.Error
	}
	webDelErr := db.conn.Delete(Website{}, "id = ?", websiteId)
	if webDelErr.Error != nil {
		log.Print(webDelErr.Error)
		return webDelErr.Error
	}

	return nil
}

func (db *DB) BatchInsertQuestions(questions []Question) error {
	// Perform the batch insert
	return db.conn.Omit("QuestionText", "QuestionTextTwo").Create(&questions).Error
}

func (db *DB) GetQuestions(websiteId uint, pageID uint, pageBlockId uint, page int, limit int) ([]Question, error) {
	if limit <= 0 {
		return nil, errors.New("invalid limit value")
	}
	if page <= 0 {
		return nil, errors.New("invalid page value")
	}
	offset := (page - 1) * limit

	var questions []Question
	err := db.conn.Raw("SELECT * FROM questions WHERE website_id = ? AND page_block_id = ? LIMIT ? OFFSET ?", websiteId, pageBlockId, limit, offset).Scan(&questions).Error

	log.Printf("questions: %v", questions)
	return questions, err
}

func (db *DB) DeletePages(websiteId uint) error {
	pageDelErr := db.conn.Delete(Page{}, "website_id = ?", websiteId)
	if pageDelErr.Error != nil {
		log.Print(pageDelErr.Error)
		return pageDelErr.Error
	}
	return nil
}

func (db *DB) GetWebsite(id uint) (*Website, error) {
	var website Website
	result := db.conn.First(&website, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &website, result.Error
}

func (db *DB) InsertPage(page Page) error {
	//log.Printf("Updating page: %+v\n", page)

	result := db.conn.Create(&page).Clauses(clause.OnConflict{UpdateAll: true})
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) UpdateWarning(pageID uint, warning string) error {
	page := Page{ID: pageID, Warning: warning}

	result := db.conn.Model(&page).Where("id = ?", pageID).Update("warning", warning)

	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were updated
	if result.RowsAffected == 0 {
		return fmt.Errorf("No page found with ID %d to update warning", pageID)
	}

	return nil
}

func (db *DB) ResetWarnings(websiteID uint) error {
	// Update the Warning field for all pages with the specified website_id
	result := db.conn.Model(&Page{}).Where("website_id = ?", websiteID).Update("warning", "")

	// Check for an error in the update operation
	if result.Error != nil {
		return result.Error
	}

	// Optionally, you can log the number of rows affected
	log.Printf("Warnings reset for %v pages of website ID %v", result.RowsAffected, websiteID)

	return nil
}

func (db *DB) ResetPage(pageID uint) error {
	// Fetch the page by ID
	var page Page
	// Reset the page fields
	page.DateUpdated = time.Now()
	page.DateCreated = time.Now()
	page.Content = ""                                                           // Assuming you want to clear the content
	page.DateProcessed = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC) // Reset to zero time if DateProcessed is a time.Time
	page.Warning = ""

	// Update the page in the database
	updateResult := db.conn.Save(&page)
	if updateResult.Error != nil {
		return updateResult.Error
	}

	return nil
}

func (db *DB) UpsertPage(page Page) error {
	log.Printf("UpsertPage - inserting page %v: %v", page.ID, page.URL)
	page.DateUpdated = time.Now()
	page.DateCreated = time.Now()
	if len(page.Content) > 0 {
		page.DateProcessed = time.Now()
	}
	hasUpdatedRow, updateErr := db.UpdatePage(page)
	if !hasUpdatedRow {
		log.Printf("No page to update - inserting page %v: %v", page.ID, updateErr)

		insertErr := db.InsertPage(page)
		if insertErr != nil {
			log.Printf("Error inserting page 2 %v: %v", page.ID, updateErr)
			return updateErr
		}
	} else if updateErr != nil {
		log.Printf("Error updateErr page 2 %v: %v", page.ID, updateErr)

		panic(updateErr)
	}
	return nil
}

func (db *DB) UpdatePage(page Page) (bool, error) {
	log.Printf("Updating page: %s\n content len: %d %s", page.URL, len(page.Content), page.Title)
	page.DateUpdated = time.Now()
	page.DateCreated = time.Now()
	result := db.conn.Model(&page).Where("pgml.page.url = ?", page.URL).Where("pgml.page.website_id= ?", page.WebsiteId).Updates(page)

	if result.Error != nil {
		return false, result.Error
	}

	// Check if any rows were updated
	if result.RowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (db *DB) ListWebsites() ([]Website, error) {
	// Note the struct tag to indicate the column name
	var websites []Website
	//result := db.conn.Where("IsSeedUrl = ?", true).Distinct("url")
	result := db.conn.Table("pgml.website").Select("*").Scan(&websites)
	if result.Error != nil {
		panic(result.Error)
	}

	return websites, nil
}

func (db *DB) GetPages(websiteId uint, page int, limit int, processAll bool, afterProcessDate time.Time, ignoreWarnings bool) ([]Page, error) {
	if limit <= 0 {
		return nil, errors.New("invalid limit value")
	}
	if page <= 0 {
		return nil, errors.New("invalid page value")
	}
	offset := (page - 1) * limit

	query := db.conn.
		Where("website_id = ?", websiteId)

	if !processAll {
		query = query.Where("LENGTH(content) = 0")
		if !ignoreWarnings {
			query = query.Where("LENGTH(warning) = 0")
		}
		//	query = query.Where("date_processed < ?", afterProcessDate)
	}

	var pages []Page
	err := query.
		Order("LENGTH(warning) ASC, date_updated ASC").
		Limit(limit).
		Offset(offset).
		Find(&pages).Error

	return pages, err
}

func (db *DB) GetCompletedPageUrls(websiteId uint) ([]string, error) {
	var urls []string
	result := db.conn.Table("pgml.page").Select("url").Where("LENGTH(content) > 0 AND website_id = ?", websiteId).Scan(&urls)
	if result.Error != nil {
		return nil, result.Error
	}
	return urls, nil
}

func (db *DB) SetLinkProcessed(url string) error {
	log.Print("SetLinkProcessed\n\n\n\nSetLinkProcessed")
	result := db.conn.Table("pgml.link").Where("pgml.link.url = ? ", url).Update("last_processed", time.Now())
	return result.Error
}

func (db *DB) ListPages(websiteId uint, page int, pageSize int) ([]Page, error) {
	var pages []Page
	offset := (page - 1) * pageSize
	result := db.conn.Where("website_id = ?", websiteId).Offset(offset).Limit(pageSize).Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pages, nil
}

func (db *DB) GetPage(websiteId uint, pageId uint) (*Page, error) {
	var page Page
	result := db.conn.Where("website_id = ? AND id = ?", websiteId, pageId).First(&page)
	if result.Error != nil {
		return nil, result.Error
	}
	return &page, nil
}

func (db *DB) GetPageBlock(pageBlockId uint) (*PageBlock, error) {
	var pageBlock PageBlock
	result := db.conn.Where("id = ?", pageBlockId).First(&pageBlock)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageBlock, nil
}

func (db *DB) InsertChat(chat Chat) (*Chat, error) {
	log.Printf("Saving chat %v", chat)
	chat.DateCreated = time.Now()
	result := db.conn.Create(&chat) //.Clauses(clause.OnConflict{UpdateAll: true})
	if result.Error != nil {
		log.Print(result.Error)
		return nil, result.Error
	}
	return &chat, nil
}

func (db *DB) GetImprovedQuestions(websiteId uint, pageId uint, pageBlockId uint) ([]ImprovedQuestion, error) {
	var questions []ImprovedQuestion
	err := db.conn.Where("website_id = ? AND page_id = ? AND page_block_id = ?", websiteId, pageId, pageBlockId).Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (db *DB) InsertImprovedQuestion(q ImprovedQuestion) error {
	log.Printf("Saving ImprovedQuestion %v", q)
	q.DateCreated = time.Now()
	result := db.conn.Create(&q) //.Clauses(clause.OnConflict{UpdateAll: true})
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) ListChats(threadId uint) ([]Chat, error) {
	var chats []Chat
	result := db.conn.Where("thread_id = ?", threadId).Find(&chats)
	if result.Error != nil {
		return nil, result.Error
	}
	return chats, nil
}

type ChatThread struct {
	ThreadId     uint
	FirstMessage string
	DateCreated  time.Time
}

// List chat threads, first message, distict by ThreadId
func (db *DB) ListChatThreads() ([]ChatThread, error) {
	var chatThreads []ChatThread

	// This SQL joins the chats table with a subquery that numbers each row per thread based on date created.
	// We then filter for rows that are numbered 1 to get the first message of each thread.
	rawSQL := `
	SELECT chats.thread_id, chats.message as first_message, chats.date_created
	FROM (
		SELECT *, ROW_NUMBER() OVER(PARTITION BY thread_id ORDER BY date_created ASC) as rn 
		FROM pgml.chat
	) AS chats
	WHERE chats.rn = 1
	ORDER BY chats.date_created ASC
	`

	result := db.conn.Raw(rawSQL).Scan(&chatThreads)

	if result.Error != nil {
		return nil, result.Error
	}
	return chatThreads, nil
}

type PageQueryResult struct {
	ID        int
	Content   string
	TidyTitle string
	URL       string
	Keywords  string
	BlockRank float64
	PageRank  float64
}

func (qr PageQueryResult) String() string {
	return qr.Content + " PageRank:" + strconv.FormatFloat(qr.PageRank, 'f', -1, 64) + " BlockRank:" + strconv.FormatFloat(qr.BlockRank, 'f', -1, 64)
}

func (db *DB) QueryWebsite(question string, websiteId uint, pageIds []uint) ([]PageQueryResult, error) {
	var queryResults []PageQueryResult

	var pageQuery string
	var pageIdStrs []string
	for _, id := range pageIds {
		pageIdStrs = append(pageIdStrs, strconv.Itoa(int(id)))
	}

	if len(pageIds) > 0 {
		pageQuery = fmt.Sprintf("AND p.id IN ARRAY[%s]", strings.Join(pageIdStrs, ","))
	} else {
		pageQuery = ""
	}

	rawSQL := fmt.Sprintf(`
    WITH request AS (
		SELECT pgml.embed(
			transformer => 'sentence-transformers/multi-qa-MiniLM-L6-cos-v1'::text, 
			text => $1,
			kwargs => '{"device": "cpu"}'::jsonb
		)::vector(384) AS embedding
	  )
	  SELECT
	  	pb.id,
		pb.content,
	  	p.tidy_title,
		p.URL,
		p.keywords,
		p.embedding <=> (SELECT embedding FROM request) AS page_cosine_similarity
	  FROM page_blocks pb
	  INNER JOIN pgml.page p ON p.id = pb.page_id
	  WHERE p.website_id = $2
	  %s
	  ORDER BY page_cosine_similarity ASC
	  LIMIT 10`, pageQuery)
	// 	p.embedding <=> (SELECT embedding FROM request) AS page_cosine_similarity,
	rows, err := db.conn.Raw(rawSQL, question, websiteId).Rows() // Here we get the raw rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// var page_cosine_similarity float64
		var page_cosine_similarity float64
		var Id int
		var url string
		var content string
		var tidy_title string
		var keywords string
		if err := rows.Scan(&Id, &content, &tidy_title, &url, &keywords, &page_cosine_similarity); err != nil { // Manual scan into individual variables
			return nil, err
		}
		log.Printf("brank %v Id %v content %s\n", page_cosine_similarity, Id, content)
		qr := PageQueryResult{ID: Id, Content: content, URL: url, Keywords: keywords, PageRank: page_cosine_similarity, TidyTitle: tidy_title, BlockRank: 0}
		queryResults = append(queryResults, qr)
	}

	return queryResults, nil
}

type LinkCountResult struct {
	TotalLinks     int
	LinksHavePages int
}

func (db *DB) CountLinksAndPages(websiteId uint) (*LinkCountResult, error) {
	log.Print(websiteId)
	// Count total links for the URL
	var totalLinks int64
	if err := db.conn.Table("pgml.page").Where("pgml.page.website_id = ?", websiteId).Count(&totalLinks).Error; err != nil {
		return nil, err
	}

	// Count links that have pages for the URL
	var linksWithPages int64
	if err := db.conn.Table("pgml.page").
		Where("pgml.page.website_id = ?", websiteId).Where("LENGTH(pgml.page.content) > 0").
		Count(&linksWithPages).Error; err != nil {
		return nil, err
	}

	log.Printf("CountLinksAndPages %d %d", totalLinks, linksWithPages)

	return &LinkCountResult{
		TotalLinks:     int(totalLinks),
		LinksHavePages: int(linksWithPages),
	}, nil
}

// For AttributeModel
func (db *DB) CreateAttributeModel(m AttributeModel) error {
	log.Printf("Creating AttributeModel %v", m)
	result := db.conn.Create(&m)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) DeleteAttributeModel(id uint) error {
	result := db.conn.Delete(&AttributeModel{}, id)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) ListAttributeModels() ([]AttributeModel, error) {
	var models []AttributeModel
	err := db.conn.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

// For AttributeSet
func (db *DB) CreateAttributeSet(s AttributeSet) (*AttributeSet, error) {
	log.Printf("Creating AttributeSet %v", s)
	s.DateCreated = time.Now()
	result := db.conn.Create(&s)
	if result.Error != nil {
		log.Print(result.Error)
		return nil, result.Error
	}
	return &s, nil
}

func (db *DB) DeleteAttributeSet(id uint) error {
	result := db.conn.Delete(&AttributeSet{}, id)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) ListAllAttributeSets() ([]AttributeSet, error) {
	var sets []AttributeSet
	err := db.conn.Preload("Attributes").Find(&sets).Error
	if err != nil {
		return nil, err
	}
	return sets, nil
}

func (db *DB) ListAttributes() ([]Attribute, error) {
	var attributes []Attribute
	err := db.conn.Find(&attributes).Error
	if err != nil {
		return nil, err
	}
	return attributes, nil
}

// For AttributeResult
func (db *DB) CreateAttributeResult(r AttributeResult) error {
	log.Printf("Creating AttributeResult %v", r)
	r.DateCreated = time.Now()
	result := db.conn.Create(&r)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) DeleteAttributeResult(id uint) error {
	result := db.conn.Delete(&AttributeResult{}, id)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) ListAttributeResults() ([]AttributeResult, error) {
	var results []AttributeResult
	err := db.conn.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// For Attribute
func (db *DB) CreateAttribute(a Attribute) error {
	log.Printf("Creating Attribute %v", a)
	result := db.conn.Create(&a)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) DeleteAttribute(id uint) error {
	result := db.conn.Delete(&Attribute{}, id)
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}
