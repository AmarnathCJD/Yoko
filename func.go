package main

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v3"
        "go.mongodb.org/mongo-driver/bson"
)

func parse_message(m *tb.Message) (string, string, []string) {
	if m.IsReply() {
		file_id, file_type := get_file(m.ReplyTo)
                buttons := ""
                if m.ReplyTo.ReplyMarkup != nil{
		  buttons = get_reply_markup(m.ReplyTo)
                }
                fmt.Println("2")
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

func make_file(id string, f string, note string) interface{}{
 if f == "document"{
    file := tb.Document{File: tb.File{FileID: id}, Caption: note}
    return file
 } else if f == "sticker"{
    file := tb.Sticker{File: tb.File{FileID: id}}
    return file
 }
 return nil
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

func unparse_message(note bson.A, media bson.A) {
 fmt.Println(note)
 file := make_file(media[0].(string), media[1].(string), "hi")
 fmt.Println(file)
}
