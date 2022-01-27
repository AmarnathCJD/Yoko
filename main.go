package main

import (
	mod "github.com/amarnathcjd/yoko/modules"
	"log"
	"os"

	bot "github.com/amarnathcjd/yoko/bot"
)

func main() {
	mod.RegisterHandlers()
	file, _ := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)
	log.Print("Bot started")
	bot.Bot.Start()
}
