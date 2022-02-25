package modules

import (
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func Purge(c tb.Context) error {
	Args := c.Message().Payload
	Reply := c.Message().ReplyTo
	var lt = Reply.ID
	if isInt(Args) {
		ArgsInt, _ := strconv.Atoi(Args)
		if ArgsInt > 0 {
			lt = ArgsInt
		}
	}
	if Reply == nil {
		return c.Reply("Reply to a message to show me where to purge from.")
	}
	for i := Reply.ID; i <= c.Message().ID-(Reply.ID-lt); i++ {
		c.Bot().Delete(&tb.Message{ID: i, Chat: c.Message().Chat})
	}
	c.Delete()
	return c.Send("Purge complete.")
}
