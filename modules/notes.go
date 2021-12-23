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
	if c.Message().Payload == string("") {
		c.Reply("Which note should I remove?")
		return nil
	}
	r := db.Del_note(c.Chat().ID, c.Message().Payload)
	if !r {
		c.Reply("You haven't saved any notes with this name yet!")
		return nil
	} else {
		c.Reply(fmt.Sprintf("Note <code>%s</code> was deleted", c.Message().Payload))
	}
	return nil
}

func clear_all(c tb.Context) error {
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == "admnistrator" {
		if !p.Rights.CanChangeInfo {
			c.Reply("You are missing the following rights to use this command: CanChangeInfo")
			return nil
		}
	} else if p.Role == "member" {
		c.Reply("You need to be an admin to do this!")
		return nil
	} else if p.Role != "creator" {
		c.Reply("You need to be chat creator to do this!")
		return nil
	}
	btns := &tb.ReplyMarkup{}
	btns.Inline(btns.Row(btns.Data("Delete all notes", "delall_notes")), btns.Row(btns.Data("Cancel", "cancel_delall")))
	c.Reply(fmt.Sprintf("Are you sure you would like to clear <b>ALL</b> notes in %s? This action cannot be undone.", c.Chat().Title), btns)
	return nil
}

func private_notes(c tb.Context) error {
	return nil
}
