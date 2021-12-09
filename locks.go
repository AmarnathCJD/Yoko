package main

import (
 tb "gopkg.in/tucnak/telebot.v2"
 "fmt"
 "strings"
)

var LOCK_TYPES = []string{"all","album","anonchannel","audio", "bot", "button",  "command", "comment", "contact", "document",  "email",  "emojigame",  "forward",  "forwardbot", "forwardchannel","forwarduser","game", "gif",  "inline", "invitelink",  "location",  "phone",  "photo", "poll", "rtl",  "sticker", "text", "url", "video", "videonote",  "voice"}

func lock(m *tb.Message){
 if m.Payload == string(""){
    b.Reply(m, "You haven't specified a type to lock.")
    return
 }
 args := strings.Split(m.Payload, " ")
 to_lock := make([]string, 0)
 for _, lock := range args {
     if stringInSlice(lock, LOCK_TYPES){
        to_lock = append(to_lock, lock)
     }
 }
 if len(to_lock) == 0{
    b.Reply(m, fmt.Sprintf("✨ Unknown lock types:- %s\nCheck <code>/locktypes</code>!", m.Payload))
    return
 }
 b.Reply(m, "Locked " + strings.Join(to_lock, "\n- "))
}

func locktypes(m *tb.Message){
 lock_types := strings.Join(LOCK_TYPES,"\n<b>~</b> ")
 b.Reply(m, "The available lock types are:" + lock_types)
}
