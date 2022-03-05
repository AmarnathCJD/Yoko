package modules

import (
	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func Gban(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	u, x := GetUser(c)
	if u.ID == 0 {
		return nil
	}
	if IsBotAdmin(u.ID) {
		return c.Reply("You can't gban a sudo user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't gban yourself !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna gban myself!")
	}
	c.Reply("Snaps the ban hammer!")
	db.GbanUser(u.ID, u.First, x, c.Sender().ID)
	for _, x := range db.GetAllChats() {
		c.Bot().Ban(&tb.Chat{ID: x.Id}, &tb.ChatMember{User: u.User()})
	}
	return nil
}

// soon
