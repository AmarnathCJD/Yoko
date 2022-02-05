package modules

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
)

func PARSET(c tb.Context) error {
	return c.Reply(ParseMD(c))

}

func ParseMD(c tb.Context) string {
	return fmt.Sprint(c.Message().Entities)

}
