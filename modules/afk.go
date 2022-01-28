package modules

import (
	"fmt"
	"math/rand"
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
				c.Reply(fmt.Sprint(a))
				reason := "."
				if r, ok := a["reason"]; ok {
					reason = fmt.Sprintf(", <b>Reason:</b> %s", r.(string))
				}
				err := c.Reply(fmt.Sprintf("<b>%s</b> is AFK !\nLast Seen: %s ago.\n%s", a["fname"].(string), get_time_value(a["time"].(int32)), reason))
				check(err)
				return true
			}
		}
	}
	return false
}
