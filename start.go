package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func start(m *tb.Message) {
	if m.Private() {
		menu.Inline(
			menu.Row(menu.URL("Support", "t.me/roseloverx_support"), menu.URL("Updates", "t.me/roseloverx_support")),
			menu.Row(menu.Data("Commands", "help_menu")),
			menu.Row(menu.URL("Add me to your group", "https://t.me/Yoko_Robot?startgroup=true")))
		b.Send(m.Sender, "Hey there! I am <b>Yoko</b>.\nIm an Anime themed Group Management Bot, feel free to add me to your groups!", menu)
		return
	}
	b.Reply(m, "Hey I'm Alive.")
}
