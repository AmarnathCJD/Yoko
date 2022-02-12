package modules

import (
	"fmt"
	"time"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func AddSudo(c tb.Context) error {
	if !IsDev(c.Sender().ID) && c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := get_user(c.Message())
	if u == nil {
		return nil
	}
	if IsSudo(u.ID) {
		return c.Reply("This user is already a sudo user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't add yourself to the bot admin list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna add myself to the bot admin list!")
	}
	db.AddSudo(u.ID, u.FirstName)
	return c.Reply("Added user to sudo list !")
}

func AddDev(c tb.Context) error {
	if c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := get_user(c.Message())
	if u == nil {
		return nil
	}
	if IsDev(u.ID) {
		return c.Reply("This user is already a dev user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't add yourself to the dev list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna add myself to the dev list!")
	}
	db.AddDev(u.ID, u.FirstName)
	return c.Reply("Added user to dev list !")
}

func ListSudo(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	return c.Reply(db.ListSudo())
}

func ListDev(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	return c.Reply(db.ListDev())
}

func RemoveSudo(c tb.Context) error {
	if !IsDev(c.Sender().ID) && c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := get_user(c.Message())
	if u == nil {
		return nil
	}
	if !IsSudo(u.ID) {
		return c.Reply("This user is not a sudo user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't remove yourself from the sudo list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna remove myself from the sudo list!")
	}
	db.RemSudo(u.ID)
	return c.Reply("Removed user from sudo list !")
}

func RemoveDev(c tb.Context) error {
	if c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := get_user(c.Message())
	if u == nil {
		return nil
	}
	if !IsDev(u.ID) {
		return c.Reply("This user is not a dev user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't remove yourself from the dev list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna remove myself from the dev list!")
	}
	db.RemDev(u.ID)
	return c.Reply("Removed user from dev list !")
}

func Logs(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	} else {
		return c.Reply(&tb.Document{
			File:     tb.File{FileLocal: "log.txt"},
			Caption:  time.Now().String(),
			FileName: "log.txt",
		})
	}
}

func Ping(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	a := time.Now()
	alive := a.Sub(StartTime)
	msg, _ := c.Bot().Send(c.Chat(), "<code>Pinging!</code>")
	b := time.Now()
	_, err := c.Bot().Edit(msg, fmt.Sprintf("<b>► Ping</b>: <code>%s</code>\n<b>► Uptime:</b> %d", b.Sub(a).String(), int(alive)))
	return err
}
