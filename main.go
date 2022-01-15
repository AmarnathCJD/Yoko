package main

import (
	mod "github.com/amarnathcjd/yoko/modules"

	bot "github.com/amarnathcjd/yoko/bot"
)

func main() {
	mod.RegHandlers()

	bot.Bot.Start()
}
