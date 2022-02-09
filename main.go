package main

import (
	"log"
	"os"

	bot "github.com/amarnathcjd/yoko/bot"
	mod "github.com/amarnathcjd/yoko/modules"
)

func main() {
	file, _ := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.Print("Bot started")
	log.SetOutput(file)
	mod.RegisterHandlers()
	// logging
	bot.Bot.Start()
}
