package modules

import (
	"fmt"
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)


func WARN(c tb.Context) error {
	cmd := strings.SplitN(c.Message().Text, " ", 2)[0]
	if cmd == "/dwarn" && !c.Message().IsReply() {
		c.Reply("You have to reply to a message to delete it and warn the user.")
		return nil
	}
	user, extra := get_user(c.Message())
	if user.ID == int64(6) {
		c.Reply("Do you really think I can do that to myself <b>:p</b>")
		return nil
	}
	p, err := c.Bot().ChatMemberOf(c.Chat(), user)
	if err != nil {
		c.Reply(err.Error())
		return nil
	}
	if stringInSlice(p.Role, []string{"administrator", "creator"}) {
		c.Reply("✨ I'm not going to warn an admin!")
		return nil
	}
	exceeded, limit, count := db.Warn_user(c.Chat().ID, user.ID, extra)
	if extra == string("") {
		extra = "No reason given."
	}
	if !exceeded {
                menu.Inline(menu.Row(menu.Data("Remove warn (admin only)", fmt.Sprintf("unwarn_btn_%d", user.ID))))
		c.Reply(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> has %d/%d warnings; be careful!\n<b>Reason</b>: %s", user.ID, user.FirstName, count, limit, extra), menu)
		return nil
	}
	return nil
}
