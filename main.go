package main

import (
	"log"

	bot "github.com/amarnathcjd/yoko/bot"
	mod "github.com/amarnathcjd/yoko/modules"
)

func main() {
	log.Print("Bot Started.")
	mod.RegisterHandlers()
	bot.Bot.Start()
}
