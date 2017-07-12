package tpbot

import (
	"log"
	//"strings"

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

type MagicBallPlugin struct {
}

func init() {
	plugin_registry.RegisterPlugin(&MagicBallPlugin{})
}

func (m *MagicBallPlugin) OnStart() {
	log.Println("[pbot] Started")
	plugin_registry.RegisterCommand("p", "pbot ready")
	bbs = pbot.NewBBS("./")
	bbs.Start()
	go func() {
		for m := range bbs.Outgoing {
			if m.Message == "text" {
				log.Print("Sending message...")
				bot := plugin_registry.Bot
				message := m.UserData.(telebot.Message)
				bot.SendMessage(message.Chat, m.PayloadString, nil)
				log.Println("Done!")
			} else {

				//fmt.Println(res2)
			}
		}
	}()
}

func (m *MagicBallPlugin) OnStop() {
	plugin_registry.UnregisterCommand("p")
}

func (m *MagicBallPlugin) Run(message telebot.Message) {
	//config := plugin_registry.Config
	//if strings.Contains(message.Text, config.CommandPrefix+"p") {
	bbs.Incoming <- &pbot.BBSmessage{Message: "text", PayloadString: message.Text, UserData: message}

	//}
}
