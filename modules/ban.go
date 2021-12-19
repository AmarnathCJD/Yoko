package modules

import (
	tb "gopkg.in/tucnak/telebot.v3"
)

func ban(c tb.Context) error {
	m := c.Message()
	if m.Private() {
		b.Reply(m, "This command is for groups.")
		return nil
	}
	user, xtra := get_user(m)
	if user == nil {
		return nil
	}
	err := b.Ban(m.Chat, &tb.ChatMember{
		User: user,
	})
	if err == nil {
		if string(xtra) != string("") {
			b.Reply(m, "<b>"+user.FirstName+"</b> was banned. ~\n<b>Reason:</b> "+xtra)
			return nil
		}
		b.Reply(m, "<b>"+user.FirstName+"</b> was banned. ~")
		return nil
	}
	b.Reply(m, "Failed to ban, "+string(err.Error()))
	return nil
}
