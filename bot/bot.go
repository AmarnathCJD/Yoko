package bot

import (
	"fmt"
	"log"
	"os"

	dotenv "github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"
)

func BotInit() tb.Bot {
	if err := dotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	if os.Getenv("TOKEN") == "" {
		log.Fatal("Please set the TOKEN environment variable")
	}
	b, _ := tb.NewBot(tb.Settings{
		URL:         "",
		Token:       os.Getenv("TOKEN"),
		Poller:      &tb.LongPoller{Timeout: 10, AllowedUpdates: []string{"message", "chat_member", "inline_query", "callback_query"}},
		Synchronous: false,
		Verbose:     false,
		ParseMode:   "HTML",
		Offline:     false,
		OnError: func(e error, c tb.Context) {
			fmt.Println(e)
			log.Println(e)
		},
	})

	return *b
}

var Bot = BotInit()
