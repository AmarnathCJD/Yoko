package modules

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
        "fmt"
        "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func Chat_bot(c tb.Context) error {
	is_chat := false
        if c.Media() != nil {
           return nil
        }
	if c.Message().IsReply() && c.Message().ReplyTo.Sender.ID == int64(5050904599) {
		is_chat = true
	} else if strings.Contains(strings.ToLower(c.Message().Text), "yoko") {
		is_chat = true
	}
	if !is_chat {
		return nil
	} else if strings.HasPrefix(c.Message().Text, "/") {
                return nil
        }
	text := strings.ReplaceAll(c.Message().Text, "yoko", "kuki")
	url_q := "https://icap.iconiq.ai/talk?&botkey=icH-VVd4uNBhjUid30-xM9QhnvAaVS3wVKA3L8w2mmspQ-hoUB3ZK153sEG3MX-Z8bKchASVLAo~&channel=7&sessionid=482070240&client_name=uuiprod-un18e6d73c-user-19422&id=true"
	req, err := http.PostForm(url_q, url.Values{"input": {text}})
	if err != nil {
		c.Reply(err.Error())
	}
	defer req.Body.Close()
	var resp mapType
	json.NewDecoder(req.Body).Decode(&resp)
	msg := resp["responses"].([]interface{})[0].(string)
        fmt.Println(msg)
	pattern := regexp.MustCompile(`<image>.+</image>`)
	media := pattern.FindAllStringSubmatch(msg, -1)
	if media != nil {
		if len(media) != 0 {
			file := strings.ReplaceAll(strings.ReplaceAll(media[0][0], "<image>", ""), "</image>", "")
                        if strings.Contains(file, "pandorabots"){
                             file = strings.ReplaceAll(strings.ReplaceAll(media[1][0], "<image>", ""), "</image>", "")
                        }
                        if strings.HasSuffix(file, "jpg") || strings.HasSuffix(file, "png") {
                             c.Reply(&tb.Photo{File: tb.FromURL(file)})
                        } else {
			     c.Reply(&tb.Animation{File: tb.FromURL(file)})
                        }
		}
	}
	chat := strings.SplitN(msg, "</image>", 2)
	var message string
	if len(chat) == 2 {
		message = chat[1]
	} else {
		message = chat[0]
	}
        message = strings.ReplaceAll(message, "kuki", "yoko")
	c.Reply(message)
	return nil
}

func Chatbot_mode(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made for group chats only.")
		return nil
	}
	args := c.Message().Payload
	if args == string("") {
		mode := db.Get_chatbot_mode(c.Chat().ID)
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
		db.Set_chatbot_mode(c.Chat().ID, true)
	} else if stringInSlice(strings.ToLower(args), []string{"disable", "off", "no", "n"}) {
		c.Reply("<b>Disabled</b> AI chatbot for this chat.")
		db.Set_chatbot_mode(c.Chat().ID, false)
	} else {
		c.Reply("Your input was not recognised as one of: yes/no/y/n/on/off")
	}
	return nil
}
