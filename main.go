package main

import (
	"io"
	"log"
	"os"

	bot "github.com/amarnathcjd/yoko/bot"
	mod "github.com/amarnathcjd/yoko/modules"
)

func main() {
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	writer := io.MultiWriter(os.Stdout, os.Stderr)
	log.SetOutput(writer)
	log.SetOutput(f)
	log.Print("Bot Started.")
	mod.RegisterHandlers()
	bot.Bot.Start()
}
