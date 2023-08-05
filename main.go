package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
	"yudole-chat/goodgame"
	"yudole-chat/trovo"
	"yudole-chat/twitch"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/chat", accept) // Чат стримера, отображает все сообщения
	http.Handle("/", http.FileServer(http.Dir("./public")))

	go twitch.Connect()
	//go twitch.Ping()
	//go goodgame.Connect()
	//go trovo.Connect()

	// Чтение общих сообщений
	go func() {
		for {
			select {
			// Twitch
			case message := <-twitch.Out:
				fmt.Println("MESSAGE:", message)

				for len(wsClients) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, ws := range wsClients {
					ws.WriteJSON(message)
				}
				break

			// GoodGame
			case message := <-goodgame.Out:
				fmt.Println("MESSAGE:", message)

				for len(wsClients) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClients {
					wsClient.WriteJSON(message)
				}
				break

			// Trovo
			case message := <-trovo.Out:
				fmt.Println("MESSAGE:", message)

				for len(wsClients) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClients {
					wsClient.WriteJSON(message)
				}
				break
			}
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5367", nil)) // Websocket main server
}
