package modules

import (
	"fmt"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func Approve(c tb.Context) error {
	User, R := GetUser(c)
	if User.ID == 0 {
		return nil
	}
	p, _ := c.Bot().ChatMemberOf(c.Chat(), User.User())
	if p.Role == tb.Administrator || p.Role == tb.Creator {
		return c.Reply("Why would i approve a admin!")
	} else if c.Sender().ID == User.ID {
		return c.Reply("Why are you trying to approve yourself ?")
	} else if User.ID == BOT_ID {
		return c.Reply("Why are you trying to approve me ?")
	} else if User.Approved(c.Chat().ID) {
		return c.Reply("This user is already approved")
	}
	db.Approve(c.Chat().ID, User.ID)
	return c.Reply(fmt.Sprintf("Successfully Approved User%s", GetReason(R)))
}

func Unapprove(c tb.Context) error {
	User, R := GetUser(c)
	if User.ID == 0 {
		return nil
	}
	p, _ := c.Bot().ChatMemberOf(c.Chat(), User.User())
	if p.Role == tb.Administrator || p.Role == tb.Creator {
		return c.Reply("Why would i unapprove a admin!")
	} else if c.Sender().ID == User.ID {
		return c.Reply("Why are you trying to unapprove yourself ?")
	} else if User.ID == BOT_ID {
		return c.Reply("Why are you trying to unapprove me ?")
	} else if !User.Approved(c.Chat().ID) {
		return c.Reply("This user is not Approved")
	}
	db.Unapprove(c.Chat().ID, User.ID)
	return c.Reply(fmt.Sprintf("Successfully Unapproved User%s", GetReason(R)))
}

func Approved(c tb.Context) error {
	All := db.GetAllApproved(c.Chat().ID)
	if All == nil {
		return c.Reply("No one is approved in this chat")
	}
	return c.Reply(fmt.Sprintf("%v", All))
}

func Approval(c tb.Context) error {
	User, _ := GetUser(c)
	if User.ID == 0 {
		return nil
	}
	p, _ := c.Bot().ChatMemberOf(c.Chat(), User.User())
	if p.Role == tb.Administrator || p.Role == tb.Creator {
		return c.Reply("Why will check status of an admin ?")
	} else if User.ID == BOT_ID {
		return c.Reply("I am not gonna check my status")
	}
	if User.Approved(c.Chat().ID) {
		return c.Reply("This User is Approved")
	} else {
		return c.Reply("This User is not Approved")
	}
}

func DisapproveAll(c tb.Context) error {
	sendErr := c.Reply("Successfully disapproved everyone in the chat.")
	db.UnapproveAll(c.Chat().ID)
	return sendErr
}
