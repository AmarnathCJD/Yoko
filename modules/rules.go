package modules

import (
	tb "gopkg.in/telebot.v3"
)

func SetRules(c tb.Context) error {
	return c.Reply("Soon")
}
