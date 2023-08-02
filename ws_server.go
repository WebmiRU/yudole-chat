package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"slices"
	"yudole-chat/messages"
)

var wsClientsAll []*websocket.Conn      // Все клиенты
var wsClientsStreamer []*websocket.Conn // Все клиенты

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10000,
	WriteBufferSize: 10000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func acceptStreamer(w http.ResponseWriter, r *http.Request) {
	accept(w, r, true)
}

func acceptStream(w http.ResponseWriter, r *http.Request) {
	accept(w, r, false)
}

func accept(w http.ResponseWriter, r *http.Request, isStreamer bool) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	wsClientsAll = append(wsClientsAll, ws)

	if isStreamer {
		wsClientsStreamer = append(wsClientsStreamer, ws)
	}

	for {
		var message messages.Channel
		err := ws.ReadJSON(&message)

		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("RECIVED: %s", message)

		err = ws.WriteJSON(message)

		if err != nil {
			log.Println(err)
			break
		}
	}

	if idx := slices.Index(wsClientsAll, ws); idx >= 0 {
		wsClientsAll = slices.Delete(wsClientsAll, idx, idx+1)
	}

	if idx := slices.Index(wsClientsStreamer, ws); idx >= 0 {
		wsClientsAll = slices.Delete(wsClientsAll, idx, idx+1)
	}

	err = ws.Close()
	if err != nil {
		log.Println(err)
	}
}
