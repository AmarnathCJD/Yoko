package main

import (
	"fmt"
	mod "github.com/amarnathcjd/yoko/modules"

	bot "github.com/amarnathcjd/yoko/bot"
)

func main() {
	mod.RegHandlers()
	bot.Bot.Start()
	fmt.Println("Bot started!")
}
