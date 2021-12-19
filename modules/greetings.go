package modules

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v3"
)

func greet_member(c tb.Context) error {
	m := c.Message()
	fmt.Println(m.Chat.Title)
	return nil
}
