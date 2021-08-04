package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func connectDB() *gorm.DB {
	// DB接続
	db, err := sqlConnect()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("DBへの接続に成功しました")
	}

	return db
}

func setNewCOnference(db *gorm.DB) {
	// 会議追加
	err := _setNewCOnference(db)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("新しい会議を追加できました")
	}
}

func findConferences(db *gorm.DB) []*Conferences {
	// SELECT * FROM conferences
	result := []*Conferences{}
	error := db.Find(&result).Error
	if error == nil && len(result) != 0 {

		// for _, conference := range result {
		// 	fmt.Printf("Conference ID %d is uploaded at %s\n", conference.ConferenceID, conference.UploadAt)
		// }
		return result
	}
}

// SQLConnect DB接続
// 接続情報は個人で試した時のものです
func sqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "root"
	PASS := "gordandsql"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "server_database"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"
	return gorm.Open(DBMS, CONNECT)
}

// 新規インスタンス生成
// 被りの無いConferenceIDが自動で生成される想定． (ConferenceIDがprimary_key制約とauto_increment制約を持つ)
// StartAtとEndAtの初期値はNULL
func _setNewCOnference(db *gorm.DB) error {
	err := db.Create(&Conferences{
		StartAt:  nil,
		EndAt:    nil,
		UploadAt: getDate(),
	}).Error
	return err
}

// 現在時刻取得
func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}

// Conferences ユーザー情報のテーブル情報
// 0列目：ConferenceID
// 1列目：StartAt
// 2列目：EndAt
// 3列目：UploadAt
type Conferences struct {
	ConferenceID int
	StartAt      *string `json:"startAt" sql:"type:date"`
	EndAt        *string `json:"endAt" sql:"type:date"`
	UploadAt     string  `json:"uploadAt" sql:"not null;type:date"`
}
