package modules

import (
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func Purge(c tb.Context) error {
	if !c.Message().IsReply() {
		return c.Reply("Reply to a message to show me where to purge from.")
	}
	for i := Reply.ID; i <= c.Message().ID; i++ {
		ID := i
		go func() {
			c.Bot().Delete(&tb.Message{ID: ID, Chat: c.Message().Chat})
		}()
	}
	c.Delete()
	return c.Send("Purge complete.")
}

func Delete(c tb.Context) error {
 
if !c.Message().IsReply() {
		return c.Reply("Reply to a message to delete it.")
	}
c.Bot().Delete(c.Message().ReplyTo)

}
