package modules

import (
	"fmt"
	"regexp"
	"strconv"
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
	mode := db.Pnote_settings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/yoko_robot?start=allnotes_%d", c.Chat().ID))))
		c.Reply("Tap here to view all notes in this chat.", menu)
		return nil
	}
	note := fmt.Sprintf("Notes in <b>%s</b>", c.Chat().Title)
	for _, x := range notes {
		note += fmt.Sprintf("\n<b>-</b> <code>#%s</code>", x.(bson.M)["name"].(string))
	}
	b.Reply(m, note+"\nYou can retrieve these notes by using <code>/get notename</code>")
	return nil
}

func Gnote(c tb.Context) error {
	note := db.Get_note(c.Chat().ID, c.Message().Payload)
	if note == nil {
		c.Reply("No note found!")
		return nil
	}
	mode := db.Pnote_settings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/missmikabot?start=notes_%d_%s", c.Chat().ID, c.Message().Payload))))
		c.Reply(fmt.Sprintf("Tap here to view '%s' in your private chat.", c.Message().Payload), menu)
		return nil
	}

	text, p := ParseString(note["note"].(string), c)

	if note["file"] != nil && len(note["file"].(bson.A)) != 0 && note["file"].(bson.A)[0] != string("") {
		f := GetFile(note["file"].(bson.A), text)
		_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
		if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}
	} else {

		if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}

	}
	return nil
}

func OnTextHandler(c tb.Context) error {
	FLOOD_EV(c)
	if err, ok := FilterEvent(c); ok {
		return err
	}
	match, _ := regexp.MatchString("\\#(\\S+)", c.Message().Text)
	if match {
		Hash_note(c)
		return nil
	}
	if afk := AFK(c); afk {
		return nil
	}
	Chat_bot(c)
	return nil
}

func Hash_note(c tb.Context) error {
        fmt.Println("#hash")
	args := strings.SplitN(c.Message().Text, "#", 2)
	note := db.Get_note(c.Message().Chat.ID, args[1])
	if note == nil {
		return nil
	}
	mode := db.Pnote_settings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/missmikabot?start=notes_%d_%s", c.Chat().ID, args[1]))))
		c.Reply(fmt.Sprintf("Tap here to view '%s' in your private chat.", args[1]), menu)
		return nil
	}
	text, p := ParseString(note["note"].(string), c)

	if note["file"] != nil && len(note["file"].(bson.A)) != 0 && note["file"].(bson.A)[0] != string("") {
		f := GetFile(note["file"].(bson.A), text)
		_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
		if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}
	} else {

		if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}

	}
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

var (
	btns      = &tb.ReplyMarkup{}
	dall_btn  = btns.Data("Delete all notes", "delall_notes")
	cdall_btn = btns.Data("Cancel", "cancel_delall")
)

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
	btns.Inline(btns.Row(dall_btn), btns.Row(cdall_btn))
	c.Reply(fmt.Sprintf("Are you sure you would like to clear <b>ALL</b> notes in %s? This action cannot be undone.", c.Chat().Title), btns)
	return nil
}

func del_all_notes_cb(c tb.Context) error {
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == "member" {
		c.Respond(&tb.CallbackResponse{Text: "you should be an admin to do this!", ShowAlert: true})
		return nil
	} else if p.Role == "administrator" {
		c.Reply(fmt.Sprintf("You need to be the chat owner of %s to do this.", c.Chat().Title))
		c.Delete()
	}
	c.Edit("Deleted all chat notes.")
	db.Del_all_notes(c.Chat().ID)
	return nil
}

func cancel_delall_cb(c tb.Context) error {
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == "member" {
		c.Respond(&tb.CallbackResponse{Text: "you should be an admin to do this!", ShowAlert: true})
		return nil
	} else if p.Role == "administrator" {
		c.Reply(fmt.Sprintf("You need to be the chat owner of %s to do this.", c.Chat().Title))
		c.Delete()
	}
	c.Edit("Clearing of all notes has been cancelled.")
	return nil
}

func private_notes(c tb.Context) error {
	if c.Message().Payload == string("") {
		mode := db.Pnote_settings(c.Chat().ID)
		if !mode {
			c.Reply("Your notes are currently being sent in the group.")
		} else {
			c.Reply("Your notes are currently being sent in private. Mika will send a small note with a button which redirects to a private chat.")
		}
	} else {
		payload := c.Message().Payload
		if stringInSlice(payload, []string{"enable", "yes", "on", "y"}) {
			c.Reply("Mika will now send a message to your chat with a button redirecting to PM, where the user will receive the note.")
			db.Set_pnote(c.Chat().ID, true)
		} else if stringInSlice(payload, []string{"off", "disable", "no"}) {
			c.Reply("Mika will now send notes straight to the group.")
			db.Set_pnote(c.Chat().ID, false)
		} else {
			c.Reply(fmt.Sprintf("failed to get boolean value from input: expected one of y/yes/on or n/no/off; got: %s", payload))
		}
	}
	return nil
}

func private_start_note(c tb.Context) {
	args := strings.SplitN(c.Message().Payload, "_", 3)
	chat_id, _ := strconv.Atoi(args[1])
	note := db.Get_note(int64(chat_id), args[2])
	if note == nil {
		c.Reply("This note was not found ~")
		return
	}
	text, p := ParseString(note["note"].(string), c)

	if note["file"] != nil && len(note["file"].(bson.A)) != 0 && note["file"].(bson.A)[0] != string("") {
		f := GetFile(note["file"].(bson.A), text)
		_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
		if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}
	} else {

		if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}

	}
}

func private_startallnotes(c tb.Context) {
	args := strings.SplitN(c.Message().Payload, "_", 2)
	chat_id, _ := strconv.Atoi(args[1])
	notes := db.Get_notes(int64(chat_id))
	if notes == nil {
		c.Reply("There are no notes in that chat. ~")
		return
	}
	out := fmt.Sprintf("<b>Notes in %s:</b>", args[1])
	for _, x := range notes {
		out += fmt.Sprintf("\n<b>-&gt;</b> <a href='t.me/missmikabot?start=notes_%s_%s'>%s</a>", args[1], x.(bson.M)["name"].(string), x.(bson.M)["name"].(string))
	}
	c.Reply(out)
}
