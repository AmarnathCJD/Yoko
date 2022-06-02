package main

import (
	"fmt"
	"log"

	bot "github.com/amarnathcjd/yoko/bot"
	mod "github.com/amarnathcjd/yoko/modules"
	db "github.com/amarnathcjd/yoko/modules/db"
)

func main() {
	fmt.Println(db.GetFiltersFromDB(100))
	return
	log.Print("Bot Started.")
	mod.RegisterHandlers()
	bot.Bot.Start()
}
