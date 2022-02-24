package modules

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

var (
	CHAT_API = "https://icap.iconiq.ai/talk?&botkey=icH-VVd4uNBhjUid30-xM9QhnvAaVS3wVKA3L8w2mmspQ-hoUB3ZK153sEG3MX-Z8bKchASVLAo~&channel=7&sessionid=482070240&client_name=uuiprod-un18e6d73c-user-19422&id=true"
	NameRe   = regexp.MustCompile(`(?i)mika`)
	MediaRe  = regexp.MustCompile(`<image>(.*?)</image>`)
	BtnRe    = regexp.MustCompile(`<button>(.*?)</button>`)
	UrlRe    = regexp.MustCompile(`<url>(.*?)</url>`)
	TextRe   = regexp.MustCompile(`<text>(.*?)</text>`)
	PostBk   = regexp.MustCompile(`<postback>(.*?)</postback>`)
	CardRe   = regexp.MustCompile(`<card>.+</card>`)
	SplitRe  = regexp.MustCompile(`<split>(.*?)</split>`)
	ChatBTN  = sel.Data("Mika", "chatbot_cb")
)

func Chat_bot(c tb.Context) error {
	if !db.IsChatbot(c.Chat().ID) {
		return nil
	}
	is_chat := false
	replace_addition := false
	if c.Message().IsReply() && c.Message().ReplyTo.Sender.ID == BOT_ID {
		is_chat = true
	} else if strings.Contains(strings.ToLower(c.Message().Text), "mika") {
		replace_addition = true
		is_chat = true
	}
	if c.Message().Private() {
		is_chat = true
	}
	if !is_chat {
		return nil
	} else if strings.HasPrefix(c.Message().Text, "/") || strings.HasPrefix(c.Message().Text, "!") || strings.HasPrefix(c.Message().Text, "?") {
		return nil
	}
	text := strings.ReplaceAll(c.Message().Text, "mika", "kuki")
	if replace_addition {
		text = NameRe.ReplaceAllString(text, "")
	}
	req, err := http.PostForm(CHAT_API, url.Values{"input": {text}})
	if err != nil {
		c.Reply("Chatbot API is down!")
	}
	defer req.Body.Close()
	var resp mapType
	json.NewDecoder(req.Body).Decode(&resp)
	msg := resp["responses"].([]interface{})[0].(string)
	return ExtractMeta(msg, c, false)
}

type BTN struct {
	Text string `json:"text,omitempty"`
	URL  string `json:"url,omitempty"`
}

func ExtractMeta(t string, c tb.Context, Edit bool) error {
	var Med []string
	var Btn []BTN
	var Inline bool
	var Card = false
	for _, x := range MediaRe.FindAllStringSubmatch(t, -1) {
		if strings.Contains(x[1], "giphylogo.png") {
			continue
		}
		Med = append(Med, x[1])
		t = strings.Replace(t, x[0], "", 1)
	}
	t = strings.ReplaceAll(strings.ReplaceAll(t, "<delay>", ""), "</delay>", "")
	t = strings.ReplaceAll(strings.ReplaceAll(t, "Kuki", "Mika"), "kuki", "mika")
	if match := CardRe.FindAllStringSubmatch(t, -1); match != nil {
		Card = true
		t = strings.Replace(t, match[0][0], "", 1)
	}
	t = strings.ReplaceAll(strings.ReplaceAll(t, "<reply>", "<button>"), "</reply>", "</button>")
	for _, y := range BtnRe.FindAllStringSubmatch(t, -1) {
		url, text := "", "Link"
		t = strings.Replace(t, string(y[0]), "", -1)
		for _, x := range UrlRe.FindAllStringSubmatch(string(y[1]), -1) {
			url = x[1]
			Inline = false
		}
		for _, x := range TextRe.FindAllStringSubmatch(string(y[1]), -1) {
			text = x[1]
		}
		for _, x := range PostBk.FindAllStringSubmatch(string(y[1]), -1) {
			url = x[1]
			Inline = true
		}
		text = strings.ReplaceAll(text, "Kuki", "Mika")
		Btn = append(Btn, BTN{text, url})
	}
	if Inline {
		var rows []tb.Row
		for _, x := range Btn {
			rows = append(rows, sel.Row(menu.Data(x.Text, "chatbot_cb", x.URL)))
		}
		sel.Inline(rows...)
	} else {
		var rows []tb.Row
		for _, x := range Btn {
			rows = append(rows, sel.Row(sel.URL(x.Text, x.URL)))
		}
		sel.Inline(rows...)
	}
	if Card {
		if len(Med) > 0 {
			media := Med[0]
			if strings.Contains(media, "png") || strings.Contains(media, "jpg") || strings.Contains(media, "jpeg") {
				if Edit {
					c.Edit(&tb.Photo{File: tb.FromURL(media), Caption: t}, sel)
				} else {
					c.Reply(&tb.Photo{File: tb.FromURL(media), Caption: t}, sel)
				}
			} else if strings.Contains(media, "gif") {
				if Edit {
					c.Edit(&tb.Animation{File: tb.FromURL(media), Caption: t}, sel)
				} else {
					c.Reply(&tb.Animation{File: tb.FromURL(media), Caption: t}, &tb.SendOptions{ReplyMarkup: sel})
				}
			} else {
				if Edit {
					c.Edit(&tb.Document{File: tb.FromURL(media), Caption: t}, sel)
				} else {
					c.Reply(&tb.Document{File: tb.FromURL(media), Caption: t}, &tb.SendOptions{ReplyMarkup: sel})
				}
			}
		}
	} else {
		for _, x := range Med {
			if strings.Contains(x, "png") || strings.Contains(x, "jpg") || strings.Contains(x, "jpeg") {
				if Edit {
					c.Edit(&tb.Photo{File: tb.FromURL(x)}, sel)
				} else {
					c.Reply(&tb.Photo{File: tb.FromURL(x)})
				}
			} else if strings.Contains(x, "gif") {
				if Edit {
					c.Edit(&tb.Animation{File: tb.FromURL(x)}, sel)
				} else {
					c.Reply(&tb.Animation{File: tb.FromURL(x)})
				}
			} else {
				if Edit {
					c.Edit(&tb.Document{File: tb.FromURL(x)}, sel)
				} else {
					c.Reply(&tb.Document{File: tb.FromURL(x)})
				}
			}
		}
		t = strings.TrimSpace(t)
		var Nmsg string
		if !Edit {
			for _, x := range SplitRe.FindAllStringSubmatch(t, -1) {
				Nmsg = x[1]
				t = strings.Replace(t, x[0], "", 1)
			}
		} else {
			t = strings.ReplaceAll(strings.ReplaceAll(t, "<split>", "\n"), "</split>", "\n")
		}
		if t != string("") && !Card {
			t = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(t, "<", ""), ">", ""), "/", "")
			if Edit {
				return c.Edit(t, sel)
			} else {
				c.Reply(t, &tb.SendOptions{ReplyMarkup: sel})
				if Nmsg != string("") {
					return c.Send(Nmsg)
				}
			}
		}

	}
	return nil
}

func ChatbotCB(c tb.Context) error {
	d := c.Callback().Data
	req, err := http.PostForm(CHAT_API, url.Values{"input": {d}})
	if err != nil {
		return c.Respond(&tb.CallbackResponse{Text: "try again!", ShowAlert: true})
	}
	defer req.Body.Close()
	var resp mapType
	json.NewDecoder(req.Body).Decode(&resp)
	msg := resp["responses"].([]interface{})[0].(string)
	return ExtractMeta(msg, c, true)
}

func Chatbot_mode(c tb.Context) error {
	args := c.Message().Payload
	if args == string("") {
		mode := db.IsChatbot(c.Chat().ID)
		if mode {
			c.Reply("AI chatbot is currently <b>enabled</b> for this chat.")
			return nil
		} else {
			c.Reply("AI chatbot is currently <b>disabled</b> for this chat.")
			return nil
		}
	}
	if stringInSlice(strings.ToLower(args), []string{"enable", "on", "yes", "y"}) {
		c.Reply("<b>Enabled</b> AI chatbot for this chat.")
		db.SetCHatBotMode(c.Chat().ID, true)
	} else if stringInSlice(strings.ToLower(args), []string{"disable", "off", "no", "n"}) {
		c.Reply("<b>Disabled</b> AI chatbot for this chat.")
		db.SetCHatBotMode(c.Chat().ID, false)
	} else {
		c.Reply("Your input was not recognised as one of: yes/no/y/n/on/off")
	}
	return nil
}
