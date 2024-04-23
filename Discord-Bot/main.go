package main

import (
	"fmt"

	"github.com/lokesh1jha/discordbot/bot"
	"github.com/lokesh1jha/discordbot/config"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
}
