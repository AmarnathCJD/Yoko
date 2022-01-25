package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

type USER struct {
	ID    int64
	COUNT int32
}

var DB = make(map[int64]map[int64]USER)

func reset_flood(chat_id int64, user_id int64) {
	for _, x := range DB[chat_id] {
		x.COUNT = 0
	}
}

func Flood(c tb.Context) error {
	if c.Message().Private() {
		return c.Reply("This command is made to be used in group chats.")
	} else {
		f := db.GetFlood(c.Chat().ID)
		if f.COUNT == 0 {
			return c.Reply("This chat is not currently enforcing flood control.")
		} else {
			return c.Reply(fmt.Sprintf("This chat is currently enforcing flood control after %d messages. Any users sending more than that amount of messages will be %s.", f.COUNT, Convert_action(f.MODE, f.TIME)))
		}
	}
}

func SetFlood(c tb.Context) error {
	if c.Message().Private() {
		return c.Reply("This command is made to be used in group chats.")
	} else {
		if c.Message().Payload == string("") {
			return c.Reply("I expected some arguments! Either off, or an integer. eg: <code>/setflood 5</code>, or <code>/setflood off?</code>")
		} else if c.Message().Payload == "off" || c.Message().Payload == "0" {
			c.Reply("Antiflood has been disabled.")
			db.SetFloodCount(c.Chat().ID, 0)
			return nil
		} else if isInt(c.Message().Payload) {
			d, _ := strconv.Atoi(c.Message().Payload)
			c.Reply(fmt.Sprintf("Antiflood settings for %s have been updated to %d", c.Chat().Title, d))
			db.SetFloodCount(c.Chat().ID, int32(d))
			return nil
		} else {
			return c.Reply(fmt.Sprintf("%s is not a valid integer.", c.Message().Payload))
		}
	}
}

func SetFloodMode(c tb.Context) error {
	if c.Message().Private() {
		return c.Reply("This command is made to be used in group chats.")
	} else {
		ctime := int64(0)
		p := c.Message().Payload
		if p == string("") {
			return c.Reply("You need to specify an action to take upon flooding. Current modes are: ban/kick/mute/tban/tmute")
		} else {
			args := strings.SplitN(p, " ", 2)
			if stringInSlice(args[0], []string{"ban", "mute", "kick", "tban", "tmute"}) {
				if strings.HasPrefix(args[0], "t") {
					if len(args) < 2 {
						c.Reply("It looks like you tried to set time value for antiflood but you didn't specified time; Try, <code>/setfloodmode [tban/tmute] <timevalue></code>.\n<b>Examples of time value:</b> <code>4m = 4 minutes</code>, <code>3h = 3 hours</code>, <code>6d = 6 days</code>, <code>5w = 5 weeks</code>.")
						return nil
					}
					ctime = Extract_time(c, args[1])
					if ctime == 0 {
						return nil
					}
				}
			} else {
				c.Reply(fmt.Sprintf("Unknown type '%s'. Please use one of: ban/kick/mute/tban/tmute", args[0]))
				return nil
			}
			c.Reply(fmt.Sprintf("Updated antiflood reaction in %s to: %s", c.Chat().Title, Convert_action(args[0], int32(ctime))))
			db.SetFloodMode(c.Chat().ID, args[0], int32(ctime))
			return nil
		}

	}

}

func FLOOD_EV(c tb.Context) bool {
	chat_id := c.Chat().ID
	f := db.GetFlood(chat_id)
	if f.COUNT == 0 {
		return false
	} else {
		if c.Sender() == nil {
			reset_flood(chat_id, 0)
			return false
		}
		if _, ok := DB[chat_id][c.Sender().ID]; !ok {
			DB[chat_id] = make(map[int64]USER)
			DB[chat_id][c.Sender().ID] = USER{ID: c.Sender().ID, COUNT: 1}
			fmt.Println(DB)
		} else {
			u := DB[chat_id][c.Sender().ID]
			DB[chat_id][c.Sender().ID] = USER{ID: c.Sender().ID, COUNT: u.COUNT + 1}
			fmt.Println(DB)
		}
		if DB[chat_id][c.Sender().ID].COUNT > f.COUNT {
			c.Reply("Flood")
			reset_flood(chat_id, 0)
			return true
		}

	}
	return false
}
