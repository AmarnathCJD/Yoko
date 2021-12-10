package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func greet_member(m *tb.Message) {
	b.Reply(m, "Hi")
	fmt.Println(m)
}
