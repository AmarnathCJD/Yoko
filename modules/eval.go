package modules

import (
	"fmt"
	"os/exec"

	tb "gopkg.in/tucnak/telebot.v3"
)

func Exec(c tb.Context) error {
	if c.Message().Payload == string("") {
		return nil
	} else {
		proc := exec.Command("bash", "-c", c.Message().Payload)
		output, err := proc.Output()
		if string(output) != string("") {
			c.Reply(fmt.Sprintf("<code>Yoko#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, string(output)))
		} else if err.Error() != string("") {
			c.Reply(fmt.Sprintf("<code>Yoko#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, err.Error()))
		} else {
			c.Reply(fmt.Sprintf("<code>Yoko#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, "No Output/Error"))
		}
	}
	return nil
}
