package modules

import (
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/telebot.v3"
)

func Promote(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	user, xtra := GetUser(c)
	if user.ID == 0 {
		return nil
	} else if user.ID == BOT_ID {
		c.Reply("Pffff, I wish I could just promote myself.")
		return nil
	}
	arg := strings.SplitN(c.Message().Text, " ", 2)[0][1:]
	if arg == "promote" {
		err := c.Bot().Promote(c.Chat(), &tb.ChatMember{
			Rights: tb.Rights{CanRestrictMembers: true, CanPinMessages: true, CanChangeInfo: false, CanDeleteMessages: true, CanInviteUsers: true},
			User:   user.User(),
		})
		if err == nil {
			c.Reply("✨ Successfully promoted! ~")
		} else if err.Error() == string("telegram: can't remove chat owner (400)") {
			c.Reply("I would love to promote the chat creator, but... well, they already have all the power.")
		} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
			c.Reply("Failed to promote, " + "Make sure I'm admin and can appoint new admins.")
		} else if err.Error() == "telegram unknown: Bad Request: CHAT_ADMIN_REQUIRED (400)" {
			c.Reply("This user has been already promoted by someone otherthan me; I can't change their permissions!")
		} else if err.Error() == "telegram unknown: Bad Request: USER_PRIVACY_RESTRICTED (400)" {
			c.Reply("Failed to promote, use was not found in this chat.")
		} else {
			c.Reply("Failed to promote, " + fmt.Sprint(err.Error()))
		}
	} else if arg == "superpromote" {
		err := c.Bot().Promote(c.Chat(), &tb.ChatMember{
			Rights: tb.Rights{CanRestrictMembers: true, CanPinMessages: true, CanChangeInfo: true, CanDeleteMessages: true, CanInviteUsers: true, CanPromoteMembers: true, CanManageVoiceChats: true},
			User:   user.User(),
		})
		if err == nil {
			c.Reply("✨ Successfully superpromoted! ~")
		} else if err.Error() == string("telegram: can't remove chat owner (400)") {
			c.Reply("I would love to promote the chat creator, but... well, they already have all the power.")
		} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
			c.Reply("Failed to promote, " + "Make sure I'm admin and can appoint new admins.")
		} else if err.Error() == "telegram unknown: Bad Request: CHAT_ADMIN_REQUIRED (400)" {
			c.Reply("This user has been already promoted by someone otherthan me; I can't change their permissions!")
		} else if err.Error() == "telegram unknown: Bad Request: USER_PRIVACY_RESTRICTED (400)" {
			c.Reply("Failed to promote, user was not found in this chat.")
		} else {
			c.Reply("Failed to promote, " + fmt.Sprint(err.Error()))
		}
	}
	if xtra != string("") {
		EditTitle(c, user.User(), xtra, true)
	}
	return nil
}

func Demote(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	user, _ := GetUser(c)
	if user.ID == 0 {
		return nil
	} else if user.ID == BOT_ID {
		c.Reply("I am not going to demote myself.")
		return nil
	}
	err := c.Bot().Promote(c.Chat(), &tb.ChatMember{
		Rights: tb.NoRestrictions(),
		User:   user.User(),
	})
	if err == nil {
		c.Reply("✨ Successfully demoted! ~")
	} else if err.Error() == "telegram: can't remove chat owner (400)" {
		c.Reply("I don't really feel like staging a mutiny today, I think the chat owner deserves to stay an admin.")
	} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
		c.Reply("Failed to demote, " + "Make sure I'm admin and can appoint new admins.")
	} else if err.Error() == "telegram unknown: Bad Request: CHAT_ADMIN_REQUIRED (400)" {
		c.Reply("This user has been already promoted by someone otherthan me; I can't change their permissions!")
	} else {
		c.Reply("Failed to demote, " + fmt.Sprint(err.Error()))
	}
	return nil
}

func Adminlist(c tb.Context) error {
	admins := fmt.Sprintf("<b>✨ Admins</b> in <b>%s</b>", c.Chat().Title)
	adminsOf, _ := c.Bot().AdminsOf(c.Chat())
	var creator = []string{}
	var admin = [][]string{}
	for _, ad := range adminsOf {
		if ad.Role == "administrator" {
			admin = append(admin, []string{ad.User.FirstName, fmt.Sprint(ad.User.ID)})
		} else {
			creator = append(creator, ad.User.FirstName)
			creator = append(creator, fmt.Sprint(ad.User.ID))
		}
	}
	admins += fmt.Sprintf("\n<b>-</b> <a href='tg://user?id=%s'>%s</a>\n", creator[1], creator[0])
	for _, ad := range admin {
		admins += fmt.Sprintf("\n<b>-</b> <a href='tg://user?id=%s'>%s</a>", ad[1], ad[0])
	}
	admins += fmt.Sprintf("\n\n<b>Admins Count:</b> %s", fmt.Sprint(len(admin)+1))
	return c.Reply(admins)
}

func Set_title(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats.")
		return nil
	}
	user, title := GetUser(c)
	if user.ID == 0 {
		return nil
	} else if title == string("") {
		return c.Reply("You need to give me a title name")

	} else if user.ID == BOT_ID {
		return c.Reply("I cannot edit my own Title!")
	}
	EditTitle(c, user.User(), title, false)
	return nil
}

func EditTitle(c tb.Context, user *tb.User, title string, promote bool) {
	if promote {
		time.Sleep(3 * time.Second)
	}
	err := c.Bot().SetAdminTitle(c.Chat(), user, title)
	if promote {
		return
	}
	if err == nil {

		c.Reply(fmt.Sprintf("<b>%s</b>'s Admin title was changed to <b>%s</b>.", user.FirstName, title))
		return
	} else if err.Error() == "telegram unknown: Bad Request: user is not an administrator (400)" {
		c.Reply("This user is not an admin!")
	} else if err.Error() == "telegram unknown: Bad Request: not enough rights to change custom title of the user (400)" {
		c.Reply("I don't have permission, because that user was promoted by other someone else.")
	} else if err.Error() == "telegram unknown: Bad Request: not enough rights (400)" {
		c.Reply("Failed to change admin title, Make sure I'm admin and can appoint new admins.")
	} else if err.Error() == "telegram unknown: Bad Request: only creator can edit their custom title (400)" {
		c.Reply("I don't have permission to edit admin title of the chat's creator.")
	} else if err.Error() == "telegram unknown: Bad Request: ADMIN_RANK_EMOJI_NOT_ALLOWED (400)" {
		c.Reply("Admin titles cannot contain emoji.")
	} else {
		c.Reply(err.Error())
	}
}
