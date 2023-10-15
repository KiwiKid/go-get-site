package main

import (
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
	DateCreated time.Time `gorm:"type:timestamp"`
	DateUpdated time.Time `gorm:"type:timestamp"`
}

func (Page) TableName() string {
	return "pgml.page"
}

type DB struct {
	conn *gorm.DB
}

// NewDB creates a new DB instance with GORM
func NewDB(connStr string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate will ONLY create tables, missing columns and missing indexes
	db.AutoMigrate(&Page{})

	return &DB{conn: db}, nil
}

// InsertPage inserts a new page record into the database using GORM
func (db *DB) InsertPage(page Page) error {
	result := db.conn.Create(&page)
	return result.Error
}
