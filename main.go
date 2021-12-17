package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tb "gopkg.in/tucnak/telebot.v3"
)

var (
	menu  = &tb.ReplyMarkup{}
	db, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://go:amar0245@lon5-c12-1.mongo.objectrocket.com:43391,lon5-c12-2.mongo.objectrocket.com:43391,lon5-c12-0.mongo.objectrocket.com:43391/go?replicaSet=24e1adf7f54a48fba7350c36009da162&retryWrites=false"))
)

var (
	b, _ = tb.NewBot(tb.Settings{
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
)

func main() {
	b.Handle("/info", info)
	b.Handle("/imdb", IMDb)
	b.Handle("/crypto", Crypto)
	b.Handle("/start", start)
	b.Handle("/save", save, change_info)
	b.Handle("/lock", lock, change_info)
	b.Handle("/locktypes", locktypes)
	b.Handle("/locks", check_locks)
	b.Handle("/unlock", unlock, change_info)
	b.Handle("/eval", evaluate)
	b.Handle("/sh", execute)
	b.Handle("/tr", translate)
	b.Handle("/saved", all_notes)
	b.Handle("/notes", all_notes)
	b.Handle("/get", gnote)
	b.Handle("/promote", promote, add_admins)
	b.Handle("/superpromote", promote, add_admins)
	b.Handle("/demote", demote, add_admins)
	b.Handle("/adminlist", adminlist)
	b.Handle(tb.OnText, hash_note, hash_regex)
	b.Handle(tb.OnChatMember, func(c tb.Context) error {
		fmt.Println(c)
		return nil
	})
	b.Handle(tb.OnUserJoined, greet_member)
	b.Start()
}
