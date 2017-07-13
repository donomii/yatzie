package tpbot

import (
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
					picture := telebot.Photo{File: photo}

					err = bot.SendPhoto(message.Chat, &picture, nil)
					if err != nil {
						log.Println("Error sending photo")
						log.Println(err)
						bot.SendMessage(message.Chat, "Could not send photo", nil)
					}
				}
			}
		}
	}()
}

func (m *Tpbot) OnStop() {
	plugin_registry.UnregisterCommand("p")
}

func (m *Tpbot) Run(message telebot.Message) {
	//config := plugin_registry.Config
	//if strings.Contains(message.Text, config.CommandPrefix+"p") {
	fmt.Println(message)
	bbs.Incoming <- &pbot.BBSmessage{Message: "text", PayloadString: message.Text, UserData: message}

	//}
}
