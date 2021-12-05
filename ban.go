package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func ban(m *tb.Message) {
	if m.Private() {
		b.Reply(m, "This command is for groups.")
		return
	}
	user_id, xtra := get_user(m)
        if user_id == string("") {
                b.Reply(m, "You dont seem to be referring to a user or the ID specified is incorrect..")
		return 
        }
        user := get_entity(m, user_id)
	if user == nil {
		return
        }
	err := b.Ban(m.Chat, &tb.ChatMember{
		User: &tb.User{ID: int(user.ID)},
	})
	if err == nil {
                if string(xtra) != string(""){
                    b.Reply(m, "<b>"+user.FirstName+"</b> was banned. ~\n<b>Reason:</b> "+xtra)
		    return
                }
		b.Reply(m, "<b>"+user.FirstName+"</b> was banned. ~")
		return
	}
	b.Reply(m, "Failed to ban, "+string(err.Error()))
}
