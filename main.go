package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var ws_client []websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  10000,
	WriteBufferSize: 10000,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/chat/streamer", echo)
	http.HandleFunc("/chat/stream", echo)
	http.HandleFunc("/", home)

	log.Fatal(http.ListenAndServe("0.0.0.0:5367", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	for {
		_, message, err := ws.ReadMessage()

		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)

		err = ws.WriteMessage(1, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

	ws.Close()
}

func home(w http.ResponseWriter, r *http.Request) {
	//homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}
