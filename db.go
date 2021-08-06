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

func setNewConferenceID(db *gorm.DB) (id int, err error) {
	// 新規インスタンス生成
	// 従来版と異なり，会議IDを返すことができます(2021/08/05 サブ活終了後に追加)．またgetDate()が不要です．
	// 被りの無いConferenceIDが自動で生成される想定． (ConferenceIDがprimary_key制約とauto_increment制約を持つ)
	// StartAtとEndAtの初期値はNULL
	conference := &Conferences{
		StartAt:  nil,
		EndAt:    nil,
		UploadAt: getDate(),
	}
	err = db.Create(conference).Error
	return conference.ConferenceID, err
}

func settingEnd(db *gorm.DB, id int, end string) error {
	// 会議IDに対応する終了時間を設定する．(2021/08/05 サブ活終了後に追加)
	var conference Conferences
	err := db.Model(&conference).Where("conference_id = ?", id).Update("end_at", end).Error
	return err
}

func settingStart(db *gorm.DB, id int, start string) error {
	// 会議IDに対応する開始時間を設定する．(2021/08/05 サブ活終了後に追加)
	var conference Conferences
	err := db.Model(&conference).Where("conference_id = ?", id).Update("start_at", start).Error
	return err
}

func findParticularConference(db *gorm.DB, id int) Conferences {
	// 会議IDに対応する構造体を返す．(2021/08/05 サブ活終了後に追加)
	var result Conferences
	error := db.Where("conference_id = ?", id).Find(&result).Error
	if error != nil {
		panic(error.Error())
	}
	return result
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

	return nil
}

// SQLConnect DB接続
// 接続情報は庄司さんが構築してくださったデータベースにしました(2021/08/05 サブ活終了後に変更)．
func sqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "b2be3f4d5c559b"
	PASS := "ae8dccb5"
	PROTOCOL := "tcp(us-cdbr-east-04.cleardb.com)"
	DBNAME := "heroku_2b2939979afb8ce"

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
func getDate() int64 {
	now := time.Now()
	return now.Unix()
}

// Conferences 会議情報のテーブル情報
type Conferences struct {
	ConferenceID int     `gorm:"primary_key"`
	StartAt      *int64  `json:"startAt"`
	EndAt        *int64  `json:"endAt"`
	UploadAt     int64   `json:"uploadAt""`
	PresenterNum int     `json:"presenterNum"`
	PTime        float64 `json:"pTime"`
	BTime        float64 `json:"bTime"`
}

// Conferences プレゼンター毎の発表情報のテーブル情報
type Presentations struct {
	ConferenceID int
	Number       int     `json:"number"`
	Presenter    string  `json:"presenter"`
	Time         float64 `json:"time"`
}
