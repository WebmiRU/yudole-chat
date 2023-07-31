package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"slices"
	"yudole-chat/messages"
)

var ws_clients []*websocket.Conn // Все клиенты

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10000,
	WriteBufferSize: 10000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func accept(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	ws_clients = append(ws_clients, ws)

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

	idx := slices.Index(ws_clients, ws)
	slices.Delete(ws_clients, idx, idx+1)

	err = ws.Close()
	if err != nil {
		log.Println(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	//homeTemplate.Execute(w, "ws://"+r.Host+"/accept")
}
