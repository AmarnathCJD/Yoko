package modules

import (
	"fmt"
	"strings"
        "time"

	tb "gopkg.in/tucnak/telebot.v3"
)

func Ban(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	m := c.Message()
	user, xtra := get_user(m)
	if user == nil {
		return nil
	}
	if user.ID == 5050904599 {
		c.Reply("You know what I'm not going to do? Ban myself.")
		return nil
	}
	arg := strings.SplitN(c.Message().Text, " ", 2)
	until_date := 0
	reason := xtra
	if arg[0] == "/tban" {
		args := strings.SplitN(xtra, " ", 2)
		if len(args) == 0 {
			c.Reply("You haven't specified a time to ban this user for!")
			return nil
		}
		until_date = int(Extract_time(c, args[0]))
		if until_date == 0 {
			return nil
		}
		if len(args) == 2 {
			reason = args[1]
		} else {
			reason = ""
		}
                fmt.Println(until_date)
                until_date = int(time.Now().Unix()) + until_date
	} else if arg[0] == "dban" {
		c.Bot().Delete(c.Message().ReplyTo)
	}
	err := b.Ban(m.Chat, &tb.ChatMember{
		User:            user,
		RestrictedUntil: int64(until_date),
	})
	if err == nil {
		if string(reason) != string("") {
			if arg[0] != "/sban" {
				c.Reply(fmt.Sprintf("<b>%s</b> was banned. ~\n<b>Reason:</b> %s", user.FirstName, reason))
			}
			return nil
		}
		if arg[0] != "/sban" {
			c.Reply(fmt.Sprintf("<b>%s</b> was banned. ~", user.FirstName))
		}
		return nil
	}
	c.Reply(fmt.Sprintf("Failed to ban, %s", err.Error()))
	return nil
}
