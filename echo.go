package main

import (
	"fmt"
	"net/http"

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
		// DBに会議追加
		conferenceID, _ := setNewConferenceID(db)
		// 会議IDを返す
		//return c.JSON(http.StatusOK, conferenceID)

		eventsEx := []Conferences{}
		// 指定した複数の条件を元に複数のレコードを引っ張ってくる
		db.Find(&eventsEx, "conference_id=?", conferenceID)
		// 会議情報を返す
		return c.JSON(http.StatusOK, eventsEx)

	})

	result := printConferences(db)
	e.GET("/conferences", func(c echo.Context) error {
		return c.JSON(http.StatusOK, result)
	})
	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":1323"))
}
