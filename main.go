// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// var addr = flag.String("addr", ":8080", "http service address")

type Memo struct {
	Messagetype string `json:"messagetype"`
	Message     string `json:"message"`
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
	fmt.Println("Start main func.")
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("$PORT must be set")
	}

	hub := newHub()
	// startEcho()
	go hub.run() // hubのゴルーチン開始

	fmt.Println("Start echo.")
	e := echo.New()

	initRouting(e, hub)

	fmt.Println("End main func.")
	// e.Logger.Fatal(e.Start(":1323"))
	e.Logger.Fatal(e.Start(":" + port))

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
