package modules

import tb "gopkg.in/tucnak/telebot.v3"
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
        if p.Role == tb.Member || p.Role == tb.Left {
return c.Respond(&tb.CallbackResponse{Text: "You should be an admin to do this", ShowAlert: true})

}
	return nil
}

// Soon, today
