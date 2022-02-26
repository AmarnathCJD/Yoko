package modules

import (
	tb "gopkg.in/telebot.v3"
)

func Purge(c tb.Context) error {
	if !c.Message().IsReply() {
		return c.Reply("Reply to a message to show me where to purge from.")
	}
	for i := c.Message().ReplyTo.ID; i <= c.Message().ID; i++ {
		ID := i
		go func() {
			err := c.Bot().Delete(&tb.Message{ID: ID, Chat: c.Message().Chat})
                        check (err)
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
        return c.Delete()
}

func PurgeFrom() {

}

func PurgeTo() {

}
