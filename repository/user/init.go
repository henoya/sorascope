package user

import (
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) (err error) {
	// Sqlite3 DB の テーブルを struct から 作成 or マイグレートする
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&UserHandle{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&UserProfile{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&UserSession{}); err != nil {
		panic(err)
	}
	return nil
}
