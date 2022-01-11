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
	user, xtra := get_user(c.Message())
	if user == nil {
		return nil
	}
	if user.ID == 5050904599 {
		c.Reply("You know what I'm not going to do? Ban myself!")
		return nil
	}
	arg := strings.SplitN(c.Message().Text, " ", 2)
	until_date := 0
	reason := xtra
        if arg[0] == "/unban" {
           err := c.Bot().Unban(c.Chat(), user, true)
           if err != nil && err.Error() == "telegram: not enough rights to restrict/unrestrict chat member (400)"{
              c.Reply("I haven't got the rights to do this.")
              return nil
           }
           if xtra == string(""){
            c.Reply(fmt.Sprintf("✨ %s Permitted to Join the Chat <b>~</b>", user.FirstName))
           } else {
            c.Reply(fmt.Sprintf("✨ %s Permitted to Join the Chat <b>~</b>\n<b>Reason:</b> %s", user.FirstName, reason))
           }
           return nil
        }
	if arg[0] == "/tban" {
                if xtra == string("") {
			c.Reply("You haven't specified a time to ban this user for!")
			return nil
		}
		args := strings.SplitN(xtra, " ", 2)
		until_date = int(Extract_time(c, args[0]))
		if until_date == 0 {
			return nil
		}
		if len(args) == 2 {
			reason = args[1]
		} else {
			reason = ""
		}
		until_date = int(time.Now().Unix()) + until_date
	} else if arg[0] == "/dban" {
		c.Bot().Delete(c.Message().ReplyTo)
	}
	err := b.Ban(c.Message().Chat, &tb.ChatMember{
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
	} else if err.Error() == "telegram: not enough rights to restrict/unrestrict chat member (400)"{
         c.Reply("I haven't got the rights to do this.")
         return nil
        }
	c.Reply(fmt.Sprintf("Failed to ban, %s", err.Error()))
	return nil
}

func Mute(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	user, xtra := get_user(c.Message())
	if user == nil {
		return nil
	}
	if user.ID == 5050904599 {
		c.Reply("You know what I'm not going to do? Mute myself.")
		return nil
	}
        arg := strings.SplitN(c.Message().Text, " ", 2)
	until_date := 0
	reason := xtra
        if arg[0] == "/unmute" {
           err := c.Bot().Restrict(c.Chat(), &tb.ChatMember{
		Rights:          tb.NoRestrictions(),
		User:            user,
	})
           if err != nil && err.Error() == "telegram: not enough rights to restrict/unrestrict chat member (400)"{
              c.Reply("I haven't got the rights to do this.")
              return nil
           }
           if xtra == string(""){
            c.Reply(fmt.Sprintf("✨ %s was unmuted. <b>~</b>", user.FirstName))
           } else {
            c.Reply(fmt.Sprintf("✨ %s was unmuted. <b>~</b>\n<b>Reason:</b> %s", user.FirstName, reason))
           }
           return nil
        }
	if arg[0] == "/tmute" {
		if xtra == string("") {
			c.Reply("You haven't specified a time to mute this user for!")
			return nil
		}
		args := strings.SplitN(xtra, " ", 2)
		until_date = int(Extract_time(c, args[0]))
		if until_date == 0 {
			return nil
		}
		if len(args) == 2 {
			reason = args[1]
		} else {
			reason = ""
		}
		until_date = int(time.Now().Unix()) + until_date
	} else if arg[0] == "/dmute" {
		c.Bot().Delete(c.Message().ReplyTo)
	}
	err := b.Restrict(c.Message().Chat, &tb.ChatMember{
		Rights:          tb.Rights{CanSendMessages: false},
		User:            user,
		RestrictedUntil: int64(until_date),
	})
	if err == nil {
		if string(reason) != string("") {
			if arg[0] != "/smute" {
				c.Reply(fmt.Sprintf("<b>%s</b> was muted. ~\n<b>Reason:</b> %s", user.FirstName, reason))
			}
			return nil
		}
		if arg[0] != "/sban" {
			c.Reply(fmt.Sprintf("<b>%s</b> was muted. ~", user.FirstName))
		}
		return nil
	} else if err.Error() == "telegram: not enough rights to restrict/unrestrict chat member (400)"{
         c.Reply("I haven't got the rights to do this.")
         return nil
        }
	c.Reply(fmt.Sprintf("Failed to mute, %s", err.Error()))
	return nil
}

func Kick(c tb.Context) error {
if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	user, xtra := get_user(c.Message())
	if user == nil {
		return nil
	}
	if user.ID == 5050904599 {
		c.Reply("You know what I'm not going to do? Kick myself.")
		return nil
	}
	reason := xtra
 err := c.Bot().Unban(c.Chat(), user, false)
 check(err)
 if err == nil {
 if reason == string(""){
  c.Reply(fmt.Sprintf("I've kicked <a href='tg://user?id=%d'>%s</a> <b>~</b>", user.ID, user.FirstName))
  return nil
 } else {
  c.Reply(fmt.Sprintf("I've kicked <a href='tg://user?id=%d'>%s</a> <b>~</b>\n<b>Reason:</b> %s", user.ID, user.FirstName, reason))
  return nil
 }
} else {
 c.Reply("Failed to kick, make sure I'm admin and can RestrictMembers.")
 return nil
}
}
