package modules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

var (
	del_all_filters        = sel.Data("Delete all filters", "stopall")
	cancel_del_all_filters = sel.Data("Cancel", "cancel_del_all")
)

func SaveFilter(c tb.Context) error {
	Message := ParseMessage(c)
	if Message.Name == string("") {
		c.Reply("You need to give the filter a name!")
		return nil
	} else if Message.Text == string("") && Message.File.FileID == string("") {
		c.Reply("You need to give the filter a name!")
		return nil
	}
	db.SaveFilter(c.Chat().ID, Message)
	return c.Reply(fmt.Sprintf("Saved filter '%s'.", Message.Name))
}

func AllFilters(c tb.Context) error {
	f := db.GetFilters(c.Chat().ID)
	if len(f) == 0 {
		return c.Reply(fmt.Sprintf("No filters in %s!", c.Chat().Title))
	} else {
		txt := fmt.Sprintf("<b>Filters in %s:</b>", c.Chat().Title)
		for _, x := range f {
			txt += fmt.Sprintf("\n- <code>%s</code>", x.Name)
		}
		return c.Reply(txt)
	}
}

func StopFilter(c tb.Context) error {
	if c.Message().Payload == string("") {
		return c.Reply("not enough arguments provided.")
	} else {
		if !db.IsFilterExists(c.Chat().ID, c.Message().Payload) {
			return c.Reply("You haven't saved any filters on this word yet!")
		}
		err := db.RemoveFilter(c.Chat().ID, c.Message().Payload)
		if err != nil {
			return c.Reply(fmt.Sprintf("Filter <code>'%s'</code> has been stopped!", c.Message().Payload))
		} else {
			return c.Reply("Error: " + err.Error())
		}
	}
}

func StopAllFIlters(c tb.Context) error {
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == tb.Member || p.Role == tb.Left {
		return c.Reply("You need to be an admin to do this!")
	} else if p.Role == tb.Administrator {
		return c.Reply("You should be the chat creator to do this!")
	} else {
		sel.Inline(sel.Row(del_all_filters), sel.Row(cancel_del_all_filters))
		return c.Reply(fmt.Sprintf("Are you sure you would like to stop <b>ALL</b> filters in %s? This action cannot be undone.", c.Chat().Title), sel)
	}
}

func DelAllFCB(c tb.Context) error {
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == tb.Member || p.Role == tb.Left {
		return c.Respond(&tb.CallbackResponse{Text: "You need to be an admin to do this!", ShowAlert: true})
	} else if p.Role == tb.Administrator {
		return c.Edit("You should be the chat creator to do this!")
	} else if p.Role == tb.Creator {
		c.Edit("Deleted all chat filters.")
		db.PurgeFilters(c.Chat().ID)
	}
	return nil
}

func CancelDALL(c tb.Context) error {
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == tb.Member || p.Role == tb.Left {
		return c.Respond(&tb.CallbackResponse{Text: "You need to be an admin to do this!", ShowAlert: true})
	} else if p.Role == tb.Administrator {
		return c.Edit("You should be the chat creator to do this!")
	} else if p.Role == tb.Creator {
		c.Edit("Stopping of all filters has been cancelled.")
	}
	return nil
}

func FilterEvent(c tb.Context) (bool, error) {
	f := db.GetFilters(c.Chat().ID)
	if len(f) == 0 {
		return false, nil
	}
	for _, x := range f {
		pattern := `( |^|[^\w])(?i)` + x.Name + `( |$|[^\w])`
		if match, _ := regexp.Match(pattern, []byte(c.Text())); match {
			text, p := ParseString(x.Text, c)
			text, btns := button_parser(text)
			if x.File.FileID != string("") {
				f := GetSendable(x)
				_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
				if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
					_, err = f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
					return true, err
				}
			} else {
				if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
					c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
				}
			}

			return true, nil
		}
	}
	return false, nil
}
