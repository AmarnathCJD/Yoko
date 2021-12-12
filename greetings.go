package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v3"
)

func greet_member(c tb.Context) error {
	m := c.Message()
	b.Reply(m, "Hi")
	fmt.Println(m)
	return nil
}
