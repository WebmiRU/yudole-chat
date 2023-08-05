package twitch

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yudole-chat/messages"
)

var Out = make(chan any, 9999)

var re = regexp.MustCompile(`^(?:@([^\r\n ]*) +|())(?::([^\r\n ]+) +|())([^\r\n ]+)(?: +([^:\r\n ]+[^\r\n ]*(?: +[^:\r\n ]+[^\r\n ]*)*)|())?(?: +:([^\r\n]*)| +())?[\r\n]*$`)
var socket net.Conn

func Connect() {
	host := os.Getenv("TWITCH_HOST")
	port := os.Getenv("TWITCH_PORT")
	login := os.Getenv("TWITCH_LOGIN")
	password := os.Getenv("TWITCH_PASSWORD")
	channel := os.Getenv("TWITCH_CHANNEL")
	var err error

	log.Println("Connecting to Twitch")

	if socket, err = net.Dial("tcp", fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Println("IRC CONNECTION ERROR: ", err)
		// @TODO Возможно тут стоить сделать задержку и переподключение
		//log.Printf("Reconnecting to Twitch after %d seconds", reconnectionDelay)
	}

	if err != nil {
		reconnect()
		return
	}

	//defer socket.Close()
	socket.SetReadDeadline(time.Now().Add(time.Second * 20))

	fmt.Fprintln(socket, "CAP REQ :twitch.tv/commands twitch.tv/tags twitch.tv/membership")
	fmt.Fprintln(socket, fmt.Sprintf("PASS %s", password))
	fmt.Fprintln(socket, fmt.Sprintf("NICK %s", login))

	var isPingSend bool

	for {
		scanner := bufio.NewScanner(bufio.NewReader(socket))
		socket.SetReadDeadline(time.Now().Add(time.Second * 20))
		//scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			isPingSend = false
			socket.SetReadDeadline(time.Now().Add(time.Second * 20))
			msg := message(scanner.Text())

			switch strings.ToLower(msg.Type) {
			case "001": // Welcome message
				fmt.Fprintln(socket, fmt.Sprintf("JOIN #%s", channel))
				break

			case "join":
				if strings.EqualFold(msg.Login, channel) {
					Out <- messages.System{
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
					Service: "twitch",
					Type:    "channel/message",
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

			case "pong":
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

		if err == nil && !isPingSend {
			isPingSend = true
			fmt.Fprintln(socket, "PING :tmi.twitch.tv")
			continue
		}

		socket.Close()
		reconnect()

		return
	}
}

func reconnect() {
	log.Println("Service TWITCH connection is broken, reconnect after 5 seconds")
	time.Sleep(time.Second * 5)
	go Connect()
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

	fmt.Println("IRC MESSAGE:", msg)

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
	msg := []rune(message.Text)
	offset := 0

	if _, ok := message.Tags["emotes"]; !ok {
		return message.Text
	}

	if len(message.Tags["emotes"]) == 0 {
		return message.Text
	}

	for _, smile := range strings.Split(message.Tags["emotes"], "/") {
		smileIdFromTo := strings.Split(smile, ":")
		smileId := smileIdFromTo[0]

		for _, fromTo := range strings.Split(smileIdFromTo[1], ",") {
			smileFromTo := strings.Split(fromTo, "-")
			smileFrom, _ := strconv.Atoi(smileFromTo[0])
			smileTo, _ := strconv.Atoi(smileFromTo[1])
			smileText := msg[smileFrom+offset : smileTo+offset+1]
			smileReplacer := []rune(fmt.Sprintf("<img class=\"smile twitch\" src=\"https://static-cdn.jtvnw.net/emoticons/v2/%s/default/dark/1.0\" alt=\"%s\"/>", smileId, string(smileText)))
			msg = append(msg[:smileFrom+offset], append(smileReplacer, msg[smileTo+1+offset:]...)...)
			offset += smileFrom - smileTo + len(smileReplacer) - 1
		}

	}

	return string(msg)
}
