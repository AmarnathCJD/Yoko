package bot

import (
	tb "gopkg.in/tucnak/telebot.v3"
)

func BotInit() tb.Bot {
	b, _ := tb.NewBot(tb.Settings{
		URL:         "",
		Token:       "5050904599:AAEqSus8jD5vTVZzQeTVo4cBc8IwsRha_I0",
		Poller:      &tb.LongPoller{Timeout: 10, AllowedUpdates: []string{"message", "chat_member"}},
		Synchronous: false,
		Verbose:     false,
		ParseMode:   "HTML",
		OnError: func(error, tb.Context) {
		},
		Offline: false,
	})
	return *b
}

var Bot = BotInit()
