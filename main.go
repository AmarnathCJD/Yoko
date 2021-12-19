package main

import (
	"context"
	"fmt"
        mod "github.com/amarnathcdj/yoko/modules"
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
	b.Handle("/info", mod.info)
	b.Handle("/imdb", mod.IMDb)
	b.Handle("/crypto", mod.Crypto)
	b.Handle("/start", mod.start)
	b.Handle("/save", mod.save, mod.change_info)
	b.Handle("/lock", mod.lock, mod.change_info)
	b.Handle("/locktypes", mod.locktypes)
	b.Handle("/locks", mod.check_locks)
	b.Handle("/unlock", mod.unlock, mod.change_info)
	b.Handle("/eval", mod.evaluate)
	b.Handle("/sh", mod.execute)
	b.Handle("/tr", mod.translate)
	b.Handle("/saved", mod.all_notes)
	b.Handle("/notes", mod.all_notes)
	b.Handle("/get", mod.gnote)
	b.Handle("/promote", mod.promote, mod.add_admins)
	b.Handle("/superpromote", mod.promote, mod.add_admins)
	b.Handle("/demote", mod.demote, mod.add_admins)
	b.Handle("/adminlist", mod.adminlist)
	b.Handle(tb.OnText, mod.hash_note, mod.hash_regex)
        b.Handle("/ud", mod.uD)
        b.Handle("/newfed", new_fed)
	b.Handle(tb.OnChatMember, func(c tb.Context) error {
		fmt.Println(c)
		return nil
	})
	b.Handle(tb.OnUserJoined, mod.greet_member)
	b.Start()
}
