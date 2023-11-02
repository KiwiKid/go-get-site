package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrArray []string

type Page struct {
	ID          uint            `gorm:"primary_key"`
	Title       string          `gorm:"size:255"`
	Keywords    string          `gorm:"size:255"`
	Content     string          `gorm:"type:text"`
	Embedding   []byte          `gorm:"-"`
	URL         string          `gorm:"size:255"`
	Links       JSONStringArray `gorm:"type:jsonb"`
	IsSeedUrl   bool            `gorm:"type:boolean;default:false;not null"`
	DateCreated time.Time       `gorm:"type:timestamp"`
	DateUpdated time.Time       `gorm:"type:timestamp"`
	WebsiteId   uint            `gorm:"index;not null"`

	_ struct{} `gorm:"unique_index:idx_website_url;column:website_id;column:url"`
}

func (Page) TableName() string {
	return "pgml.page"
}

func (p Page) ToProcess() bool {
	return len(p.Content) == 0 || p.DateUpdated.Before(time.Now().Add(-7*24*time.Hour))
}

type Website struct {
	ID                 uint      `gorm:"primary_key"`
	CustomQueryParam   string    `gorm:"size:1024"`
	BaseUrl            string    `gorm:"size:1024"`
	StartUrl           string    `gorm:"size:1024"`
	DateCreated        time.Time `gorm:"type:timestamp"`
	LoginName          string    `gorm:"size:255"`
	LoginPass          string    `gorm:"size:255"`
	RequestCookieName  string    `gorm:"size:255"`
	RequestCookieValue string    `gorm:"size:255"`
}

func (Website) TableName() string {
	return "pgml.website"
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
	return fmt.Sprintf("/site/%d", w.ID)
}

func (w *Website) websitePagesURL() string {
	return fmt.Sprintf("/site/%d/pages", w.ID)
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
	connStr := os.Getenv("DB_CONN_STR")
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &DB{conn: db}, nil
}

func (db *DB) Migrate() {
	// AutoMigrate will ONLY create tables, missing columns and missing indexes
	db.conn.AutoMigrate(&Page{})
	db.conn.AutoMigrate(&Website{})
	db.conn.AutoMigrate(&Chat{})
	//db.conn.AutoMigrate(&Embedding{})

	// Check if the index exists
	var count int64
	db.conn.Raw(`
        SELECT COUNT(*) 
        FROM pg_indexes 
        AND indexname = ?
	`, "pgml.page", "idx_website_url").Scan(&count)

	if count == 0 {
		db.conn.Exec("CREATE INDEX idx_website_url ON pgml.page (url)")
	}

	err := db.conn.Exec(`
		DO $$ 

		BEGIN 
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='page' AND column_name='embedding') THEN
				ALTER TABLE pgml.page DROP COLUMN embedding;
			END IF;
		END $$;
		
		-- Add the new generated embedding column if it doesn't exist
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='page' AND column_name='embedding') THEN
				ALTER TABLE pgml.page 
				ADD COLUMN embedding vector(384) GENERATED ALWAYS AS 
				(pgml.embed(transformer => 'sentence-transformers/multi-qa-MiniLM-L6-cos-v1'::text, 
							text => content, 
							kwargs => '{"device": "cuda"}'::jsonb)) STORED;
			END IF;
		END $$;
		`).Error

	if err != nil {
		panic(err)
	}
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
	page.DateUpdated = time.Now()
	page.DateCreated = time.Now()
	result := db.conn.Create(&page).Clauses(clause.OnConflict{UpdateAll: true})
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) UpsertPage(page Page) error {
	hasUpdatedRow, updateErr := db.UpdatePage(page)
	if !hasUpdatedRow {
		insertErr := db.InsertPage(page)
		if insertErr != nil {
			log.Printf("Error saving page 2 %v: %v", page.ID, updateErr)
			return updateErr
		}
	} else if updateErr != nil {
		panic(updateErr)
	}
	return nil
}

func (db *DB) UpdatePage(page Page) (bool, error) {
	log.Printf("Updating page: %s\n content len: %+v", page.URL, page)
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

func (db *DB) GetPages(websiteId uint, page int, limit int, processAll bool, afterDate time.Time) ([]Page, error) {
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
		query = query.Where("LENGTH(content) > 0")
		query = query.Where("date_updated > ?", afterDate)
	}

	var pages []Page
	err := query.
		Order("date_updated ASC").
		Limit(limit).
		Offset(offset).
		Find(&pages).Error

	return pages, err
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

func (db *DB) InsertChat(chat Chat) error {
	log.Printf("Saving chat %v", chat)
	chat.DateCreated = time.Now()
	result := db.conn.Create(&chat) //.Clauses(clause.OnConflict{UpdateAll: true})
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
	ID       int
	Content  string
	Title    string
	URL      string
	Keywords string
	Rank     float64
}

func (qr PageQueryResult) String() string {
	return qr.Content + " Rank:" + strconv.FormatFloat(qr.Rank, 'f', -1, 64)
}

func (db *DB) QueryWebsite(question string, websiteId uint) ([]PageQueryResult, error) {
	var queryResults []PageQueryResult

	rawSQL := `
    WITH request AS (
		SELECT pgml.embed(
			transformer => 'sentence-transformers/multi-qa-MiniLM-L6-cos-v1'::text, 
			text => $1,
			kwargs => '{"device": "cuda"}'::jsonb
		)::vector(384) AS embedding
	  )
	  SELECT
		id,
		content,
		title,
		url,
		keywords,
		embedding <=> (SELECT embedding FROM request) AS cosine_similarity
	  FROM pgml.page
	  WHERE website_id = $2
	  ORDER BY cosine_similarity ASC
	  LIMIT 3`

	rows, err := db.conn.Raw(rawSQL, question, websiteId).Rows() // Here we get the raw rows
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rank float64
		var Id int
		var url string
		var content string
		var title string
		var keywords string
		if err := rows.Scan(&Id, &content, &title, &url, &keywords, &rank); err != nil { // Manual scan into individual variables
			return nil, err
		}
		log.Printf("rank %v Id %v content %s\n", rank, Id, content)
		qr := PageQueryResult{ID: Id, Content: content, URL: url, Keywords: keywords, Rank: rank, Title: title}
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

	log.Print("CountLinksAndPages")

	// Count links that have pages for the URL
	var linksWithPages int64
	if err := db.conn.Table("pgml.page").
		Where("pgml.page.website_id = ?", websiteId).Where("LENGTH(pgml.page.content) > 0").
		Count(&linksWithPages).Error; err != nil {
		return nil, err
	}

	return &LinkCountResult{
		TotalLinks:     int(totalLinks),
		LinksHavePages: int(linksWithPages),
	}, nil
}
