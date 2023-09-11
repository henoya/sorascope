package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// sql.go
// sql まわりの操作関数

func openDB() (db *gorm.DB, err error) {
	// DBファイルのオープン
	db, err = gorm.Open(sqlite.Open("sorascope.db"), &gorm.Config{})
	return db, err
}

func migrateDB(db *gorm.DB) (err error) {
	// Sqlite3 DB の テーブルを struct から 作成 or マイグレートする
	if err := db.AutoMigrate(&Image{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&EmbedImages{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&EmbedExternal{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&AuthorRecord{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&PostHistory{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&PostHistoryStatus{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&PostRecord{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&PostStatus{}); err != nil {
		panic(err)
	}
	return nil
}

// getMaxId 最大のIDを取得する
func getMaxId(db *gorm.DB, table interface{}) (maxId int64, err error) {
	stmt := &gorm.Statement{DB: db}
	_ = stmt.Parse(table)
	fmt.Println(stmt.Schema.Table) // Output: users
	tableName := stmt.Schema.Table
	// 最大のIDを取得
	var count int64
	db.Model(table).Count(&count)
	if count == 0 {
		return 0, nil
	}
	maxId = 0
	err = db.Raw(fmt.Sprintf("SELECT MAX(id) FROM '%s'", tableName)).Scan(&maxId).Error
	return maxId, err
}

// truncateTable 指定テーブルのレコードを全削除する
func truncateTable(db *gorm.DB, tableName interface{}) error {
	db.Unscoped().Where("1 = 1").Delete(tableName)
	return db.Error
}
