package tpbot

import (
	"strings"
	"gopkg.in/h2non/filetype.v1"
	"io/ioutil"
	"net/http"
	"log"
	//"strings"
	"fmt"

	"github.com/donomii/pbot"
	"github.com/donomii/yatzie/shared/registry"
	//"github.com/donomii/yatzie/shared/utils"
	"github.com/tucnak/telebot"
)

var bbs *pbot.BBSdata

type Tpbot struct {
}

func init() {
	plugin_registry.RegisterPlugin(&Tpbot{})
}

var chats map[string]telebot.Chat

func (m *Tpbot) OnStart() {
	log.Println("[pbot] Started")
	plugin_registry.RegisterCommand("p", "pbot ready")
	bbs = pbot.NewBBS("./")
	bbs.Start()
	fmt.Printf("%+v\n",plugin_registry.Bot)
	go func() {
		for m := range bbs.Outgoing {
			bot := plugin_registry.Bot
			var message telebot.Message
			if m.UserData != nil {
				message = m.UserData.(telebot.Message)
				chats[message.Chat.Username] = message.Chat
			} else {
				message.Chat = chats[bbs.Config.Get("TelegramOwner")]
			}
			if m.Message == "text" {
				log.Print("Sending message...", m.PayloadString)
				err := bot.SendMessage(message.Chat, m.PayloadString, nil)
				log.Println("Done!", err)
			} else {
				data := m.PayloadBytes
				name := bbs.TempDir+"/"+m.PayloadString
				err := ioutil.WriteFile(name, data, 0644)
				if err != nil {
					log.Println("Error creating the new file ")
					log.Println(err)
					bot.SendMessage(message.Chat, "Error creating the new file ", nil)

				} else {
					//WTF
					photo, err := telebot.NewFile(name)
					mtype := fileType(name)
					if strings.HasPrefix(mtype, "image") {
						picture := telebot.Photo{File: photo}

						err = bot.SendPhoto(message.Chat, &picture, &telebot.SendOptions{ReplyTo: message})
						if err != nil {
							log.Println("Error sending photo")
							log.Println(err)
							bot.SendMessage(message.Chat, "Could not send photo", nil)
						}
					} else {
						d:= telebot.Document{ File: photo, FileName: m.PayloadString, Mime: mtype}
						err = bot.SendDocument(message.Chat, &d, nil)
						if err != nil {
							log.Println("Error sending document")
							log.Println(err)
							bot.SendMessage(message.Chat, "Could not send photo", nil)
						}
						
					}
				}
			}
		}
	}()
}

func (m *Tpbot) OnStop() {
	plugin_registry.UnregisterCommand("p")
}

func getUrl (url string) []byte {
resp, err := http.Get(url)
if err != nil {
	// handle error
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
return body
}

func fileType (path string) string{
  buf, _ := ioutil.ReadFile(path)

  kind, unkwown := filetype.Match(buf)
  if unkwown != nil {
    return "application/octet-stream"
  }
  return kind.MIME.Value
}

func extension (buf []byte) string{
  kind, unkwown := filetype.Match(buf)
  if unkwown != nil {
    return "application/octet-stream"
  }
  return kind.Extension
}


func getTelegramFile(fileId string) []byte {
	fmt.Printf("FileID: %v\n", fileId)
	file, _ := plugin_registry.Bot.GetFile(fileId)
	fmt.Printf("%+v\n",file)
	url := "https://api.telegram.org/file/bot" + plugin_registry.Config.Token + "/" + file.FilePath
	data := getUrl(url)
	return data
}



func (m *Tpbot) Run(message telebot.Message) {
	config := plugin_registry.Config
	fmt.Printf("%+v\n",config)
	//if strings.Contains(message.Text, config.CommandPrefix+"p") {
	if message.Chat.Username == bbs.Config.Get("TelegramOwner") && message.Chat.Type == "private" {
	chats = map[string]telebot.Chat{}
	chats[message.Chat.Username] = message.Chat
	if message.Text == "" {
		fmt.Printf("%+v\n", message)
		photo := message.Photo
		if photo != nil && len(photo)>2 {
			photo_id := photo[len(photo)-1].File.FileID			
			data := getTelegramFile(photo_id)
			bbs.Files.PutBytes(fmt.Sprintf("%v.pic",photo_id), data)
		} else {
			fileId := message.Document.File.FileID
			data := getTelegramFile(fileId)
			bbs.Files.PutBytes(message.Document.FileName, data) 
			}
	} else {
		bbs.Incoming <- &pbot.BBSmessage{Message: "text", PayloadString: message.Text, UserData: message}
	}
	} else {
		bot := plugin_registry.Bot
		bot.SendMessage(message.Chat, "User '"+ message.Chat.Username + "' not recognised or chat not private", nil)
	}
	//}
}
