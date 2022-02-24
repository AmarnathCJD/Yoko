package modules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

func SetRules(c tb.Context) error {
	Rules := GetArgs(c)
	if Rules == "" {
		return c.Reply("You haven't specified any rules!")
	} else {
		db.SetRules(c.Chat().ID, Rules)
		return c.Reply("Rules have been set!")
	}
}

func Rules(c tb.Context) error {
	Rules, btns := db.GetRules(c.Chat().ID)
	if Rules == "" {
		return c.Reply("The group admins haven't set any rules for this chat yet.\nThis probably doesn't mean it's lawless though...!")
	} else {
		Private := db.PrivateRules(c.Chat().ID)
		if Private == true {
			sel.Inline(sel.Row(sel.URL(btns, fmt.Sprintf("t.me/%s?start=rules_%d", BOT_USERNAME, c.Chat().ID))))
			return c.Reply("Tap here to view the rules in your private chat.", sel)
		} else {
			return c.Reply(Rules)
		}
	}
}

func ResetRules(c tb.Context) error {
	db.ResetRules(c.Chat().ID)
	return c.Reply("Rules have been reset!")
}

func SetRulesButton(c tb.Context) error {
	Button := GetArgs(c)
	if Button == "" {
		return c.Reply("You haven't specified any button name!")
	} else {
		db.SetRulesButton(c.Chat().ID, Button)
		return c.Reply("Custom rules button name has been set!")
	}
}

func ResetRulesButton(c tb.Context) error {
	db.ResetRulesButton(c.Chat().ID)
	return c.Reply("Rules button name has been reset!")
}

func PrivateRules(c tb.Context) error {
	Private := GetArgs(c)
	if Private == "" {
		P := db.PrivateRules(c.Chat().ID)
		if P == true {
			return c.Reply("Private rules are currently enabled!")
		} else {
			return c.Reply("Private rules are currently disabled!")
		}
	} else {
		if Private == "true" {
			db.SetPrivateRules(c.Chat().ID, true)
			return c.Reply("Private rules has been enabled!")
		} else if Private == "false" {
			db.SetPrivateRules(c.Chat().ID, false)
			return c.Reply("Private rules has been disabled!")
		} else {
			return c.Reply("Invalid input!")
		}
	}
}

func SendPrivateRules(c tb.Context) error {
	args := strings.SplitN(c.Message().Payload, "_", 2)
	ChatID, _ := strconv.Atoi(args[1])
	Rules, _ := db.GetRules(int64(ChatID))
	if Rules == "" {
		return c.Reply("The group admins haven't set any rules for this chat yet.\nThis probably doesn't mean it's lawless though...!")
	} else {
		return c.Reply(Rules)
	}
}
