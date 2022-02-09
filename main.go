package main

import (
	"log"
	"os"

	bot "github.com/amarnathcjd/yoko/bot"
	mod "github.com/amarnathcjd/yoko/modules"
)

func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
if err != nil {
    log.Print("error opening file: Hi")
}
defer f.Close()

log.SetOutput(f)
log.Print("Bot Started")
	mod.RegisterHandlers()
	// logging
	bot.Bot.Start()
}
