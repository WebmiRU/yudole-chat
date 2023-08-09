package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"slices"
	"yudole-chat/messages"
)

var wsClients []*websocket.Conn // Все клиенты

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

	wsClients = append(wsClients, ws)

	for {
		var message messages.Channel
		err := ws.ReadJSON(&message)

		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("RECIVED: %s", message)

		if err = ws.WriteJSON(message); err != nil {
			log.Println(err)
			break
		}
	}

	if idx := slices.Index(wsClients, ws); idx >= 0 {
		wsClients = slices.Delete(wsClients, idx, idx+1)
	}

	if err = ws.Close(); err != nil {
		log.Println(err)
	}
}

func chat(w http.ResponseWriter, req *http.Request) {
	dir, _ := os.Getwd()
	theme := req.URL.Query().Get("theme")

	_, err := os.Stat(fmt.Sprintf("%s/public/themes/%s/index.html", dir, theme))
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Theme \"%s\" not found. Using system theme", theme)

		Out <- messages.System{
			Type:    "error/theme/notfound",
			Service: "system",
			Value:   theme,
		}

		theme = "system"
	}

	buf, _ := os.ReadFile(fmt.Sprintf("%s/public/themes/%s/index.html", dir, theme))

	w.Write(buf)
}
