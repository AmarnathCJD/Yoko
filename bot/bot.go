package bot

import (
	"fmt"
"os"
	tb "gopkg.in/telebot.v3"
)


func BotInit() tb.Bot {
	b, _ := tb.NewBot(tb.Settings{
		URL:         "",
		Token:       "5181620060:AAF6FCOW9M6tbUPPtLs_b2UkXKGz7ksaggo",
		Poller:      &tb.Webhook{
  Listen: os.Getenv("PORT"),
  Endpoint: &tb.WebhookEndpoint{
    PublicURL: "https://golang-yoko.herokuapp.com/",
  },
  MaxConnections: 100,
  AllowedUpdates: []string{"message", "chat_member", "inline_query", "callback_query"},
}
		Synchronous: false,
		Verbose:     false,
		ParseMode:   "HTML",
		Offline:     false,
		OnError: func(e error, c tb.Context) {
			fmt.Println(e)
		},
	})

	return *b
}

var Bot = BotInit()
