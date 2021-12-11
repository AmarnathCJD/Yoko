package main

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

var LOCK_TYPES = []string{"all", "album", "anonchannel", "audio", "bot", "button", "command", "comment", "contact", "document", "email", "emojigame", "forward", "forwardbot", "forwardchannel", "forwarduser", "game", "gif", "inline", "invitelink", "location", "phone", "photo", "poll", "rtl", "sticker", "text", "url", "video", "videonote", "voice"}

func lock(m *tb.Message) {
	if m.Payload == string("") {
		b.Reply(m, "You haven't specified a type to lock.")
		return
	}
	args := strings.Split(m.Payload, " ")
	to_lock := make([]string, 0)
	if stringInSlice("all", args) {
		b.Reply(m, "Locked <code>all</code>.")
		lock_item(m.Chat.ID, LOCK_TYPES)
		return
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
		return
	}
	if len(to_lock) == 1 {
		b.Reply(m, fmt.Sprintf("Locked <code>%s</code>.", to_lock[0]))
		return
	}
	b.Reply(m, "Locked "+locked_msg)
	lock_item(m.Chat.ID, to_lock)
}

func locktypes(m *tb.Message) {
	lock_types := ""
	for _, lock := range LOCK_TYPES {
		lock_types += fmt.Sprintf("\n<b>~</b> %s", lock)
	}
	b.Reply(m, "The available lock types are:"+lock_types)
}

func check_locks(m *tb.Message) {
	locked := "✨ Chat LockSettings"
	lock_c := get_locks(m.Chat.ID)
	for _, lock := range LOCK_TYPES {
		if isTrue(lock, lock_c) {
			locked += fmt.Sprintf("\n<b>-&gt; %s:</b> true", lock)
		} else {
			locked += fmt.Sprintf("\n<b>-&gt; %s:</b> false", lock)
		}
	}
	b.Reply(m, locked)
}

func unlock(m *tb.Message) {
	if m.Payload == string("") {
		b.Reply(m, "You haven't specified a type to unlock.")
		return
	}
	args := strings.Split(m.Payload, " ")
	to_unlock := make([]string, 0)
	if stringInSlice("all", args) {
		b.Reply(m, "Unlocked all.")
		unlock_item(m.Chat.ID, LOCK_TYPES)
		return
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
		return
	}
	if len(to_unlock) == 1 {
		b.Reply(m, fmt.Sprintf("Unlocked <code>%s</code>.", to_unlock[0]))
		unlock_item(m.Chat.ID, to_unlock)
		return
	}
	b.Reply(m, "Unlocked "+unlocked_msg)
	unlock_item(m.Chat.ID, to_unlock)
}
