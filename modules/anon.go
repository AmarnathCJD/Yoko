package modules

import tb "gopkg.in/tucnak/telebot.v3"
import "encoding/json"
import "fmt"

var (
	ANON        = make(map[int]Update)
	anon_button = sel.Data("Click to prove admin", "anon_btn")
)

type Update struct {
	Id    int
	Func  func(tb.Context) error
	Right string
	C     tb.Context
}

func AnonAdmin(next tb.HandlerFunc, p string, c tb.Context) {
	sel.Inline(sel.Row(anon_button))
	msg, _ := c.Bot().Send(c.Chat(), "It looks like you're anonymous. Tap this button to confirm your identity.", &tb.SendOptions{ReplyMarkup: sel, ReplyTo: c.Message()})
	ANON[msg.ID] = Update{c.Message().ID, next, p, c}
	fmt.Println(ANON)
}

func AnonCB(c tb.Context) error {
	update, exist := ANON[c.Callback().Message.ID]
	if !exist {
		fmt.Println("Not exist")
	}
	p, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	check(err)
	if p.Role == tb.Member || p.Role == tb.Left {
		return c.Respond(&tb.CallbackResponse{Text: "You should be an admin to do this", ShowAlert: true})

	} else if p.Role == tb.Administrator {
		if update.Right == "ban_users" && !p.Rights.CanRestrictMembers {
			return c.Edit("You are missing the following rights to use this command: CanRestrictMembers")
		} else if update.Right == "change_info" && !p.Rights.CanChangeInfo {
			return c.Edit("You are missing the following rights to use this command: CanChangeInfo")
		} else if update.Right == "pin_messages" && !p.Rights.CanPinMessages {
			return c.Edit("You are missing the following rights to use this command: CanPinMessages")
		} else if update.Right == "add_admins" && !p.Rights.CanPromoteMembers {
			return c.Edit("You are missing the following rights to use this command: CanPromoteMembers")
		}
	}
	b, err := json.Marshal(update.C.Sender())
	s := tb.Update{ID: update.C.Message().ID,
		Message: &tb.Message{
			Sender:  c.Message().Sender,
			Chat:    update.C.Message().Chat,
			Payload: update.C.Message().Payload,
			Text:    update.C.Message().Text,
		},
	}
	fmt.Println(string(b), s)
	return nil
}

// Soon, today
