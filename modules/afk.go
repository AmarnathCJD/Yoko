package modules

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func AFK(c tb.Context) bool {
	if c.Sender() == nil || c.Media() != nil {
		return false
	}
	prefix := strings.SplitN(c.Text(), " ", 2)
	if stringInSlice(strings.ToLower(prefix[0]), []string{"/afk", "brb"}) {
		reason := ""
		if len(prefix) == 2 {
			reason = prefix[1]
		}
		c.Reply(fmt.Sprintf("<b>%s</b> is now AFK !", c.Sender().FirstName))
		db.Set_afk(c.Sender().ID, c.Sender().FirstName, reason)
		return true
	} else if db.IsAfk(c.Sender().ID) {
		rand.Seed(time.Now().Unix())
		c.Reply(fmt.Sprintf(AFK_STR[rand.Intn(len(AFK_STR))], c.Sender().FirstName))
		db.Unset_afk(c.Sender().ID)
		return true
	} else {
		var user_id int64
		if c.Message().IsReply() {
			user_id = int64(c.Message().ReplyTo.Sender.ID)
		} else {
			for _, e := range c.Message().Entities {
				if e.Type == tb.EntityMention || e.Type == tb.EntityTMention {
					if e.User == nil {
						u, _ := getJson(c.Message().EntityText(e))
						user_id = int64(u["id"].(float64))
					} else {
						user_id = e.User.ID
					}
				}
			}
		}
		if user_id != int64(0) {
			if db.IsAfk(user_id) {
				a := db.GetAfk(user_id)
				fmt.Println("hi")
				var reason = ""
				if r, ok := a["reason"]; ok && r != nil {
					reason = fmt.Sprintf("<b>Reason:</b> %s", r.(string))
				}
				c.Reply(fmt.Sprintf("<b>%s</b> is AFK !\nLast Seen: %s ago.\n%s", a["fname"].(string), get_time_value(a["time"].(int32)), reason))
				return true
			}
		}
	}
	return false
}

func Sed_reg(c tb.Context) bool {
	if !c.Message().IsReply() {
		return false
	}
	msg := c.Message().ReplyTo.Text
	if msg == string("") {
		return false
	} else if to_parse := c.Message().Text; to_parse == string("") {
		return false
	}
	to_parse := c.Message().Text
	api_url := "https://polar-refuge-17864.herokuapp.com/sed"
	req, _ := http.NewRequest("GET", api_url, nil)
	q := req.URL.Query()
	q.Add("text", to_parse)
	q.Add("sed", msg)
	req.URL.RawQuery = q.Encode()
	resp, err := myClient.Do(req)
	check(err)
	var body mapType
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&body)
	if text, ok := body["text"]; ok {
		text = text.(string)
		c.Send(text, tb.SendOptions{ReplyTo: c.Message().ReplyTo})
	}
	return true
}
