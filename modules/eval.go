package modules

import (
	"bytes"
	"fmt"
	"os/exec"

	tb "gopkg.in/tucnak/telebot.v3"
)

func Exec(c tb.Context) error {
	if c.Message().Payload == string("") {
		return nil
	} else {
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		proc := exec.Command("bash", "-c", c.Message().Payload)
		proc.Stdout = &stdout
		proc.Stderr = &stderr
		err := proc.Run()
		if stdout.String() != string("") {
			c.Reply(fmt.Sprintf("<code>Yoko#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, stdout.String()))
		} else if stderr.String() != string("") {
			c.Reply(fmt.Sprintf("<code>Yoko#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, stderr.String()))
		} else if err != nil {
			c.Reply(fmt.Sprintf("<code>Yoko#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, err.Error()))
		}
	}
	return nil
}
