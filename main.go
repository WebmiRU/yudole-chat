package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"yudole-chat/twitch"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/chat", accept)
	http.HandleFunc("/", home)

	go twitch.Connect()

	go func() {
		for {
			select {
			case message := <-twitch.Out:
				fmt.Println("MESSAGE:", message)

				for len(ws_clients) == 0 {
					continue
				}

				for _, ws := range ws_clients {
					ws.WriteJSON(message)
				}
				break
			case system := <-twitch.OutSystem:
				fmt.Println("SYSTEM MESSAGE:", system)

				for len(ws_clients) == 0 {
					continue
				}

				for _, ws := range ws_clients {
					ws.WriteJSON(system)
				}
				break
			}
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5367", nil)) // Websocket main server
}
