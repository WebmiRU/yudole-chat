package goodgame

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"yudole-chat/messages"
)

var Out = make(chan messages.Channel, 9999)
var OutSystem = make(chan messages.System, 9999)
var smiles = make(map[string]string)

func Connect() {
	data, _ := os.ReadFile("goodgame_smiles.json")
	var smilesData []Smile
	err := json.Unmarshal(data, &smilesData)

	if err != nil {
		log.Fatal(err)
	}

	for _, smile := range smilesData {
		url := smile.Images.Big

		if len(smile.Images.Gif) > 0 {
			url = smile.Images.Gif
		}

		smiles[":"+smile.Key+":"] = url
	}

	url := os.Getenv("GOODGAME_CHAT_URL")
	channelUrl := os.Getenv("GOODGAME_CHANNEL_URL")

	resp, err := http.Get(channelUrl)

	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(resp.Body)

	r := regexp.MustCompile(`"channel_id"\s?:\s?"(\d+)"`).FindSubmatch(body)
	//channelId, err := strconv.Atoi(string(r[1]))
	//if err != nil {
	//	log.Fatal("Error while getting GoodGame ChannelID value")
	//}

	if len(r) < 2 {
		log.Fatal("Error while getting GoodGame ChannelID value")
	}

	channelId := string(r[1])

	client, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		log.Fatal("Goodgame chat server connection error")
	}

	for {
		var message MessageIncome
		err := client.ReadJSON(&message)

		if err != nil {
			log.Println(err)
		}

		log.Println(message)

		switch strings.ToLower(message.Type) {
		case "welcome":
			msgAuth := MessageAuth{
				Type: "auth",
				Data: MessageAuthData{
					SiteId: 1,
					UserId: 0,
					Token:  "",
				},
			}

			client.WriteJSON(msgAuth)
			break

		case "success_auth":
			msgJoin := MessageJoin{
				Type: "join",
				Data: MessageJoinData{
					ChannelId: channelId,
					Hidden:    false,
					Mobile:    false,
				},
			}

			client.WriteJSON(msgJoin)
			break

		case "success_join":
			OutSystem <- messages.System{
				Service: "goodgame",
				Type:    "channel/join/success",
				Text:    fmt.Sprintf("Успешное подключение к каналу %s", message.Data.ChannelId),
			}
			log.Println("SUCCESS JOIN")
			break

		case "private_message":
			// TODO Добавить обработку приватных сообщений
			break

		case "message":
			Out <- messages.Channel{
				Service: "goodgame",
				Type:    "channel/message",
				User: messages.User{
					Login:     message.Data.UserName,
					Nick:      message.Data.UserName,
					AvatarUrl: "",
					Color:     "",
				},
				Message: messages.Message{
					Text: message.Data.Text,
					Html: smile(message.Data.Text),
				},
			}
			break
		}
	}
}

func smile(text string) string {
	for k, v := range smiles {
		text = strings.Replace(text, k, fmt.Sprintf("<img src=\"%s\" alt=\"%s\"/>", v, v), -1)
	}

	return text
}
