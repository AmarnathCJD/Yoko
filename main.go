package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	menu  = &tb.ReplyMarkup{}
	db, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://go:amar0245@lon5-c12-1.mongo.objectrocket.com:43391,lon5-c12-2.mongo.objectrocket.com:43391,lon5-c12-0.mongo.objectrocket.com:43391/go?replicaSet=24e1adf7f54a48fba7350c36009da162&retryWrites=false"))
)

var (
	b, _ = tb.NewBot(tb.Settings{
		Token:       "5050904599:AAEqSus8jD5vTVZzQeTVo4cBc8IwsRha_I0",
		Synchronous: false,
		Poller:      &tb.LongPoller{Timeout: 10},
		ParseMode:   "HTML",
	})
)

func main() {
	b.Handle("/eval", evaluate)
	b.Handle("/start", start)
	b.Handle("/info", info)
	b.Handle("/sh", execute)
	b.Handle("/ban", ban)
	b.Handle("/gp", gp)
	b.Handle(tb.OnChatMember, greet_member)
	b.Handle("/imdb", IMDb)
	b.Handle("/lock", lock)
	b.Handle("/locktypes", locktypes)
        b.Handle("/locks", check_locks)
	b.Start()
}
