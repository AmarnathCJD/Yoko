package modules

import (
	"fmt"
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

var LOCK_TYPES = []string{"all", "album", "anonchannel", "audio", "bot", "button", "command", "comment", "contact", "document", "email", "emojigame", "forward", "forwardbot", "forwardchannel", "forwarduser", "game", "gif", "inline", "invitelink", "location", "phone", "photo", "poll", "rtl", "sticker", "text", "url", "video", "videonote", "voice"}

func Lock(c tb.Context) error {
	m := c.Message()
	if m.Payload == string("") {
		b.Reply(m, "You haven't specified a type to lock.")
		return nil
	}
	args := strings.Split(m.Payload, " ")
	to_lock := make([]string, 0)
	if stringInSlice("all", args) {
		b.Reply(m, "Locked <code>all</code>.")
		db.Lock_item(m.Chat.ID, LOCK_TYPES)
		return nil
	}
	locked_msg := ""
	for _, lock := range args {
		if stringInSlice(lock, LOCK_TYPES) {
			to_lock = append(to_lock, lock)
			locked_msg += fmt.Sprintf("\n- <code>%s</code>", lock)
		}
	}
	if len(to_lock) == 0 {
		b.Reply(m, fmt.Sprintf("✨ Unknown lock types:- %s\nCheck /locktypes !", m.Payload))
		return nil
	}
	if len(to_lock) == 1 {
		b.Reply(m, fmt.Sprintf("Locked <code>%s</code>.", to_lock[0]))
		return nil
	}
	b.Reply(m, "Locked "+locked_msg)
	db.Lock_item(m.Chat.ID, to_lock)
	return nil
}

func Locktypes(c tb.Context) error {
	m := c.Message()
	lock_types := ""
	for _, lock := range LOCK_TYPES {
		lock_types += fmt.Sprintf("\n<b>~</b> %s", lock)
	}
	b.Reply(m, "The available lock types are:"+lock_types)
	return nil
}

func Check_locks(c tb.Context) error {
	m := c.Message()
	locked := "✨ Chat LockSettings"
	lock_c := db.Get_locks(m.Chat.ID)
	for _, lock := range LOCK_TYPES {
		if db.IsTrue(lock, lock_c) {
			locked += fmt.Sprintf("\n<b>-&gt; %s:</b> true", lock)
		} else {
			locked += fmt.Sprintf("\n<b>-&gt; %s:</b> false", lock)
		}
	}
	b.Reply(m, locked)
	return nil
}

func Unlock(c tb.Context) error {
	m := c.Message()
	if m.Payload == string("") {
		b.Reply(m, "You haven't specified a type to unlock.")
		return nil
	}
	args := strings.Split(m.Payload, " ")
	to_unlock := make([]string, 0)
	if stringInSlice("all", args) {
		b.Reply(m, "Unlocked all.")
		db.Unlock_item(m.Chat.ID, LOCK_TYPES)
		return nil
	}
	unlocked_msg := ""
	for _, lock := range args {
		if stringInSlice(lock, LOCK_TYPES) {
			to_unlock = append(to_unlock, lock)
			unlocked_msg += fmt.Sprintf("\n- <code>%s</code>", lock)
		}
	}
	if len(to_unlock) == 0 {
		b.Reply(m, fmt.Sprintf("✨ Unknown lock types:- %s\nCheck /locktypes !", m.Payload))
		return nil
	}
	if len(to_unlock) == 1 {
		b.Reply(m, fmt.Sprintf("Unlocked <code>%s</code>.", to_unlock[0]))
		db.Unlock_item(m.Chat.ID, to_unlock)
		return nil
	}
	b.Reply(m, "Unlocked "+unlocked_msg)
	db.Unlock_item(m.Chat.ID, to_unlock)
	return nil
}

func anonchannel(c tb.Context) error {
	fmt.Println("hii")
	if c.Chat().Type != "channel" {
		return nil
	}
	m := c.Message()
	lock_c := db.Get_locks(m.Chat.ID)
	x := false
	for _, y := range lock_c {
		if y == "anonchannel" {
			x = true
		}
	}
	if !x {
		return nil
	}
	fmt.Println("CHannel Msg", c.Sender())
	return nil
}
