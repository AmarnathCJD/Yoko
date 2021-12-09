package main

import (
        "context"
        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
    "fmt"
)

var (
	menu = &tb.ReplyMarkup{}
        db, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://go:amar0245@lon5-c12-1.mongo.objectrocket.com:43391,lon5-c12-2.mongo.objectrocket.com:43391,lon5-c12-0.mongo.objectrocket.com:43391/go?replicaSet=24e1adf7f54a48fba7350c36009da162"))
)

var (
	b, _ = tb.NewBot(tb.Settings{
		Token:       "5050904599:AAG-YrM6KN4EJJx8peQOn901qHhLCkFo5QA",
		Synchronous: false,
		Poller:      &tb.LongPoller{Timeout: 10},
		ParseMode:   "HTML",
	})
)


func main() {
        fmt.Println(db)
	b.Handle("/eval", evaluate)
	b.Handle("/start", func(m *tb.Message) {
		if m.Private() {
			menu.Inline(
				menu.Row(menu.URL("Support", "t.me/roseloverx_support"), menu.URL("Updates", "t.me/roseloverx_support")),
				menu.Row(menu.Data("Commands", "help_menu")),
				menu.Row(menu.URL("Add me to your group", "https://t.me/Yoko_Robot?startgroup=true")))
			b.Send(m.Sender, "Hey there! I am <b>Yoko</b>.\nIm an Anime themed Group Management Bot, feel free to add me to your groups!", menu)
			return
		}
		b.Reply(m, "Hey I'm Alive.")
	})
	b.Handle("/info", info)
	b.Handle("/sh", execute)
	b.Handle("/ban", ban)
        b.Handle("/gp", gp)
        b.Handle(tb.OnChatMember, greet_member)
        b.Handle("/imdb", IMDb)
        b.Handle("/lock", lock)
        b.Handle("/locktypes", locktypes)
	b.Start()
}
