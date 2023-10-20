package main

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Page struct {
	ID          uint      `gorm:"primary_key"`
	Title       string    `gorm:"size:255"`
	Keywords    string    `gorm:"size:255"`
	Content     string    `gorm:"type:text"`
	URL         string    `gorm:"size:255"`
	IsSeedUrl   bool      `gorm:"type:boolean;default:false;not null"`
	DateCreated time.Time `gorm:"type:timestamp"`
	DateUpdated time.Time `gorm:"type:timestamp"`
}

func (Page) TableName() string {
	return "pgml.page"
}

type Link struct {
	ID            uint      `gorm:"primary_key"`
	URL           string    `gorm:"size:255"`
	DateCreated   time.Time `gorm:"type:timestamp"`
	LastProcessed time.Time `gorm:"type:timestamp"`
}

func (Link) TableName() string {
	return "pgml.link"
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
}

func (db *DB) InsertPage(page Page) error {
	result := db.conn.Create(&page)
	return result.Error
}

func (db *DB) ListWebsites() ([]string, error) {
	// Note the struct tag to indicate the column name
	var websites []struct {
		WebsiteUrl string `gorm:"column:url"`
	}
	//result := db.conn.Where("IsSeedUrl = ?", true).Distinct("url")
	result := db.conn.Table("pgml.page").Select("distinct(url)").Where("is_seed_url = true").Scan(&websites)
	if result.Error != nil {
		panic(result.Error)
	}

	websiteUrls := make([]string, len(websites))
	for _, website := range websites {
		websiteUrls = append(websiteUrls, website.WebsiteUrl)
	}

	return websiteUrls, nil
}

func (db *DB) ListPages(websiteUrl string, page int, pageSize int) ([]Page, error) {
	var pages []Page
	offset := (page - 1) * pageSize
	result := db.conn.Where("URL = ?", websiteUrl).Offset(offset).Limit(pageSize).Find(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return pages, nil
}

func (db *DB) InsertLink(link Link) error {
	result := db.conn.Create(&link)
	return result.Error
}
