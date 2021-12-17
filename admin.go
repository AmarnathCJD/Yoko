package main

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v3"
)

func add_admins(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		p, _ := b.ChatMemberOf(c.Chat(), c.Sender())
		if p.Role == "member" {
			b.Reply(c.Message(), "You need to be an admin to do this!")
			return nil
		} else if p.Role == "creator" {
			return next(c)
		} else if p.Role == "administrator" {
			if p.Rights.CanPromoteMembers {
				return next(c)
			} else {
				b.Reply(c.Message(), "You are missing the following rights to use this command: CanPromoteUsers")
				return nil
			}
		}
		return nil
	}
}

func promote(c tb.Context) error {
	user, xtra := get_user(c.Message())
	if user == nil {
		return nil
	} else if user.ID == 5050904599 {
		b.Reply(c.Message(), "Pffff, I wish I could just promote myself.")
		return nil
	}
	title := "Admin"
	if xtra != string("") {
		title = xtra
	}
	arg := strings.SplitN(c.Message().Text, " ", 2)[0]
	if arg == "/promote" {
		err := b.Promote(c.Chat(), &tb.ChatMember{
			Rights: tb.Rights{CanRestrictMembers: true, CanPinMessages: true, CanChangeInfo: false, CanDeleteMessages: true, CanInviteUsers: true},
			User:   user,
			Title:  title,
		})
		if err == nil {
			b.Reply(c.Message(), "✨ Successfully superpromoted! ~")
		} else if err.Error() == string("telegram: can't remove chat owner (400)") {
			b.Reply(c.Message(), "I would love to promote the chat creator, but... well, they already have all the power.")
		} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
			b.Reply(c.Message(), "Failed to promote, "+"Make sure I'm admin and can appoint new admins.")
		} else if err.Error() == "telegram unknown: Bad Request: CHAT_ADMIN_REQUIRED (400)" {
			b.Reply(c.Message(), "This user has been already promoted by someone otherthan me; I can't change their permissions!")
		} else if err.Error() == "telegram unknown: Bad Request: USER_PRIVACY_RESTRICTED (400)" {
			b.Reply(c.Message(), "Failed to promote, use was not found in this chat.")
		} else {
			b.Reply(c.Message(), "Failed to promote, "+fmt.Sprint(err.Error()))
		}
	} else if arg == "/superpromote" {
		err := b.Promote(c.Chat(), &tb.ChatMember{
			Rights: tb.Rights{CanRestrictMembers: true, CanPinMessages: true, CanChangeInfo: true, CanDeleteMessages: true, CanInviteUsers: true, CanPromoteMembers: true, CanManageVoiceChats: true},
			User:   user,
			Title:  title,
		})
		if err == nil {
			b.Reply(c.Message(), "✨ Successfully superpromoted! ~")
		} else if err.Error() == string("telegram: can't remove chat owner (400)") {
			b.Reply(c.Message(), "I would love to promote the chat creator, but... well, they already have all the power.")
		} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
			b.Reply(c.Message(), "Failed to promote, "+"Make sure I'm admin and can appoint new admins.")
		} else if err.Error() == "telegram unknown: Bad Request: CHAT_ADMIN_REQUIRED (400)" {
			b.Reply(c.Message(), "This user has been already promoted by someone otherthan me; I can't change their permissions!")
		} else if err.Error() == "telegram unknown: Bad Request: USER_PRIVACY_RESTRICTED (400)" {
			b.Reply(c.Message(), "Failed to promote, user was not found in this chat.")
		} else {
			b.Reply(c.Message(), "Failed to promote, "+fmt.Sprint(err.Error()))
		}
	}
	return nil
}

func demote(c tb.Context) error {
	user, _ := get_user(c.Message())
	if user == nil {
		return nil
	} else if user.ID == 5050904599 {
		b.Reply(c.Message(), "I am not going to demote myself.")
		return nil
	}
	err := b.Promote(c.Chat(), &tb.ChatMember{
		Rights: tb.NoRestrictions(),
		User:   user,
	})
	if err == nil {
		b.Reply(c.Message(), "✨ Successfully demoted! ~")
	} else if err.Error() == "telegram: can't remove chat owner (400)" {
		b.Reply(c.Message(), "I don't really feel like staging a mutiny today, I think the chat owner deserves to stay an admin.")
	} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
		b.Reply(c.Message(), "Failed to demote, "+"Make sure I'm admin and can appoint new admins.")
	} else if err.Error() == "telegram unknown: Bad Request: CHAT_ADMIN_REQUIRED (400)" {
		b.Reply(c.Message(), "This user has been already promoted by someone otherthan me; I can't change their permissions!")
	} else {
		b.Reply(c.Message(), "Failed to demote, "+fmt.Sprint(err.Error()))
	}
	return nil
}
