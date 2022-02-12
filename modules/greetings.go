package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func Welcome_set(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	if c.Message().Payload == string("") {
		text, file, mode := db.Get_welcome(c.Chat().ID)
		msg, err := c.Bot().Send(c.Chat(), fmt.Sprintf("<b>Greetings config in this chat</b>:\n- Should greet new members: <code>%s</code>\n- Delete old welcome message: <code>%s</code>\n- Delete welcome service: <code>%s</code>\n\nWelcome message:", strings.Title((strconv.FormatBool(mode))), "True", "True"), &tb.SendOptions{ReplyTo: c.Message()})
		check((err))
		if mode {
			if len(file) != 0 {
				text, btns := button_parser(text)
				f := GetFile(file, text)
				f.Send(c.Bot(), c.Chat(), &tb.SendOptions{ReplyTo: msg, ReplyMarkup: btns, DisableWebPagePreview: true})
			} else {
				c.Send(text, &tb.SendOptions{ReplyTo: msg, DisableWebPagePreview: true})
			}

		}
	} else if strings.ToLower(c.Message().Payload) == "raw" || strings.ToLower(c.Message().Payload) == "noformat" {
		text, _, _ := db.Get_welcome(c.Chat().ID)
		if text != string("") {
			c.Reply(fmt.Sprint(text))
		}
	} else if stringInSlice(strings.ToLower(c.Message().Payload), []string{"yes", "on", "enable", "y"}) {
		c.Reply("I'll be welcoming all new members from now on!")
		db.Set_welcome_mode(c.Chat().ID, true)
	} else if stringInSlice(strings.ToLower(c.Message().Payload), []string{"no", "off", "disable", "n"}) {
		c.Reply("I'll stay quiet when new members join.")
		db.Set_welcome_mode(c.Chat().ID, false)
		return nil
	} else {
		c.Reply("Your input was not recognised as one of: yes/no/on/off")
		return nil
	}
	return nil
}

func Set_welcome(c tb.Context) error {
	Text, Text2, File := parse_message(c.Message())
	if c.Message().IsReply() {
		Text += Text2
	}
	if Text == string("") && File == nil {
		c.Reply("You need to give the welcome message some content!")
		return nil
	}
	c.Reply("The new welcome message has been saved!")
	db.Set_welcome(c.Chat().ID, Text, File)
	return nil
}

func ResetWelcome(c tb.Context) error {
	c.Reply("The welcome message has been reset!")
	db.Reset_welcome(c.Chat().ID)
	return nil
}

func OnChatMemberHandler(c tb.Context) error {
	upd := c.ChatMember()
	if upd.NewChatMember != nil && upd.OldChatMember != nil {
		if upd.NewChatMember.Role == tb.Member && upd.OldChatMember.Role == tb.Left {
			text, file, mode := db.Get_welcome(c.Chat().ID)
			if !mode {
				return nil
			}
			text, btns := button_parser(text)
			if len(file) != 0 {
				text, p := ParseString(text, c)
				f := GetFile(file, text)
				f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns})
			} else {
				text, p := ParseString(text, c)
				c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns})
			}
		}
	}
	return nil
}
