package main

import (
	"fmt"
	"github.com/henoya/sorascope/repository/post"
	"github.com/henoya/sorascope/repository/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDBConnection() (db *gorm.DB, err error) {
	// DBファイルのオープン
	db, err = openDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect database")
	}

	err = post.MigrateDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database")
	}
	err = user.MigrateDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database")
	}
	return db, nil
}

func openDB() (db *gorm.DB, err error) {
	// DBファイルのオープン
	db, err = gorm.Open(sqlite.Open("sorascope.db"), &gorm.Config{})
	return db, err
}
