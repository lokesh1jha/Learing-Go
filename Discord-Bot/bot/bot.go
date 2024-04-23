package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/lokesh1jha/discordbot/config"
)

var BotId string
var goBot *discordgo.Session

func Start() {

	fmt.Println("Starting bot...")

	var err error
	goBot, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("New error" + err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println("User error" + err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// bot should not respond to itself
	if m.Author.ID == BotId {
		return
	}

	// Command is Ping, Response is Pong
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
