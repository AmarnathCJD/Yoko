package modules

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v3"
)

var menu = &tb.ReplyMarkup{}

func Start(c tb.Context) error {
	m := c.Message()
	if strings.HasPrefix(m.Payload, "notes") {
		private_start_note(c)
		return nil
	} else if strings.HasPrefix(m.Payload, "allnotes") {
		private_startallnotes(c)
		return nil
	}
	if m.Private() {
		menu.Inline(
			menu.Row(menu.URL("Support", "t.me/roseloverx_support"), menu.URL("Updates", "t.me/roseloverx_support")),
			menu.Row(menu.Data("Commands", "help_menu")),
			menu.Row(menu.URL("Add me to your group", "https://t.me/Aiko_Robot?startgroup=true")))
		b.Send(m.Sender, "Hey there! I am <b>Aiko</b>.\nIm an Anime themed Group Management Bot, feel free to add me to your groups!", menu)
		return nil
	}
	b.Reply(m, "Hey I'm Alive.")
	return nil
}
