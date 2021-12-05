package main

import (
	"fmt"
	"strings"

	eval "github.com/apaxa-go/eval"
	tb "gopkg.in/tucnak/telebot.v2"
)

func evaluate(m *tb.Message) {
	if m.Sender.ID != 1833850637 {
		return
	}
	if string(m.Payload) == string("") {
		return
	}
	expr, _ := eval.ParseString(string(m.Payload), "")
	a := eval.Args{"fmt.Sprint": eval.MakeDataRegularInterface(fmt.Sprint), "bot": eval.MakeDataRegularInterface(b), "e": eval.MakeDataRegularInterface(m)}
	r, err := expr.EvalToInterface(a)
	b.Reply(m, "<b>► EVALGo</b>\n"+string(m.Payload)+"\n\n<b>► OUTPUT</b>\n<code>"+fmt.Sprint(err.Error())+fmt.Sprint(r)+"</code>")
}

func execute(m *tb.Message) {
	if m.Sender.ID != 1833850637 {
		return
	}
	if string(m.Payload) == string("") {
		b.Reply(m, "No CMD given.")
		return
	}
	out, err := Shellout(m.Payload)
	output := "<code>Go#~:" + string(out) + string(err) + "</code>"
	b.Reply(m, output)
}
