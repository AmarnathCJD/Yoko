package main

import (
 tb "gopkg.in/tucnak/telebot.v2"
 "fmt"
 "strings"
)

var LOCK_TYPES = [31]string{"all","album","audio", "bot", "button",  "command", "comment", "contact", "document",  "email",  "emojigame",  "forward",  "forwardbot", "forwardchannel","forwarduser","game", "gif",  "inline", "invitelink",  "location",  "phone",  "photo", "poll", "rtl",  "sticker", "text", "url", "video", "videonote",  "voice", "anonchannel"}

func lock(m *tb.Message){
 if m.Payload == string(""){
    b.Reply(m, "You haven't specified a type to lock.")
 }
 lk = strings.Split(m.Payload, " ")
 for lock := range lk {
     if stringInSlice(lock, LOCK_TYPE){
        fmt.Println(lock)
     }
 }
}
