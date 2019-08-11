package tpbot

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"gopkg.in/h2non/filetype.v1"
	//"strings"
	"fmt"

	"github.com/donomii/pbot"
	"github.com/donomii/svarmrgo"
	"github.com/donomii/yatzie/shared/registry"
	//"github.com/donomii/yatzie/shared/utils"
	"gopkg.in/tucnak/telebot.v1"
)

var bbs *pbot.BBSdata

type Tpbot struct {
}

func init() {
	plugin_registry.RegisterPlugin(&Tpbot{})
}

var chats map[string]telebot.Chat

func fileType(path string) string {
	buf, _ := ioutil.ReadFile(path)

	kind, unkwown := filetype.Match(buf)
	if unkwown != nil {
		return "application/octet-stream"
	}
	return kind.MIME.Value
}

type OnlyFuckingRetardedProgrammersUseInterfaces struct {
	S string
}

func (s OnlyFuckingRetardedProgrammersUseInterfaces) Destination() string {
	return s.S
}
func handleOneOutgoing(m *pbot.BBSmessage) {
	bot := plugin_registry.Bot
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%s: %s", r, debug.Stack())
			log.Println("Error handling outgoing messages: ", r)

		}
	}()

	//time.Sleep(1 * time.Second)
	log.Printf("Handling message: %+v", m)
	//	var err
	var Chat telebot.Chat
	userSettings := m.UserData.(map[string]string)
	if userSettings == nil {
		log.Println("No user settings found, cannot send message")
		return
	}
	targetUser := userSettings["User"]
	targetChat := userSettings["ChatID"]
	if targetChat == "" {
		log.Println("Could not find chat id, trying to find via username")
		Chat = chats[targetUser]
		targetChat = fmt.Sprintf("%v", Chat.ID)
	}

	if targetChat == "" {
		log.Println("Could not get chat id, unable to send message")
		return
	}

	if m.Message == "text" {
		log.Printf("Sending message... .%v.", m.PayloadString)
		//bot.SendMessage(Chat, m.PayloadString, nil)
		log.Printf("%+v\n", Chat)
		fuckingRetards := OnlyFuckingRetardedProgrammersUseInterfaces{targetChat}
		log.Println(plugin_registry.Bot.SendMessage(fuckingRetards, m.PayloadString, nil))
		log.Println("Done!")
	} else {

		name := "asdfasdf" + m.PayloadString

		log.Print("Sending file... " + name)
		/*
			file := bytes.NewReader(m.PayloadBytes)
			fileLength := int64(len(m.PayloadBytes))
			fileName := m.PayloadString
		*/
		ioutil.WriteFile(name, m.PayloadBytes, 0644)
		log.Printf("File type: '%v'", fileType)
		//fileType := m.Message

		photo, err := telebot.NewFile(name)
		mtype := fileType(name)

		if strings.HasPrefix(mtype, "image") {
			picture := telebot.Photo{File: photo}
			err = bot.SendPhoto(Chat, &picture, &telebot.SendOptions{})
			if err != nil {
				log.Println("Error sending photo")
				log.Println(err)
				bot.SendMessage(Chat, "Could not send photo", nil)
			}

		} else {

			d := telebot.Document{File: photo, FileName: m.PayloadString, Mime: "binary/unknown"}
			err = bot.SendDocument(Chat, &d, nil)
			if err != nil {
				log.Println("Error sending document")
				log.Println(err)
				bot.SendMessage(Chat, "Could not send photo", nil)
			}

		}
	}
}

func handleMessage(m svarmrgo.Message) []svarmrgo.Message {

	out := []svarmrgo.Message{}
	switch m.Selector {
	case "reveal-yourself":
		m.Respond(svarmrgo.Message{Selector: "announce", Arg: "BBS"})
	case "shutdown":
		os.Exit(0)
	case "outgoing-message":
		log.Println("Yatziebot handling outgoing message")
		var msg pbot.BBSmessage
		msg.UserData = m.NamedArgs
		msg.Message = m.NamedArgs["Message"]
		msg.PayloadString = m.NamedArgs["PayloadString"]

		msg.PayloadBytes, _ = base64.StdEncoding.DecodeString(m.NamedArgs["PayloadBytes"])

		//msg.Message = "text"
		//msg.PayloadString = m.Arg

		handleOneOutgoing(&msg)
	}
	return out
}

func (m *Tpbot) OnStart() {
	chats = map[string]telebot.Chat{}
	log.Println("[pbot] Started")
	plugin_registry.RegisterCommand("p", "pbot ready")

	log.Printf("%+v\n", plugin_registry.Bot)

	runtime.GOMAXPROCS(2)
	conn := svarmrgo.CliConnect()
	svarmrgo.HandleInputLoop(conn, handleMessage)
}

func (m *Tpbot) OnStop() {
	plugin_registry.UnregisterCommand("p")
}

func getUrl(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func extension(buf []byte) string {
	kind, unkwown := filetype.Match(buf)
	if unkwown != nil {
		return "application/octet-stream"
	}
	return kind.Extension
}

func getTelegramFile(fileId string) []byte {
	log.Printf("FileID: %v\n", fileId)
	file, _ := plugin_registry.Bot.GetFile(fileId)
	log.Printf("%+v\n", file)
	url := "https://api.telegram.org/file/bot" + plugin_registry.Config.Token + "/" + file.FilePath
	data := getUrl(url)
	return data
}

func (m *Tpbot) Run(message telebot.Message) {
	config := plugin_registry.Config
	log.Printf("%+v\n", config)

	//if strings.Contains(message.Text, config.CommandPrefix+"p") {
	if message.Chat.Type == "private" {
		log.Println("Yatziebot handling incoming message from user")
		chats[message.Chat.Username] = message.Chat
		if message.Text == "" {
			log.Printf("%+v\n", message)
			photo := message.Photo
			if photo != nil && len(photo) > 2 {
				photo_id := photo[len(photo)-1].File.FileID

				args := map[string]string{}
				args["Message"] = "photo"
				args["Sender"] = message.Chat.Username
				args["Service"] = "Telegram"
				args["PayloadString"] = fmt.Sprintf("%v.pic", photo_id)
				args["PayloadBytes"] = string(getTelegramFile(photo_id))
				svarmrgo.SendMessage(nil, svarmrgo.Message{Selector: "incoming-message", NamedArgs: args})

			} else {
				fileId := message.Document.File.FileID

				args := map[string]string{}
				args["Message"] = "file"
				args["Sender"] = message.Chat.Username
				args["Service"] = "Telegram"
				args["PayloadString"] = message.Document.FileName
				args["PayloadBytes"] = string(getTelegramFile(fileId))
				svarmrgo.SendMessage(nil, svarmrgo.Message{Selector: "incoming-message", NamedArgs: args})

			}
		} else {
			args := map[string]string{}
			args["Message"] = "text"
			args["Sender"] = message.Chat.Username
			args["Service"] = "Telegram"
			args["PayloadString"] = message.Text
			svarmrgo.SendMessage(nil, svarmrgo.Message{Selector: "incoming-message", NamedArgs: args})
		}

	} else {
		bot := plugin_registry.Bot
		bot.SendMessage(message.Chat, "Chat is not private", nil)
	}
	//}
}
