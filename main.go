// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// var addr = flag.String("addr", ":8080", "http service address")

// type Memo struct {
// 	Messagetype string `json:"messagetype"`
// 	MeetingID   int    `json:"meetingid"`
// 	Message     string `json:"message"`
// }
type Memo struct {
	Messagetype string `json:"messagetype"`
	Message     string `json:"message"`
}

type Setting struct {
	Messagetype   string   `json:"messagetype"`
	Presenterlist []string `json:"presenterlist"`
	TimeSetting   []int    `json:"timesetting"`
	Starttime     int      `json:"starttime"`
	Endtime       int      `json:"endtime"`
	Presentime    int      `json:"presentime"`
	Breaktime     int      `json:"breaktime"`
}

type ChangePresenter struct {
	Messagetype   string `json:"messagetype"`
	Nextpresenter int    `json:"nowpresenter"`
	TimeSetting   []int  `json:"timesetting"`
}

//　webページに移動
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	hub := newHub()
	// startEcho()
	go hub.run() // hubのゴルーチン開始

	e := echo.New()

	initRouting(e, hub)

	e.Logger.Fatal(e.Start(":1323"))

	// http.HandleFunc("/", serveHome) // TOP画面の表示周り(それ以外はNot Found)
	// // websockerの扱い(直接アクセスはBad Request)
	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	serveWs(hub, w, r)
	// })
	// err := http.ListenAndServe(*addr, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}
