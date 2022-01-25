package main

import (
	"log"

	mod "github.com/amarnathcjd/yoko/modules"

	bot "github.com/amarnathcjd/yoko/bot"
)

func main() {
	mod.RegisterHandlers()
	log.Print("Bot started")
	bot.Bot.Start()
}
