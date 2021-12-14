package main

import (
	"fmt"
	"strings"
        "go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v3"
)

func parse_message(m *tb.Message) (string, string, []string) {
	if m.IsReply() {
		file_id, file_type := get_file(m.ReplyTo)
                buttons := ""
                if m.ReplyTo.ReplyMarkup != nil{
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

func unparse_message(file interface{}, note string, m *tb.Message){
 if len(file.(bson.A)) != 0{
  id, f := file.(bson.A)[0].(string), file.(bson.A)[1].(string)
  if f == "document"{
    file := &tb.Document{File: tb.File{FileID: id}, Caption: note}
    _, err := b.Reply(m, file)
    fmt.Println(err)
  } else if f == "sticker"{
    file := &tb.Sticker{File: tb.File{FileID: id}, }
    b.Reply(m, file)
  } else if f == "photo"{
    file := &tb.Photo{File: tb.File{FileID: id}, Caption: note}
    _, err := b.Reply(m, file)
    fmt.Println(err)
  } else if f == "audio"{
    file := &tb.Audio{File: tb.File{FileID: id}, Caption: note}
    b.Reply(m, file)
  } else if f == "voice"{
    file := &tb.Voice{File: tb.File{FileID: id}, Caption: note}
    b.Reply(m, file)
  } else if f == "video"{
    file := &tb.Video{File: tb.File{FileID: id}, Caption: note}
    b.Reply(m, file)
  } else if f == "animation"{
    file := &tb.Animation{File: tb.File{FileID: id}, Caption: note}
    b.Reply(m, file)
  } else if f == "videonote"{
    file := &tb.VideoNote{File: tb.File{FileID: id}, }
    b.Reply(m, file)
  }
 } else if note != string(""){
   b.Reply(m, note)
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

var BTN_URL_REGEX = regexp.MustCompile(`(\\[([^\\[]+?)\\]\\((btnurl|buttonurl):(?:/{0,2})(.+?)(:same)?\\))`)

func button_parser(){
 rg := "Hi[Google](buttonurl://google.com) [Yahoo](buttonurl://google.com:same)"
 c := BTN_URL_REGEX.FindAllStringSubmatch(rg, -1)
 fmt.Println(c)
}
