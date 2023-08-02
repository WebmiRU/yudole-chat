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

	http.HandleFunc("/chat", acceptStreamer)      // Чат стримера, отображает все сообщения
	http.HandleFunc("/chat/stream", acceptStream) // Чат для стрима, отображает только общие сообщения
	http.Handle("/", http.FileServer(http.Dir("./public")))

	go twitch.Connect()
	go goodgame.Connect()
	go trovo.Connect()

	// Чтение общих сообщений
	go func() {
		for {
			select {
			// Twitch
			case message := <-twitch.OutAll:
				fmt.Println("MESSAGE:", message)

				for len(wsClientsAll) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, ws := range wsClientsAll {
					ws.WriteJSON(message)
				}
				break

			// GoodGame
			case message := <-goodgame.OutAll:
				fmt.Println("MESSAGE:", message)

				for len(wsClientsAll) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClientsAll {
					wsClient.WriteJSON(message)
				}
				break

			// Trovo
			case message := <-trovo.OutAll:
				fmt.Println("MESSAGE:", message)

				for len(wsClientsAll) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClientsAll {
					wsClient.WriteJSON(message)
				}
				break
			}
		}
	}()

	// Чтение сообщений на стримерском канале
	go func() {
		for {
			select {
			// Twitch
			case system := <-twitch.OutStreamer:
				fmt.Println("SYSTEM MESSAGE:", system)

				for len(wsClientsStreamer) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClientsStreamer {
					wsClient.WriteJSON(system)
				}
				break

			// GoodGame
			case system := <-goodgame.OutStreamer:
				fmt.Println("SYSTEM MESSAGE:", system)

				for len(wsClientsStreamer) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClientsStreamer {
					wsClient.WriteJSON(system)
				}
				break

			// Trovo
			case system := <-trovo.OutStreamer:
				fmt.Println("SYSTEM MESSAGE:", system)

				for len(wsClientsStreamer) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				for _, wsClient := range wsClientsStreamer {
					wsClient.WriteJSON(system)
				}
				break
			}
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5367", nil)) // Websocket main server
}
