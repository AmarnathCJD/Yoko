package modules

import (
	"fmt"
	"github.com/amarnathcjd/yoko/modules/db"
	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v3"
	"regexp"
	"strings"
)

var (
	del_all_filters        = sel.Data("Delete all filters", "stopall")
	cancel_del_all_filters = sel.Data("Cancel", "cancel_del_all")
)

func SaveFilter(c tb.Context) error {
	name, note, file := parse_message(c.Message())
	if name == string("") {
		c.Reply("You need to give the filter a name!")
		return nil
	} else if note == string("") && file == nil {
		c.Reply("You need to give the filter a name!")
		return nil
	}
	c.Reply(fmt.Sprintf("Saved filter '%s'.", name))
	db.Save_filter(c.Chat().ID, name, note, file)
	return nil
}

func AllFilters(c tb.Context) error {
	f := db.Get_filters(c.Chat().ID)
	if len(f) == 0 {
		return c.Reply(fmt.Sprintf("No filters in %s!", c.Chat().Title))
	} else {
		txt := fmt.Sprintf("<b>Filters in %s:</b>", c.Chat().Title)
		for _, x := range f {
			txt += fmt.Sprintf("\n- <code>%s</code>", x)
		}
		return c.Reply(txt)
	}
}

func StopFilter(c tb.Context) error {
	if c.Message().Payload == string("") {
		return c.Reply("not enough arguments provided.")
	} else {
		stop := db.Del_filter(c.Chat().ID, c.Message().Payload)
		if !stop {
			return c.Reply("You haven't saved any filters on this word yet!")
		} else {
			return c.Reply(fmt.Sprintf("Filter <code>'%s'</code> has been stopped!", c.Message().Payload))
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
		db.Del_all_filters(c.Chat().ID)
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

func FilterEvent(c tb.Context) (error, bool) {
	f := db.Get_filters(c.Chat().ID)
	if len(f) == 0 {
		return nil, false
	}
	for _, x := range f {
		pattern := `( |^|[^\w])(?i)` + x + `( |$|[^\w])`
		if match, _ := regexp.Match(pattern, []byte(c.Text())); match {
			filter := db.Get_filter(c.Chat().ID, x)
			text, p := ParseString(filter["note"].(string), c)

			if filter["file"] != nil && len(filter["file"].(bson.A)) != 0 && filter["file"].(bson.A)[0] != string("") {
				f := GetFile(filter["file"].(bson.A), text)
				_, err := f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()})
				if err != nil && strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
					_, err = f.Send(c.Bot(), c.Chat(), &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
					return err, true
				}
			} else {

				if err := c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message()}); strings.Contains(err.Error(), "telegram unknown: Bad Request: can't parse entities") {
					c.Send(text, &tb.SendOptions{DisableWebPagePreview: p, ReplyMarkup: btns, ReplyTo: c.Message(), ParseMode: "Markdown"})
				}
			}

			return nil, true
		}
	}
	return nil, false
}
