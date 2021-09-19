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

// 現在の発表者を更新する
// 経過時間→pt, 会議ID→id
func changePresenter(db *gorm.DB, pt float64, id int) error {
	var result1 Conferences
	var result2 Presentations
	var err error
	err1 := db.Where("conference_id = ?", id).Find(&result1).Error
	err2 := db.Model(&result2).Where("conference_id = ? AND number = ?", id, result1.PresenterNum).Update("time", pt).Error
	if err1 == nil && err2 == nil {
		var result3 Conferences
		next := result1.PresenterNum + 1
		err = db.Model(&result3).Where("conference_id = ?", id).Update("presenter_num", next).Error
	} else {
		err = err1
	}
	return err
}

/* 20210919追加・更新　ここから */
// id(会議ID)がデータベースに持つ情報を，引数の情報(発表時間p_time, 休憩時間b_time, 発表者名配列names, 発表者人数num)で全て更新．
// 8/26最終回やり残しまとめ バックエンド 「元の発表者を消して新たに追加しなおす」を満たすことを想定しています．
func resetParticularConference(db *gorm.DB, id int, ptime float64, btime float64, names []string, num int) (err error) {

	// conferencesに現在保持している情報を削除
	db.Debug().Delete(&Conferences{}, "conference_id = ?", id)
	// presentationsに現在保持している情報を削除
	db.Debug().Delete(&Presentations{}, "conference_id = ?", id)

	// 会議インスタンス再作成（setNewConference準拠）
	now := getDate()
	conference := &Conferences{
		ConferenceID: id,
		StartAt:      nil,
		EndAt:        nil,
		UploadAt:     now,
		PresenterNum: -1,
		PTime:        ptime,
		BTime:        btime,
	}
	err = db.Create(conference).Error
	if err == nil {
		for i := 0; i < num; i++ {
			var t float64
			if names[i] == "break" {
				t = conference.BTime
			} else {
				t = conference.PTime
			}
			presentation := &Presentations{
				ConferenceID: conference.ConferenceID,
				Number:       i,
				Presenter:    names[i],
				Time:         t,
			}
			err = db.Create(presentation).Error
		}
	}
	return err
}

// 設定発表時間→ptime, 設定休憩時間→btime, 発表者名スライス→names(休憩は文字列"break"で受け取る), 発表者数→num
// 8/26最終回やり残しまとめ バックエンド 「presenter_numを-1スタートにする」を満たすことを想定しています．
func setNewConferenceID(db *gorm.DB, ptime float64, btime float64, names []string, num int) (id int, err error) {
	now := getDate()
	conference := &Conferences{
		StartAt:      nil,
		EndAt:        nil,
		UploadAt:     now,
		PresenterNum: -1,
		PTime:        ptime,
		BTime:        btime,
	}
	err = db.Create(conference).Error
	if err == nil {
		for i := 0; i < num; i++ {
			var t float64
			if names[i] == "break" {
				t = conference.BTime
			} else {
				t = conference.PTime
			}
			presentation := &Presentations{
				ConferenceID: conference.ConferenceID,
				Number:       i,
				Presenter:    names[i],
				Time:         t,
			}
			err = db.Create(presentation).Error
		}
	}
	return conference.ConferenceID, err
}

/* 20210919追加・更新　ここまで */

func settingEnd(db *gorm.DB, id int, end int64) error {
	var conference Conferences
	err := db.Model(&conference).Where("conference_id = ?", id).Update("end_at", end).Error
	return err
}

func settingStart(db *gorm.DB, id int, start int64) error {
	var conference Conferences
	err := db.Model(&conference).Where("conference_id = ?", id).Update("start_at", start).Error
	return err
}

/* 20210826追加 １　ここから */
// 会議IDに対応する発表者リストを返す．
func findParticularPresenters(db *gorm.DB, id int) []string {
	tmp := []*Presentations{}
	error := db.Where("conference_id = ?", id).Select("presenter").Find(&tmp).Error
	if error != nil {
		panic(error.Error())
	}
	result := []string{}
	for _, presen := range tmp {
		result = append(result, presen.Presenter)
	}
	return result
}

// 会議IDに対応する発表構造体を返す．
func findParticularPresentations(db *gorm.DB, id int) []*Presentations {
	result := []*Presentations{}
	error := db.Where("conference_id = ?", id).Find(&result).Error
	if error != nil {
		panic(error.Error())
	}
	return result
}

// 全発表の全情報の構造体配列を返す
func findPresentations(db *gorm.DB) []*Presentations {
	// SELECT * FROM conferences
	result := []*Presentations{}
	error := db.Find(&result).Error
	if error == nil && len(result) != 0 {
		return result
	}
	return nil
}

/* 20210826追加　ここまで */

/* 20210919追加　ここから */
// メモの取得
func findParticularMemos(db *gorm.DB, id int) []*Memos {
	result := []*Memos{}
	error := db.Where("conference_id = ?", id).Find(&result).Error
	if error != nil {
		panic(error.Error())
	}
	return result
}

// メモの追加
func settingMemo(db *gorm.DB, id int, writer string, content string) error {
	memo := &Memos{
		ConferenceID: id,
		Writer:       writer,
		Content:      content,
	}
	err := db.Create(memo).Error
	return err
}

/* 20210919追加　ここまで */

func findParticularConference(db *gorm.DB, id int) Conferences {
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

/* 20210919追加　ここから */
// Memos  メモのテーブル情報
type Memos struct {
	ConferenceID int
	Writer       string `json:"writer"`
	Content      string `json:"content"`
}

/* 20210919追加　ここまで */
