package modules

import tb "gopkg.in/tucnak/telebot.v3"
import "fmt"

var ANON = make(map[int]Update)

type Update struct {
	Id    int
	Func  func(tb.Context) error
	Right string
        C tb.Context
}

func AnonAdmin(next tb.HandlerFunc, p string, c tb.Context) {
	ANON[c.Message().ID] = Update{c.Message().ID, next, p, c}
        fmt.Println(ANON)
}

// Soon, today
