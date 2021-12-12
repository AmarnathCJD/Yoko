package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v3"
)

func save(c tb.Context) error {
	m := c.Message()
	name, note, file := parse_message(m)
	if name == string("") {
		b.Reply(m, "You need to give the note a name!")
		return nil
	} else if note == string("") && file == nil {
		b.Reply(m, "You need to give the note some content!")
		return nil
	}
	b.Reply(m, fmt.Sprintf("✨ Saved note <code>%s</code>", name))
	save_note(m.Chat.ID, name, note, file)
	get_notes(m.Chat.ID)
	return nil
}

func rgtest(c tb.Context) error {
	m := c.Message()
	fmt.Println("hello")
	b.Reply(m, m.Payload)
	return nil
}
