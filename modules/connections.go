package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func ConnectChat(c tb.Context) error {
	if !c.Message().Private() {
		sel.Inline(sel.Row(sel.URL("Connect to chat", fmt.Sprintf("t.me/missmikabot?start=connect_%d", c.Chat().ID))))
		c.Reply("Tap the following button to connect to this chat in PM", sel)
	} else if c.Message().Payload == string("") {
		c.Reply("I need a chat id to connect to!")
		return nil
	}
	chat_id, err := strconv.Atoi(c.Message().Payload)
	if err != nil {
		return c.Reply(err.Error())
	}
	Chat, _ := b.ChatByID(int64(chat_id))
	db.ConnectChat(int64(chat_id), c.Sender().ID)
	P, _ := b.ChatMemberOf(Chat, c.Sender())
	if P.Role == tb.Member {
		sel.Inline(sel.Row(sel.Data("User commands", "connect_us_cmd")))
	} else {
		sel.Inline(sel.Row(sel.Data("Admin Commands", "connect_ad_cmd")), sel.Row(sel.Data("User commands", "connect_us_cmd")))
	}
	return c.Reply(fmt.Sprintf("You have been connected to %s!", Chat.Title), sel)
}

func PrivateConnect(c tb.Context) error {
	args := strings.SplitN(c.Message().Payload, "_", 2)
	chat_id, _ := strconv.Atoi(args[1])
	chat, err := c.Bot().ChatByID(int64(chat_id))
	check(err)
	P, _ := b.ChatMemberOf(chat, c.Sender())
	if P.Role == tb.Member {
		sel.Inline(sel.Row(sel.Data("User commands", "connect_us_cmd")))
	} else {
		sel.Inline(sel.Row(sel.Data("Admin Commands", "connect_ad_cmd")), sel.Row(sel.Data("User commands", "connect_us_cmd")))
	}
	c.Reply(fmt.Sprintf("You have been connected to %s!", chat.Title), sel)
	db.ConnectChat(chat.ID, c.Sender().ID)
	return nil
}

func ChatContext(c tb.Context) tb.Context {
	if !c.Message().Private() {
		return c
	}
	chat_id := db.GetChat(c.Sender().ID)
	fmt.Println(chat_id)
	if chat_id == int64(0) {
		return c
	}
	cmd := strings.Split(c.Text(), " ")[0][1:]
	fmt.Println(cmd)
	if !stringInSlice(cmd, CNT) {
		return c
	}
	Chat, _ := b.ChatByID(chat_id)
	return c.Bot().NewContext(tb.Update{ID: c.Message().ID,
		Message: &tb.Message{
			Sender:      c.Sender(),
			Chat:        Chat,
			Payload:     c.Message().Payload,
			Text:        c.Message().Text,
			ReplyTo:     c.Message().ReplyTo,
			Audio:       c.Message().Audio,
			Video:       c.Message().Video,
			Document:    c.Message().Document,
			Photo:       c.Message().Photo,
			Sticker:     c.Message().Sticker,
			Voice:       c.Message().Voice,
			Animation:   c.Message().Animation,
			ReplyMarkup: c.Message().ReplyMarkup,
			ID:          c.Message().ID,
		},
	})
}
