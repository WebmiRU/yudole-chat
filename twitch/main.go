package twitch

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"yudole-chat/messages"
)

var Out = make(chan messages.Channel, 999)
var OutSystem = make(chan messages.System, 999)

var re = regexp.MustCompile(`^(?:@([^\r\n ]*) +|())(?::([^\r\n ]+) +|())([^\r\n ]+)(?: +([^:\r\n ]+[^\r\n ]*(?: +[^:\r\n ]+[^\r\n ]*)*)|())?(?: +:([^\r\n]*)| +())?[\r\n]*$`)
var socket net.Conn

func Connect() {
	host := os.Getenv("TWITCH_HOST")
	port := os.Getenv("TWITCH_PORT")
	login := os.Getenv("TWITCH_LOGIN")
	password := os.Getenv("TWITCH_PASSWORD")
	channel := os.Getenv("TWITCH_CHANNEL")
	var err error

	socket, err = net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))

	if err != nil {
		log.Fatal("IRC CONNECTION ERROR: ", err)
	}

	fmt.Fprintln(socket, "CAP REQ :twitch.tv/commands twitch.tv/tags twitch.tv/membership")
	fmt.Fprintln(socket, fmt.Sprintf("PASS %s", password))
	fmt.Fprintln(socket, fmt.Sprintf("NICK %s", login))

	for {
		scanner := bufio.NewScanner(bufio.NewReader(socket))
		//scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			msg := message(scanner.Text())

			switch strings.ToLower(msg.Type) {
			case "001": // Welcome message
				fmt.Fprintln(socket, fmt.Sprintf("JOIN #%s", channel))
				break

			case "join":
				if strings.EqualFold(msg.Channel, channel) {
					OutSystem <- messages.System{
						Service: "twitch",
						Type:    "channel/join/success",
						Text:    fmt.Sprintf("Успешное подключение к каналу %s", msg.Channel),
					}
				}
				break

			case "part":
				break

			case "privmsg":
				Out <- messages.Channel{
					Type: "channel/message",
					User: messages.User{
						Login:     msg.Login,
						Nick:      msg.Nick,
						AvatarUrl: "",
						Color:     "",
					},
					Message: messages.Message{
						Text: msg.Text,
						Html: smiles(msg),
					},
				}
				break

			case "ping":
				_, err := fmt.Fprintln(socket, "PONG :"+msg.Text)
				if err != nil {
					log.Fatal("IRC CONNECTION ERROR", err)
				}
				break

			case "roomstate":
				break

			case "userstate":
				break

			case "globaluserstate":
				break

			case "002", "003", "004", "353", "366", "372", "375", "376", "cap":
				// Ignore this message types
				break

			default:
				log.Println("UNKNOWN IRC MESSAGE TYPE: ", msg.Type)
			}
		}
	}
}

func tags(tags string) map[string]string {
	result := make(map[string]string)

	if len(tags) == 0 {
		return result
	}

	for _, v := range strings.Split(tags, ";") {
		kv := strings.SplitN(v, "=", 2)

		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}

	return result
}

func message(msg string) Message {
	matches := re.FindStringSubmatch(msg)
	tags := tags(matches[1])
	var nick string

	//fmt.Println("IRC MESSAGE:", msg)

	if v, ok := tags["display-name"]; ok {
		nick = v
	}

	var message = Message{
		Login:   strings.Split(matches[3], "!")[0],
		Nick:    nick,
		Type:    matches[5],
		Channel: strings.Replace(matches[6], "#", "", 1),
		Text:    matches[8],
		Tags:    tags,
		Prefix:  matches[3],
	}

	return message
}

func smiles(message Message) string {

	return ""
}