package modules

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func AFK(c tb.Context) bool {
	if c.Sender() == nil && c.Message().SenderChat == nil {
		return false
	}
	prefix := strings.SplitN(c.Text(), " ", 2)
	if c.Text() != string("") && stringInSlice(strings.ToLower(prefix[0]), []string{"/afk", "brb"}) {
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
			if c.Text() == string("") {
				return false
			}
			for _, e := range c.Message().Entities {
				if e.Type == tb.EntityMention || e.Type == tb.EntityTMention {
					if e.User == nil {
						u := ResolveUsername(c.Message().EntityText(e))
						if u.ID != 0 {
							user_id = u.ID
						} else {
							return false
						}
					} else {
						user_id = e.User.ID
					}
				}
			}
		}
		if user_id != int64(0) {
			if db.IsAfk(user_id) {
				a := db.GetAfk(user_id)
				reason := "."
				if r, ok := a["reason"]; ok && r.(string) != "" {
					reason = GetReason(r.(string))
				}
				since := get_readable_time(time.Unix(int64(a["time"].(int64)), 0), time.Now())
				if err := c.Reply(fmt.Sprintf("<b>%s</b> is AFK !\nLast Seen: %s ago.%s", FormatString(a["fname"].(string)), since.Truncate(time.Second).String(), FormatString(reason))); err != nil {
					log.Println(err)
				}
				return true
			}
		}
	}
	return false
}
