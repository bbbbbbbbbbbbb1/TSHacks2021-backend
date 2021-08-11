// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub // 親となるHub

	// The websocket connection.
	conn *websocket.Conn // 自分のwebsocket

	// Buffered channel of outbound messages.
	send chan []byte // broadcastのメッセージを受け取るチャネル
}

type Memo struct {
	Messagetype string `json:"messagetype"`
	Message     string `json:"message"`
	Presenter   int    `json:"presenter"`
	Sender      string `json:"sender"`
}

type Setting struct {
	Messagetype   string   `json:"messagetype"`
	Presenterlist []string `json:"presenterlist"`
	TimeSetting   []int    `json:"timesetting"`
	Starttime     int      `json:"starttime"`
	Endtime       int      `json:"endtime"`
	Presentime    int      `json:"presenttime"`
	Breaktime     int      `json:"breaktime"`
}

type ChangePresenter struct {
	Messagetype   string `json:"messagetype"`
	Nextpresenter int    `json:"nextpresenter"`
	TimeSetting   []int  `json:"timesetting"`
}

func loadJson(byteArray []byte) (interface{}, error) {
	var jsonObj interface{}
	err := json.Unmarshal(byteArray, &jsonObj)
	return jsonObj, err
}

func presenlist(name_list []interface{}) ([]string, int) {
	user_count := len(name_list)
	//fmt.Printf("%T", user_count)
	//fmt.Println(user_count)
	presenter_list := make([]string, user_count)
	for i := 0; i < user_count; i++ {
		presenter_list[i] = name_list[i].(string)
	}
	return presenter_list, user_count
}

func timelist(name_list []interface{}, user_count int, presen_time int, break_time int) []int {
	time_list := make([]int, user_count)
	//開始時間と終了時間を送る

	//timesettingの配列
	for i := 0; i < user_count; i++ {
		if name_list[i].(string) != "break" {
			time_list[i] = presen_time
		} else {
			time_list[i] = break_time
		}
	}
	return time_list
}

func modify(time_list []interface{}, nextpresenter float64) (int, []int) {
	next_presenter := int(nextpresenter)

	user_count := len(time_list)
	//残りの休憩回数
	break_count := 0
	//残りのpresenter
	left_presenter := 0

	time_setting := make([]int, user_count)
	//とりあえず格納する、休憩回数のカウント
	for i := 0; i < user_count; i++ {
		time_setting[i] = int(time_list[i].(float64))
		if i >= int(nextpresenter) && time_list[i].(float64) == 10 {
			break_count += 1
		} else if i >= int(nextpresenter) && time_list[i].(float64) != 10 {
			left_presenter += 1
		}
	}
	//fmt.Println(break_count)

	//発表が終わったところまでの合計時間
	var finish_time int
	for i := 0; i < int(nextpresenter); i++ {
		finish_time = finish_time + int(time_list[i].(float64))
	}

	//meetingの時間を変更しない場合の合計時間
	var time_sum int
	for i := 0; i < user_count; i++ {
		time_sum = time_sum + int(time_list[i].(float64))
	}
	//fmt.Println(time_sum)

	//開始時刻と終了時刻、発表者の順番をDBからもらう
	//設定の発表時間と休憩時間をもらう

	var meeting_time int
	meeting_time = 150

	//残りの一人あたりの発表時間
	var left_presen_person int

	meetingtime_left := meeting_time - time_sum
	if meetingtime_left < 0 {
		//残り時間
		left_time := meeting_time - finish_time
		break_time_left := 10 * break_count
		left_presen_time := left_time - int(break_time_left)
		left_presen_person = left_presen_time / int(left_presenter)
	}

	for j := next_presenter; j < int(user_count); j++ {
		if time_setting[j] != 10 {
			time_setting[j] = left_presen_person
		}
	}

	return next_presenter, time_setting
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()

		// messagestruct := struct{}{}

		jsonObj, jsonerr := loadJson(message)
		if jsonerr != nil {
			continue
		}
		fmt.Printf(string(message) + "\n")
		message_type := jsonObj.(map[string]interface{})["messagetype"].(string)

		// エラー処理
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// messagestruct := Memo{"memo", int(1), string(message)}
		// messagejson, _ := json.Marshal(messagestruct)

		var messagestruct interface{}

		if message_type == "memo" {
			message_jsonobj := jsonObj.(map[string]interface{})["message"].(string)
			presenter := jsonObj.(map[string]interface{})["presenter"].(float64)
			presenter_id := int(presenter)
			sender := jsonObj.(map[string]interface{})["sender"].(string)
			//MeetingID := jsonObj.(map[string]interface{})["MeetingID"].(int)
			//messagestruct = Memo{"memo", MeetingID, message}
			messagestruct = Memo{"memo", message_jsonobj, presenter_id, sender}
			//messagejson, _ := json.Marshal(messagestruct)
		} else if message_type == "setting" {
			name_list := (jsonObj.(map[string]interface{})["presenterlist"]).([]interface{})

			presenter_list, user_count := presenlist(name_list)
			//DBにpresenterのlistを送る

			starttime := (jsonObj.(map[string]interface{})["starttime"]).(float64)
			endtime := (jsonObj.(map[string]interface{})["endtime"]).(float64)
			start_time := int(starttime)
			end_time := int(endtime)

			presentime := (jsonObj.(map[string]interface{})["presentime"]).(float64)
			breaktime := (jsonObj.(map[string]interface{})["breaktime"]).(float64)
			presen_time := int(presentime)
			break_time := int(breaktime)

			time_list := timelist(name_list, user_count, presen_time, break_time)

			messagestruct = Setting{"setting", presenter_list, time_list, start_time, end_time, presen_time, break_time}
			//messagejson, _ := json.Marshal(messagestruct)
			//fmt.Println(string(messagejson))
		} else if message_type == "change" {
			nextpresenter := (jsonObj.(map[string]interface{})["nextpresenter"]).(float64)
			time_list := (jsonObj.(map[string]interface{})["timesetting"]).([]interface{})

			next_presenter, time_setting := modify(time_list, nextpresenter)

			//fmt.Println(left_presen_person)
			messagestruct = ChangePresenter{"change", next_presenter, time_setting}
			//messagejson, _ := json.Marshal(messagestruct)

			//fmt.Println(string(messagejson))
		} else {
			return
		}

		messagejson, _ := json.Marshal(messagestruct)

		// 自分のメッセージをhubのbroadcastチャネルに送り込む
		fmt.Printf("%+v\n", messagestruct)
		c.hub.broadcast <- messagejson
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			// タイムアウト時間の設定
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// エラー処理
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("unsuccessed upgrade.")
		log.Println(err)
		return
	} else {
		fmt.Println("successed upgrade!")
	}
	// sendは他の人からのメッセージが投入される
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client // hubのregisterチャネルに自分のClientを登録

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
