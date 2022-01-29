package modules

import tb "gopkg.in/tucnak/telebot.v3"

func Gban(c tb.Context) error {
	if !IS_SUDO(c.Sender().ID) {
		return nil
	}
	u, _ := get_user(c.Message())
	if u == nil {
		return nil
	}
	if IS_SUDO(u.ID) {
		return c.Reply("You can't gban a sudo user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't gban yourself !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna gban myself!")
	}
	return nil
}

// soon
