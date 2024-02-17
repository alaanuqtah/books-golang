package models

import "gorm.io/gorm"

type Book struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:author`
	Title     *string `json:title`
	Publisher *string `json:publisher`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{}) //creates the database for us bec you need to manually create it for postgres
	return err
}
