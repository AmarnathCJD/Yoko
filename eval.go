package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	eval "github.com/apaxa-go/eval"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Shellout(command string) (string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), err.Error()
}

func evaluate(m *tb.Message) {
	if m.Sender.ID != 1833850637 {
		return
	}
	if string(m.Payload) == string("") {
		return
	}
	expr, _ := eval.ParseString(string(m.Payload), "")
	a := eval.Args{"fmt.Sprint": eval.MakeDataRegularInterface(fmt.Sprint), "bot": eval.MakeDataRegularInterface(b), "e": eval.MakeDataRegularInterface(m), "import": eval.MakeDataRegularInterface("import"), "get_user": eval.MakeDataRegularInterface(get_user), "db": eval.MakeDataRegularInterface(db), "ctx": eval.MakeDataRegularInterface(context.TODO())}
	r, err := expr.EvalToInterface(a)
	if r != nil {
		_, err := b.Reply(m, "<b>► EVALGo</b>\n"+string(m.Payload)+"\n\n<b>► OUTPUT</b>\n<code>"+fmt.Sprint(r)+"</code>")
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		return
	}
	b.Reply(m, "<b>► EVALGo</b>\n"+string(m.Payload)+"\n\n<b>► OUTPUT</b>\n<code>"+fmt.Sprint(err.Error())+"</code>")
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
