package modules

import (
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func WARN(c tb.Context) {
	cmd := strings.SplitN(c.Message().Text, " ", 2)[0]
	if cmd == "/dwarn" && !c.Message().IsReply() {
		c.Reply("You have to reply to a message to delete it and warn the user.")
		return nil
	}
	user, extra := get_user(c.Message())
	if user.ID == 6 {
		c.Reply("Do you really think I can do that to myself <b>:p</b>")
		return nil
	}
	p := c.Bot().ChatMemberOf(c.Chat(), user)
	if p.Role != "member" {
		c.Reply("✨ I'm not going to warn an admin!")
		return nil
	}
	exceeded := db.Warn_user(c.Chat().ID, user.ID, extra)
        fmt.Println(exceeded)
}
