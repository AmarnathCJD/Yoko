package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v3"
        "go.mongodb.org/mongo-driver/bson"
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

func all_notes(c tb.Context) error {
 m := c.Message()
 notes := get_notes(c.Chat().ID)
 if notes == nil{
    b.Reply(m, fmt.Sprintf("There are no notes in %s!", m.Chat.Title))
    return nil
 }
 note := fmt.Sprintf("Notes in <b>%s</b>", c.Chat().Title)
 for _, x := range notes{
  note += fmt.Sprintf("\n<b>-&gt;</b> <code>%s</code>", x.(bson.M)["name"].(string))
 }
 b.Reply(m, note + "\nYou can retrieve these notes by using <code>/get notename</code>")
 return nil
}

func gnote(c tb.Context) error {
 m := c.Message()
 note := get_note(m.Chat.ID, m.Payload)
 if note == nil{
    b.Reply(m, "No note found!")
    return nil
 }
 f := unparse_message(fmt.Sprint(note["note"]), note["file"].(bson.A))
 _, err := b.Reply(m, f.(tb.Sticker))
 fmt.Println(err)
 return nil
}
 
