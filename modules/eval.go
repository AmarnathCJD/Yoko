package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	tb "gopkg.in/telebot.v3"
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
			b, _ = json.MarshalIndent(Reply.Audio, "", "    ")
		} else if Reply.Document != nil {
			b, _ = json.MarshalIndent(Reply.Document, "", "    ")
		} else if Reply.Animation != nil {
			b, _ = json.MarshalIndent(Reply.Animation, "", "    ")
		} else if Reply.Video != nil {
			b, _ = json.MarshalIndent(Reply.Video, "", "    ")
		} else if Reply.Photo != nil {
			b, _ = json.MarshalIndent(Reply.Photo, "", "    ")
		} else if Reply.Sticker != nil {
			b, _ = json.MarshalIndent(Reply.Sticker, "", "    ")
		} else if Reply.Voice != nil {
			b, _ = json.MarshalIndent(Reply.Voice, "", "    ")
		} else if Reply.Contact != nil {
			b, _ = json.MarshalIndent(Reply.Voice, "", "    ")
		} else if Reply.Location != nil {
			b, _ = json.MarshalIndent(Reply.Voice, "", "    ")
		} else if Reply.VideoNote != nil {
			b, _ = json.MarshalIndent(Reply.VideoNote, "", "    ")
		} else {
			b, _ = json.MarshalIndent(Reply, "", "    ")
		}
		return c.Reply("<code>" + string(b) + "</code>")
	}

}

type FmtResp struct {
	Body  string `json:"body"`
	Error string `json:"error"`
}

type CompileResp struct {
	Errors string `json:"Errors"`
	Events []struct {
		Message string `json:"Message"`
		Kind    string `json:"Kind"`
		Delay   int    `json:"Delay"`
	} `json:"Events"`
	Err string `json:"Err"`
	Out string `json:"Out"`
}

func Eval(c tb.Context) error {
	Addon := `package main\n`
	Code := GetArgs(c)
	if c.Sender().ID != OWNER_ID {
		return nil
	}
	if Code == "" {
		return nil
	}
	if !strings.Contains(Code, "package main") {
		Code = Addon + Code
	}
	FmtApi := "https://go.dev/_/fmt?backend="
	GoCompile := "https://go.dev/_/compile?backend="
	req, _ := http.NewRequest("POST", FmtApi, strings.NewReader(`body=`+Code+`&imports=true`))
	req.Header.Set("Content-Type", "application/json")
	formatted, _ := Client.Do(req)
	if formatted.StatusCode != 200 {
		return c.Reply("Error formatting code!")
	}
	var resp FmtResp
	json.NewDecoder(formatted.Body).Decode(&resp)
	if resp.Error != "" {
		return c.Reply(resp.Error)
	}
	goreq, _ := http.NewRequest("POST", GoCompile, strings.NewReader(`version=2&body=`+resp.Body+`&withVet=true`))
	goreq.Header.Set("Content-Type", "application/json")
	compiled, _ := Client.Do(goreq)
	if compiled.StatusCode != 200 {
		return c.Reply("Error compiling code!")
	}
	var resp2 CompileResp
	json.NewDecoder(compiled.Body).Decode(&resp2)
	var Evaluation string
	if resp2.Errors != "" {
		Evaluation = resp2.Errors
	} else if resp2.Events != nil {
		for _, event := range resp2.Events {
			Evaluation += event.Message + "\n"
		}
	} else if resp2.Out != "" {
		Evaluation = resp2.Out
	} else if resp2.Err != "" {
		Evaluation = resp2.Err
	} else {
		Evaluation = "No output"
	}
	FinalOutput := fmt.Sprintf(
		"__►__ <b>EVALGo</b>\n```%s``` \n\n __►__ <b>OUTPUT</b>: \n```%s``` \n",
		resp.Body,
		Evaluation,
	)
	return c.Reply(FinalOutput)
}
