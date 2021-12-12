package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	eval "github.com/apaxa-go/eval"
	tb "gopkg.in/tucnak/telebot.v3"
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

func evaluate(c tb.Context) error {
	m := c.Message()
	if m.Sender.ID != 1833850637 {
		return nil
	}
	if string(c.Message().Payload) == string("") {
		return nil
	}
	expr, _ := eval.ParseString(string(c.Message().Payload), "")
	a := eval.Args{"fmt.Sprint": eval.MakeDataRegularInterface(fmt.Sprint), "bot": eval.MakeDataRegularInterface(b), "e": eval.MakeDataRegularInterface(c.Message()), "import": eval.MakeDataRegularInterface("import"), "get_user": eval.MakeDataRegularInterface(get_user), "db": eval.MakeDataRegularInterface(db), "ctx": eval.MakeDataRegularInterface(context.TODO())}
	r, err := expr.EvalToInterface(a)
	if r != nil {
		_, err := b.Reply(m, "<b>► EVALGo</b>\n"+string(c.Message().Payload)+"\n\n<b>► OUTPUT</b>\n<code>"+fmt.Sprint(r)+"</code>")
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		return nil
	}
	b.Reply(m, "<b>► EVALGo</b>\n"+string(m.Payload)+"\n\n<b>► OUTPUT</b>\n<code>"+fmt.Sprint(err.Error())+"</code>")
	return nil
}

func execute(c tb.Context) error {
	m := c.Message()
	if m.Sender.ID != 1833850637 {
		return nil
	}
	if string(m.Payload) == string("") {
		b.Reply(m, "No CMD given.")
		return nil
	}
	out, err := Shellout(m.Payload)
	output := "<code>Go#~:" + string(out) + string(err) + "</code>"
	b.Reply(m, output)
	return nil
}
