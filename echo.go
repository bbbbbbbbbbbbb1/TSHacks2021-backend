package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func printConferences(db *gorm.DB) []*Conferences {
	conferences := findConferences(db)

	// return c.JSON(http.StatusOK, conferences_json)
	return conferences
}

func initRouting(e *echo.Echo, hub *Hub, db *gorm.DB) {

	e.GET("/", func(c echo.Context) error {
		// return c.String(http.StatusOK, "Hello, World!")
		serveHome(c.Response(), c.Request())
		// return c.JSON(http.StatusOK, {"ok": true})
		return nil
	})

	e.GET("/ip", func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(("<h3>あなたのIPアドレスは %s</h3>"), c.RealIP()))
	})

	e.GET("/users/:id", func(c echo.Context) error {
		jsonMap := map[string]string{
			"name": "okutani",
			"hoge": "piyo",
		}
		return c.JSON(http.StatusOK, jsonMap)
	})

	//e.GET("/meeting/create", createMeeting)
	e.GET("/meeting/create", func(c echo.Context) error {
		// 設定情報を受け取る
		// URLからパラメータ取得

		// 設定情報（仮）
		ptime := 0.0
		btime := 0.0
		num := 0
		names := []string{"Golang", "Java"}

		// DBに会議追加
		conferenceID, _ := setNewConferenceID(db, ptime, btime, names, num)
		// 会議IDを返す
		//return c.JSON(http.StatusOK, conferenceID)

		eventsEx := []Conferences{}
		// 指定した複数の条件を元に複数のレコードを引っ張ってくる
		db.Find(&eventsEx, "conference_id=?", conferenceID)
		// 会議情報を返す
		return c.JSON(http.StatusOK, eventsEx)

	})

	e.GET("/conferences", func(c echo.Context) error {
		result := printConferences(db)
		return c.JSON(http.StatusOK, result)
	})

	e.GET("/presentation", func(c echo.Context) error {
		result := findPresentations(db)
		return c.JSON(http.StatusOK, result)
	})

	e.GET("/setting/:id", func(c echo.Context) error {
		type InitSetting struct {
			Presenterlist []string `json:"presenterlist"`
			TimeSetting   []int    `json:"timesetting"`
			Starttime     int      `json:"starttime"`
			Endtime       int      `json:"endtime"`
			Presentime    int      `json:"presenttime"`
			Breaktime     int      `json:"breaktime"`
			PresenterNum  int      `json:"presenternum"`
		}

		// id := 175
		id, _ := strconv.Atoi(c.Param("id"))
		result := findParticularConference(db, id)
		starttime := int(result.StartAt)
		endtime := int(result.EndAt)
		presentime := int(result.PTime)
		breaktime := int(result.BTime)
		presenter_num := result.PresenterNum
		presenter_list := findParticularPresenters(db, id)
		time_list := timelist(presenter_list, len(presenter_list), presentime, breaktime)
		message := InitSetting{presenter_list, time_list, starttime, endtime, presentime, breaktime, presenter_num}

		// messagejson, _ := json.Marshal(message)

		return c.JSON(http.StatusOK, message)
	})

	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})

	// e.Logger.Fatal(e.Start(":1323"))
}
