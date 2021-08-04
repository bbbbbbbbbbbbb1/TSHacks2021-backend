package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func createMeeting(c echo.Context) error {
	// DBに会議追加
	// 会議情報を返す
	return nil
}

func initRouting(e *echo.Echo, hub *Hub) {

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

	e.GET("/meeting/create", createMeeting)

	e.GET("/ws", func(c echo.Context) error {
		serveWs(hub, c.Response(), c.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":1323"))
}
