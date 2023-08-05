package trovo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"yudole-chat/messages"
)

const (
	UrlGetUsers        = "https://open-api.trovo.live/openplatform/getusers"
	UrlGetChannelToken = "https://open-api.trovo.live/openplatform/chat/channel-token/"
	UrlTrovoWS         = "wss://open-chat.trovo.live/chat"
)

var Out = make(chan any, 9999)
var regexSmile = regexp.MustCompile(`:\b(\S+)\b`)

func Connect() {
	clientId := os.Getenv("TROVO_CLIENT_ID")
	channel := os.Getenv("TROVO_CHANNEL")

	req, err := http.NewRequest("POST", UrlGetUsers, bytes.NewBuffer([]byte(fmt.Sprintf("{\"user\":[\"%s\"]})", channel))))
	req.Header.Set("Client-ID", clientId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var data MessageResponseUsers
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Fatalln("Error while parsing response from Trovo:", err)
	}

	if data.Total < 1 || len(data.Users) < 1 {
		log.Fatalln("Error response from Trovo while channel search", data)
	}

	channelId := data.Users[0].ChannelId

	req, err = http.NewRequest("GET", UrlGetChannelToken+channelId, nil)
	req.Header.Set("Client-ID", clientId)
	resp, err = client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, err = io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var respToken MessageResponseToken
	err = json.Unmarshal(body, &respToken)

	if err != nil {
		log.Fatalln("Error response from Trovo while getting channel token", string(body))
	}

	wsClient, _, err := websocket.DefaultDialer.Dial(UrlTrovoWS, nil)

	if err != nil {
		log.Fatal("Trovo chat server connection error", err)
	}

	err = wsClient.WriteJSON(MessageAuth{
		Type:  "AUTH",
		Nonce: string(time.Now().Unix()),
		Data: MessageAuthData{
			Token: respToken.Token,
		},
	})

	if err != nil {
		log.Fatal("Trovo chat server send message error", err)
	}

	for {
		var message Message
		if err := wsClient.ReadJSON(&message); err != nil {
			log.Println(err)
		}

		f, _ := os.OpenFile("log/trovo.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

		s, _ := json.Marshal(message)
		f.Write(s)
		f.WriteString("\n")
		f.Close()

		fmt.Println(message)

		switch strings.ToLower(message.Type) {
		case "response":
			Out <- messages.System{
				Service: "trovo",
				Type:    "channel/join/success",
				Text:    fmt.Sprintf("Успешное подключение к каналу %s", channel),
			}
			log.Println("SUCCESS JOIN (TROVO)")

			go func() {
				for {
					time.Sleep(20 * time.Second)

					err := wsClient.WriteJSON(MessagePing{
						Type:  "PING",
						Nonce: string(time.Now().Unix()),
					})

					if err != nil {
						log.Fatalln("Error Trovo send PING command", err)
					}
				}
			}()
			break

		case "chat":
			for _, chat := range message.Data.Chats {
				if chat.Type == 5007 {
					continue
				}

				Out <- messages.Channel{
					Service: "trovo",
					Type:    "channel/message",
					User: messages.User{
						Login:     chat.UserName,
						Nick:      chat.UserName,
						AvatarUrl: "",
						Color:     "",
					},
					Message: messages.Message{
						Text: chat.Content,
						Html: smile(chat.Content),
					},
				}
			}
			break

		case "pong":
			fmt.Println("Trovo PONG message")
		}
	}
}

func smile(message string) string {
	return regexSmile.ReplaceAllString(message, `<img src="https://img.trovo.live/emotes/$1.png?imageView2/1/w/72/h/72/format/webp&max_age=31536000" alt="$1"/>`)
}
