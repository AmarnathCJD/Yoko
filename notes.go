package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func save(m *tb.Message) {
	name, note, file := parse_message(m)
	if name == string("") {
		b.Reply(m, "You need to give the note a name!")
		return
	} else if note == string("") && file == nil {
		b.Reply(m, "You need to give the note some content!")
		return
	}
	b.Reply(m, fmt.Sprintf("✨ Saved note <code>%s</code>", name))
	save_note(m.Chat.ID, name, note, file)
	get_notes(m.Chat.ID)
}

func rgtest(m *tb.Message) {
	fmt.Println("hello")
	b.Reply(m, m.Payload)
}
