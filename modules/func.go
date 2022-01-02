package modules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/amarnathcjd/yoko/bot"
	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v3"
)

var b = bot.Bot
var BOT_SUDO = []int{}

func parse_message(m *tb.Message) (string, string, []string) {
	if m.IsReply() {
		file_id, file_type := get_file(m.ReplyTo)
		buttons := ""
		if m.ReplyTo.ReplyMarkup != nil {
			buttons = get_reply_markup(m.ReplyTo)
		}
		args := strings.SplitN(m.Text, " ", 3)
		if len(args) == 3 {
			note, name := args[2], args[1]
			note += buttons
			return name, note, []string{file_id, file_type}
		} else if len(args) == 2 {
			if m.ReplyTo.Text != string("") {
				note, name := m.ReplyTo.Text, args[1]
				note += buttons
				return name, note, []string{file_id, file_type}
			} else {
				return args[1], buttons, []string{file_id, file_type}
			}
		} else {
			return "", buttons, []string{file_id, file_type}
		}
	} else {
		args := strings.SplitN(m.Text, " ", 3)
		if len(args) == 3 {
			return args[1], args[2], nil
		} else if len(args) == 2 {
			return args[1], "", nil
		} else {
			return "", "", nil
		}
	}
}

func get_file(m *tb.Message) (string, string) {
	if m.Document != nil {
		return m.Document.FileID, "document"
	} else if m.Photo != nil {
		return m.Photo.FileID, "photo"
	} else if m.Sticker != nil {
		return m.Sticker.FileID, "sticker"
	} else if m.Audio != nil {
		return m.Audio.FileID, "audio"
	} else if m.Voice != nil {
		return m.Voice.FileID, "voice"
	} else if m.Animation != nil {
		return m.Animation.FileID, "animation"
	} else if m.Video != nil {
		return m.Video.FileID, "video"
	} else if m.VideoNote != nil {
		return m.VideoNote.FileID, "videonote"
	} else {
		return "", ""
	}
}

func unparse_message(file interface{}, note string, m *tb.Message) {
	text, buttons := button_parser(note)
	if file != nil && file.(bson.A)[0] != string("") {
		id, f := file.(bson.A)[0].(string), file.(bson.A)[1].(string)
		if f == "document" {
			file := &tb.Document{File: tb.File{FileID: id}, Caption: text}
			b.Reply(m, file, buttons)
		} else if f == "sticker" {
			file := &tb.Sticker{File: tb.File{FileID: id}}
			b.Reply(m, file, buttons)
		} else if f == "photo" {
			file := &tb.Photo{File: tb.File{FileID: id}, Caption: text}
			b.Reply(m, file, buttons)
		} else if f == "audio" {
			file := &tb.Audio{File: tb.File{FileID: id}, Caption: text}
			b.Reply(m, file, buttons)
		} else if f == "voice" {
			file := &tb.Voice{File: tb.File{FileID: id}, Caption: text}
			b.Reply(m, file, buttons)
		} else if f == "video" {
			file := &tb.Video{File: tb.File{FileID: id}, Caption: text}
			b.Reply(m, file, buttons)
		} else if f == "animation" {
			file := &tb.Animation{File: tb.File{FileID: id}, Caption: text}
			b.Reply(m, file, buttons)
		} else if f == "videonote" {
			file := &tb.VideoNote{File: tb.File{FileID: id}}
			b.Reply(m, file, buttons)
		}
	} else {
		b.Reply(m, text, buttons)
	}
}

func get_reply_markup(m *tb.Message) string {
	reply_mark := ""
	mark := m.ReplyMarkup
	for _, y := range mark.InlineKeyboard {
		btn_num := 0
		btn := ""
		for _, x := range y {
			btn_num++
			btn += fmt.Sprintf("\n[%s](buttonurl://%s:same)", x.Text, x.URL)
		}
		if btn_num < 2 {
			btn = strings.Replace(btn, ":same", "", 2)
		}
		reply_mark += btn
	}
	return reply_mark
}

func button_parser(text string) (string, *tb.ReplyMarkup) {
	BTN_URL_REGEX := regexp.MustCompile(`(\[([^\[]+?)\]\((btnurl|buttonurl):(?:/{0,2})(.+?)(:same)?\))`)
	c := BTN_URL_REGEX.FindAllStringSubmatch(text, -1)
	var rows []tb.Row
	btns := &tb.ReplyMarkup{Selective: true}
	for _, m := range c {
		if m[5] != string("") && len(rows) != 0 {
			rows[len(rows)-1] = append(rows[len(rows)-1], btns.URL(m[2], m[4]))
		} else {
			rows = append(rows, btns.Row(btns.URL(m[2], m[4])))
		}
	}
	btns.Inline(rows...)
	note := BTN_URL_REGEX.Split(text, -1)[0]
	return note, btns
}

func Change_info(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		p, _ := b.ChatMemberOf(c.Chat(), c.Sender())
		if p.Role == "member" {
			b.Reply(c.Message(), "You need to be an admin to do this!")
			return nil
		} else if p.Role == "creator" {
			return next(c)
		} else if p.Role == "administrator" {
			if p.Rights.CanChangeInfo {
				return next(c)
			} else {
				b.Reply(c.Message(), "You are missing the following rights to use this command: CanChangeInfo")
				return nil
			}
		}
		return nil
	}
}

func Add_admins(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if c.Message().Private() {
			return next(c)
		}
		p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
		if p.Role == "member" {
			c.Reply("You need to be an admin to do this!")
			return nil
		} else if p.Role == "creator" {
			return next(c)
		} else if p.Role == "administrator" {
			if p.Rights.CanPromoteMembers {
				return next(c)
			} else {
				c.Reply("You are missing the following rights to use this command: CanPromoteUsers")
				return nil
			}
		}
		return nil
	}
}

func Ban_users(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if c.Message().Private() {
			return next(c)
		}
		p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
		if p.Role == "member" {
			c.Reply("You need to be an admin to do this!")
			return nil
		} else if p.Role == "creator" {
			return next(c)
		} else if p.Role == "administrator" {
			if p.Rights.CanRestrictMembers {
				return next(c)
			} else {
				c.Reply("You are missing the following rights to use this command: CanRestrictMembers")
				return nil
			}
		}
		return nil
	}
}

func Pin_messages(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if c.Message().Private() {
			return next(c)
		}
		p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
		if p.Role == "member" {
			c.Reply("You need to be an admin to do this!")
			return nil
		} else if p.Role == "creator" {
			return next(c)
		} else if p.Role == "administrator" {
			if p.Rights.CanPinMessages {
				return next(c)
			} else {
				c.Reply("You are missing the following rights to use this command: CanPinMessages")
				return nil
			}
		}
		return nil
	}
}

func Extract_time(c tb.Context, time_val string) int64 {
	has_suffix := false
	for _, arg := range []string{"m", "h", "d", "w"} {
		if strings.HasSuffix(time_val, arg) {
			has_suffix = true
		}
	}
	if !has_suffix {
		c.Reply(fmt.Sprintf("Invalid time type specified. Expected m,h,d or w, got: %s", string(time_val)))
		return 0
	}
	if !unicode.IsDigit(rune(time_val[0])) {
		c.Reply("Invalid time amount specified.")
	}
	unit := string(time_val[len(time_val)-1])
	time_num, _ := strconv.Atoi(string(time_val[0]))
	bantime := 0
	if unit == "m" {
		bantime = time_num * 60
	} else if unit == "h" {
		bantime = time_num * 60 * 60
	} else if unit == "d" {
		bantime = time_num * 60 * 60 * 24
	} else if unit == "w" {
		bantime = time_num * 60 * 60 * 24 * 7
	}
	if bantime > 31622400 {
		c.Reply("Invalid time specified, temporary actions have to be in between 1 minute and 366 days.")
		return 0
	}
	return int64(bantime)
}
