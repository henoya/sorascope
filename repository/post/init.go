package post

import (
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) (err error) {
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
