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
	mode := db.Pnote_settings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/yoko_robot?start=notes_%d_%s", c.Chat().ID, m.Payload))))
		c.Reply(fmt.Sprintf("Tap here to view '%s' in your private chat.", m.Payload), menu)
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

func EchoPerReg(next tele.HandlerFunc) tele.HandlerFunc {
         return func(c tele.Context) error {
              exp := regexp.MustCompile("^اکو.*")
              if exp.MatchString(c.Message().Text) {
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
	mode := db.Pnote_settings(c.Chat().ID)
	if mode {
		menu := &tb.ReplyMarkup{}
		menu.Inline(menu.Row(menu.URL("Click me!", fmt.Sprintf("t.me/yoko_robot?start=notes_%d_%s", c.Chat().ID, args[1]))))
		c.Reply(fmt.Sprintf("Tap here to view '%s' in your private chat.", args[1]), menu)
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
			c.Reply("Your notes are currently being sent in private. Aiko will send a small note with a button which redirects to a private chat.")
		}
	} else {
		payload := c.Message().Payload
		if stringInSlice(payload, []string{"enable", "yes", "on", "y"}) {
			c.Reply("Aiko will now send a message to your chat with a button redirecting to PM, where the user will receive the note.")
			db.Set_pnote(c.Chat().ID, true)
		} else if stringInSlice(payload, []string{"off", "disable", "no"}) {
			c.Reply("Aiko will now send notes straight to the group.")
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
	unparse_message(note["file"], note["note"].(string), c.Message())
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
		out += fmt.Sprintf("\n<b>-&gt;</b> <a href='t.me/yoko_robot?start=notes_%s_%s'>%s</a>", args[1], x.(bson.M)["name"].(string), x.(bson.M)["name"].(string))
	}
	c.Reply(out)
}

var notes_help = `✨ Here is the help for <b>Notes:<b>
<b>Command for Members</b>
<b>-></b> <code>/get notename</code>: get the note with this notename
<b>-</b> <code>#notename</code>: same as /get
<b>-></b> <code>/notes<code>: list all saved notes in this chat
<b>Command for Admins</b>
<b>-></b> <code>/save notename notedata</code>: saves notedata as a note with name notename, reply to a message or document to save it
<b>-></b> <code>/clear notename</code>: clear note with this name
<b>-></b> <code>/privatenote on/yes/off/no</code>: whether or not to send the note in PM. Write del besides on/off to delete hashtag message on group.
<b>Note</b>
<b>-</b> Only admins can use This module
<b>-</b> To save a document (like photo, audio, etc.), reply to a document or media then type /save
<b>-</b> Need help for parsing text? Check /markdownhelp
Save data for future users with notes!
Notes are great to save random tidbits of information; a phone number, a nice gif, a funny picture - anything!
Also you can save a text/document with buttons, you can even save it in here.`
