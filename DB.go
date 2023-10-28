package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Page struct {
	ID          uint      `gorm:"primary_key"`
	Title       string    `gorm:"size:255"`
	Keywords    string    `gorm:"size:255"`
	Content     string    `gorm:"type:text"`
	URL         string    `gorm:"size:255;unique"`
	Links       []string  `gorm:"type:text[]"`
	IsSeedUrl   bool      `gorm:"type:boolean;default:false;not null"`
	DateCreated time.Time `gorm:"type:timestamp"`
	DateUpdated time.Time `gorm:"type:timestamp"`
	WebsiteId   uint      `gorm:"index;not null"`
}

func (Page) TableName() string {
	return "pgml.page"
}

type Link struct {
	ID            uint      `gorm:"primary_key"`
	SourceURL     string    `gorm:"size:255"`
	URL           string    `gorm:"size:255;unique"`
	DateCreated   time.Time `gorm:"type:timestamp"`
	LastProcessed time.Time `gorm:"type:timestamp"`
	WebsiteId     uint      `gorm:"index;not null"`
}

func (Link) TableName() string {
	return "pgml.link"
}

type Website struct {
	ID               uint      `gorm:"primary_key"`
	CustomQueryParam string    `gorm:"size:255"`
	BaseUrl          string    `gorm:"size:255"`
	DateCreated      time.Time `gorm:"type:timestamp"`
}

func (Website) TableName() string {
	return "pgml.website"
}

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

func (w *ChatThread) ChatThreadURL() string {
	return fmt.Sprintf("/chat/%d", w.ThreadId)
}

func (w *Chat) ChatURL() string {
	return fmt.Sprintf("/chat/%d", w.ThreadId)
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
	db.conn.AutoMigrate(&Link{})
	db.conn.AutoMigrate(&Website{})
	db.conn.AutoMigrate(&Chat{})
	// Check if the index exists
	var count int64
	db.conn.Raw(`
        SELECT COUNT(*) 
        FROM pg_indexes 
        WHERE tablename = ? 
        AND indexname = ?
	`, "pgml.page", "idx_website_url").Scan(&count)

	if count == 0 {
		db.conn.Exec("CREATE INDEX idx_website_url ON pgml.page (url)")
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

func (db *DB) GetWebsite(id uint) (*Website, error) {
	var website Website
	result := db.conn.First(&website, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // return nil when not found
		}
		log.Print(result.Error)
		return nil, result.Error
	}
	return &website, nil
}

func (db *DB) InsertPage(page Page) error {
	page.DateCreated = time.Now()
	result := db.conn.Create(&page).Clauses(clause.OnConflict{UpdateAll: true})
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}

func (db *DB) UpdatePage(page Page) error {
	page.DateUpdated = time.Now()
	result := db.conn.Model(&page).Updates(page)
	return result.Error
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

func (db *DB) GetUnprocessedLinks(url string, limit int) ([]Link, error) {
	var links []Link

	urlQuer := url + "%"
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	// Refactor the JOIN and WHERE clauses to reflect the new requirement
	result := db.conn.Table("pgml.link").
		Joins("LEFT JOIN pgml.page ON pgml.page.url = pgml.link.url"). // Use LEFT JOIN
		Where("pgml.link.url LIKE ?", urlQuer).
		Where("pgml.page.url IS NULL OR LENGTH(pgml.page.content) = 0 OR pgml.page.content IS NULL"). // Either no matching record in page or content is NULL or has zero length
		Where("(pgml.link.last_processed IS NULL OR pgml.link.last_processed <= ?)", oneMonthAgo).    // Either last_processed is NULL or it's older than a month
		Limit(limit).
		Find(&links)

	if result.Error != nil {
		return nil, result.Error
	}

	return links, nil
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

type QueryResult struct {
	Result []byte
}

func (qr QueryResult) String() string {
	return string(qr.Result)
}

func (db *DB) QueryWebsite(question string) ([]QueryResult, error) {

	var queryResult []QueryResult

	// This SQL joins the chats table with a subquery that numbers each row per thread based on date created.
	// We then filter for rows that are numbered 1 to get the first message of each thread.
	rawSQL := `
	SELECT pgml.transform(

		inputs => ARRAY[
	
			'I am so excited to benchmark deep learning models in SQL. I can not wait to see the results!'
	
		],
	
		task   => '{
	
			"task": "text-classification", 
	
			"model": "cardiffnlp/twitter-roberta-base-sentiment"
	
		}'::JSONB
	
	);
	`
	/*
		SELECT chats.thread_id, chats.message as first_message, chats.date_created
		FROM (
			SELECT *, ROW_NUMBER() OVER(PARTITION BY thread_id ORDER BY date_created ASC) as rn
			FROM pgml.chat
		) AS chats
		WHERE chats.rn = 1
		ORDER BY chats.date_created ASC
	*/
	result := db.conn.Raw(rawSQL).Scan(&queryResult)

	if result.Error != nil {
		return nil, result.Error
	}

	return queryResult, nil
}

func (db *DB) InsertLink(link Link) error {
	result := db.conn.Create(&link).Clauses(clause.OnConflict{UpdateAll: true})
	return result.Error
}

func (db *DB) unProcessLink(url string) error {
	log.Print("unProcessLink")
	return db.conn.Model(&Link{}).Where("url = ?", url).Update("last_processed", "1970-01-01 00:00:00 UTC").Error
}

type LinkCountResult struct {
	TotalLinks     int
	LinksHavePages int
}

func (db *DB) CountLinksAndPages(websiteId uint) (*LinkCountResult, error) {
	log.Print(websiteId)
	// Count total links for the URL
	var totalLinks int64
	if err := db.conn.Table("pgml.link").Where("website_id = ?", websiteId).Count(&totalLinks).Error; err != nil {
		return nil, err
	}

	log.Print("CountLinksAndPages")
	log.Print(websiteId)

	// Count links that have pages for the URL
	var linksWithPages int64
	if err := db.conn.Table("pgml.link").
		Joins("INNER JOIN pgml.page ON pgml.page.website_id = pgml.link.website_id").
		Where("pgml.link.website_id = ?", websiteId).Where("LENGTH(pgml.page.content) > 0").
		Count(&linksWithPages).Error; err != nil {
		return nil, err
	}

	return &LinkCountResult{
		TotalLinks:     int(totalLinks),
		LinksHavePages: int(linksWithPages),
	}, nil
}
