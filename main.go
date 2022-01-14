package main

import (
	"fmt"
        "strconv"
	mod "github.com/amarnathcjd/yoko/modules"

	bot "github.com/amarnathcjd/yoko/bot"
)

func main() {
	mod.RegHandlers()
	bot.Bot.Start()
	fmt.Println("Bot started!")
}
