package modules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v3"
)

func Save(c tb.Context) error {
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
	db.Save_note(m.Chat.ID, name, note, file)
	return nil
}

func All_notes(c tb.Context) error {
	m := c.Message()
	notes := db.Get_notes(c.Chat().ID)
	if notes == nil {
		b.Reply(m, fmt.Sprintf("There are no notes in %s!", m.Chat.Title))
		return nil
	}
	note := fmt.Sprintf("Notes in <b>%s</b>", c.Chat().Title)
	for _, x := range notes {
		note += fmt.Sprintf("\n<b>-&gt;</b> <code>%s</code>", x.(bson.M)["name"].(string))
	}
	b.Reply(m, note+"\nYou can retrieve these notes by using <code>/get notename</code>")
	return nil
}

func Gnote(c tb.Context) error {
	m := c.Message()
	note := db.Get_note(m.Chat.ID, m.Payload)
	if note == nil {
		b.Reply(m, "No note found!")
		return nil
	}
	unparse_message(note["file"], note["note"].(string), m)
	return nil
}

func Hash_regex(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		match, _ := regexp.MatchString("\\#(\\S+)", c.Message().Text)
		if match {
			return next(c)
		}
		return nil
	}
}

func Hash_note(c tb.Context) error {
	args := strings.SplitN(c.Message().Text, "#", 2)
	note := db.Get_note(c.Message().Chat.ID, args[1])
	if note == nil {
		return nil
	}
	unparse_message(note["file"], note["note"].(string), c.Message())
	return nil
}

func clear_note(c tb.Context) error {
	return nil
}
