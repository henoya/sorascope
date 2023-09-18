package sql

import (
	"fmt"
	"gorm.io/gorm"
)

// sql.go
// sql まわりの操作関数

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
