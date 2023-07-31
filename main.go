package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
	"yudole-chat/twitch"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/chat/streamer", accept)
	http.HandleFunc("/chat/stream", accept)
	http.HandleFunc("/chat", accept)
	http.HandleFunc("/", home)

	go twitch.Connect()

	go func() {
		for {
			fmt.Println("WS CLIENTS:", ws_clients)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case message := <-twitch.Out:
				fmt.Println("MESSAGE:", message)

				for _, ws := range ws_clients {
					ws.WriteJSON(message)
				}
				break
			case system := <-twitch.OutSystem:
				fmt.Println("SYSTEM MESSAGE:", system)
				for _, ws := range ws_clients {
					ws.WriteJSON(system)
				}
				break
			}
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5367", nil)) // Websocket main server
}
