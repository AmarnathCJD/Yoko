package modules

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
)

func Connect_chat(c tb.Context) error {
	if !c.Message().Private() {
		sel.Inline(sel.Row(sel.URL("Connect to chat", fmt.Sprintf("http://t.me/yoko_robot?start=connect_%d", c.Chat().ID))))
		c.Reply("Tap the following button to connect to this chat in PM", sel)
	}
	return nil
}
