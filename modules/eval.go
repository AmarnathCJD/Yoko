package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos72/gomacro/fast"
	tb "gopkg.in/telebot.v3"
	"io/ioutil"
	"os"
	"os/exec"
)

func Exec(c tb.Context) error {
	if c.Sender().ID != int64(1833850637) {
		return nil
	}
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
			c.Reply(fmt.Sprintf("<code>Mika#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, stdout.String()))
		} else if stderr.String() != string("") {
			c.Reply(fmt.Sprintf("<code>Mika#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, stderr.String()))
		} else if err != nil {
			c.Reply(fmt.Sprintf("<code>Mika#~</code>: <code>%s</code>\n<code>%s</code>", c.Message().Payload, err.Error()))
		}
	}
	return nil
}

func MediaInfo(c tb.Context) error {
	if c.Sender().ID != int64(1833850637) {
		return nil
	}
	if !c.Message().IsReply() {
		return c.Reply("Reply to a media!")
	} else {
		var b []byte
		Reply := c.Message().ReplyTo
		if Reply.Audio != nil {
			b, _ = json.Marshal(Reply.Audio)
		} else if Reply.Document != nil {
			b, _ = json.Marshal(Reply.Document)
		} else if Reply.Animation != nil {
			b, _ = json.Marshal(Reply.Animation)
		} else if Reply.Video != nil {
			b, _ = json.Marshal(Reply.Video)
		} else if Reply.Photo != nil {
			b, _ = json.Marshal(Reply.Photo)
		} else if Reply.Sticker != nil {
			b, _ = json.Marshal(Reply.Sticker)
		} else if Reply.Voice != nil {
			b, _ = json.Marshal(Reply.Voice)
		} else if Reply.Contact != nil {
			b, _ = json.Marshal(Reply.Voice)
		} else if Reply.Location != nil {
			b, _ = json.Marshal(Reply.Voice)
		} else if Reply.VideoNote != nil {
			b, _ = json.Marshal(Reply.VideoNote)
		} else {
			return c.Reply("Reply to a media!")
		}
		return c.Reply("<code>" + string(b) + "</code>")
	}

}

func Eval(c tb.Context) error {
	code := c.Message().Payload
	out := EvalCmd(code)
	return c.Reply(fmt.Sprintf("Eval: %s\nOut: %s", code, out))
}

func EvalCmd(code string) string {
	interp := fast.New()
	rd := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	interp.Eval(code)
        w.Close()
	out, _ := ioutil.ReadAll(r)
	fmt.Println(out)
	os.Stdout = rd
	return string(out)
}
