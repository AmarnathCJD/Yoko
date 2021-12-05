package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func ban(m *tb.Message) {
	if m.Private() {
		b.Reply(m, "This command is for groups.")
		return
	}
	user_id, _ := get_user(m)
        if user_id == nil {
		return 
        }
        user = get_entity(user_id)
	if user == nil {
		return
        }
	err := b.Ban(m.Chat, &tb.ChatMember{
		User: user,
	})
	if err == nil {
		b.Reply(m, "<b>"+user.FirstName+"</b> was banned. ~")
		return
	}
	b.Reply(m, "Failed to ban, "+string(err.Error()))
}
