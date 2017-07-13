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

var quips = []string{
	"It is certain", "It is decidedly so", "Without a doubt",
	"Yes definitely", "You may rely on it", "As I see it, yes",
	"Most likely", "Outlook good", "Yes", "Signs point to yes",
	"Reply hazy try again", "Ask again later",
	"Better not tell you now", "Cannot predict now",
	"Concentrate and ask again", "Don\"t count on it",
	"My reply is no", "My sources say no", "Outlook not so good",
	"Very doubtful",
}

type Tpbot struct {
}

func init() {
	plugin_registry.RegisterPlugin(&Tpbot{})
}

func (m *Tpbot) OnStart() {
	log.Println("[pbot] Started")
	plugin_registry.RegisterCommand("p", "pbot ready")
	bbs = pbot.NewBBS("./")
	bbs.Start()
	go func() {
		for m := range bbs.Outgoing {
			bot := plugin_registry.Bot
				message := m.UserData.(telebot.Message)
			if m.Message == "text" {
				log.Print("Sending message...", m.PayloadString)
				bot.SendMessage(message.Chat, m.PayloadString, nil)
				log.Println("Done!")
			} else {
				photo, err := telebot.NewFile("botfiles/files/"+m.PayloadString)
				if err != nil {
					log.Println("Error creating the new file ")
					log.Println(err)
					bot.SendMessage(message.Chat, "Error creating the new file ", nil)

				} else {
					mtype := fileType("botfiles/files/"+m.PayloadString)
					if strings.HasPrefix(mtype, "image") {
						picture := telebot.Photo{File: photo}

						err = bot.SendPhoto(message.Chat, &picture, nil)
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

func (m *Tpbot) Run(message telebot.Message) {
	//config := plugin_registry.Config
	//if strings.Contains(message.Text, config.CommandPrefix+"p") {
	if message.Chat.Username == "donomii" && message.Chat.Type == "private" {
	if message.Text == "" {
		fmt.Printf("%+v\n", message)
		fileId := message.Document.File.FileID
		fmt.Printf("FileID: %v\n", fileId)
		file, _ := plugin_registry.Bot.GetFile(fileId)
		fmt.Printf("%+v\n",file)
		url := "https://api.telegram.org/file/bot" + plugin_registry.Config.Token + "/" + file.FilePath
		data := getUrl(url)
		err := ioutil.WriteFile("botfiles/files/"+message.Document.FileName, data, 0644) 
		fmt.Println(err)
	} else {
		bbs.Incoming <- &pbot.BBSmessage{Message: "text", PayloadString: message.Text, UserData: message}
	}
	}
	//}
}
