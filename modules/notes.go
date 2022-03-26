package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func Save(c tb.Context) error {
	Message := ParseMessage(c)
	fmt.Println(Message)
	if Message.Name == string("") {
		c.Reply("You need to give the note a name!")
		return nil
	} else if Message.Text == string("") && Message.File.FileID == string("") {
		c.Reply("You need to give the note some content!")
		return nil
	}
	c.Reply(fmt.Sprintf("✨ Saved note <code>%s</code>", Message.Name))
	db.SaveNote(c.Chat().ID, Message)
	return nil
}

func AllNotes(c tb.Context) error {
	var Title string
	if c.Message().Private() {
		Title = c.Sender().FirstName
	} else {
		Title = c.Chat().Title
	}
	m := c.Message()
	notes := db.GetNotes(c.Chat().ID)
	if notes == nil {
		return c.Reply(fmt.Sprintf("There are no notes in %s!", Title))
	}
	mode := db.PnoteSettings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/%s?start=allnotes_%d", BOT_USERNAME, c.Chat().ID))))
		c.Reply("Tap here to view all notes in this chat.", menu)
		return nil
	}
	note := fmt.Sprintf("Notes in <b>%s</b>", c.Chat().Title)
	for _, x := range notes {
		note += fmt.Sprintf("\n<b>-</b> <code>#%s</code>", x.Name)
	}
	b.Reply(m, note+"\nYou can retrieve these notes by using <code>/get notename</code>")
	return nil
}

func Getnote(c tb.Context) error {
	note := db.GetNote(c.Chat().ID, c.Message().Payload)
	if note.Name == string("") {
		c.Reply("No note found!")
		return nil
	}
	mode := db.PnoteSettings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/missmikabot?start=notes_%d_%s", c.Chat().ID, c.Message().Payload))))
		c.Reply(fmt.Sprintf("Tap here to view '%s' in your private chat.", c.Message().Payload), menu)
		return nil
	}
	text, p := ParseString(note.Text, c)
	text, btns := button_parser(text)
	if note.File.FileID != string("") {
		f := GetSendable(note)
		_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
		if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}
	} else {
		if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}

	}
	return nil
}

func HashNote(c tb.Context) error {
	args := strings.SplitN(c.Message().Text, "#", 2)
	note := db.GetNote(c.Chat().ID, args[1])
	if note.Name == string("") {
		return nil
	}
	mode := db.PnoteSettings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/%s?start=notes_%d_%s", BOT_USERNAME, c.Chat().ID, args[1]))))
		c.Reply(fmt.Sprintf("Tap here to view '%s' in your private chat.", args[1]), menu)
		return nil
	}
	text, p := ParseString(note.Text, c)
	text, btns := button_parser(text)
	if note.File.FileID != string("") {
		f := GetSendable(note)
		_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
		if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}
	} else {
		if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
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
	r := db.DelNote(c.Chat().ID, c.Message().Payload)
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
	db.DelAllNotes(c.Chat().ID)
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
		mode := db.PnoteSettings(c.Chat().ID)
		if !mode {
			c.Reply("Your notes are currently being sent in the group.")
		} else {
			c.Reply("Your notes are currently being sent in private. Mika will send a small note with a button which redirects to a private chat.")
		}
	} else {
		payload := c.Message().Payload
		if stringInSlice(payload, []string{"enable", "yes", "on", "y"}) {
			c.Reply("Mika will now send a message to your chat with a button redirecting to PM, where the user will receive the note.")
			db.SetPnote(c.Chat().ID, true)
		} else if stringInSlice(payload, []string{"off", "disable", "no"}) {
			c.Reply("Mika will now send notes straight to the group.")
			db.SetPnote(c.Chat().ID, false)
		} else {
			c.Reply(fmt.Sprintf("failed to get boolean value from input: expected one of y/yes/on or n/no/off; got: %s", payload))
		}
	}
	return nil
}

func PrivateStartNote(c tb.Context) error {
	args := strings.SplitN(c.Message().Payload, "_", 3)
	chat_id, _ := strconv.Atoi(args[1])
	note := db.GetNote(int64(chat_id), args[2])
	if note.Name == "" {
		return c.Reply("This note was not found ~")
	}
	text, p := ParseString(note.Text, c)

	if note.File.FileID != "" {
		f := GetSendable(note)
		_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
		if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
			return err
		}
	} else {

		if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
			return c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
		}

	}
	return nil
}

func PrivateStartNotes(c tb.Context) error {
	args := strings.SplitN(c.Message().Payload, "_", 2)
	chat_id, _ := strconv.Atoi(args[1])
	notes := db.GetNotes(int64(chat_id))
	if notes == nil {
		return c.Reply("There are no notes in that chat. ~")
	}
	out := fmt.Sprintf("<b>Notes in %s:</b>", args[1])
	for _, x := range notes {
		out += fmt.Sprintf("\n<b>-&gt;</b> <a href='t.me/missmikabot?start=notes_%s_%s'>%s</a>", args[1], x.Name, x.Name)
	}
	return c.Reply(out)
}
