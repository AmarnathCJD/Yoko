package bot

import (
	"fmt"
	tb "gopkg.in/telebot.v3"
	"log"
	"os"
)

func BotInit() tb.Bot {
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
