package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	tb "gopkg.in/tucnak/telebot.v3"
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

func Eval(c tb.Context) error {
	return nil
	cmd := strings.SplitN(c.Message().Text, " ", 2)
	if len(cmd) == 1 {
		return nil
	}
	api_url := "https://api.roseloverx.in/go"
	req, _ := http.NewRequest("GET", api_url, nil)
	q := req.URL.Query()
	q.Add("code", "package main\n"+cmd[1])
	req.URL.RawQuery = q.Encode()
	resp, err := myClient.Do(req)
	if err != nil {
		c.Reply(err.Error())
	}
	defer resp.Body.Close()
	var body mapType
	json.NewDecoder(resp.Body).Decode(&body)
	if body["errors"].(string) != string("") {
		c.Reply(fmt.Sprintf("<b>► EvalGo</b>\n<code>%s</code>\n\n<b>►</b> Output\n<code>%s</code>", cmd[1], body["errors"].(string)))
	} else if body["output"] != string("") {
		c.Reply(fmt.Sprintf("<b>► EvalGo</b>\n<code>%s</code>\n\n<b></b>► Output\n<code>%s</code>", cmd[1], body["output"].(string)))
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
		return c.Reply(string(b), &tb.SendOptions{ParseMode: "Markdown"})
	}

}
