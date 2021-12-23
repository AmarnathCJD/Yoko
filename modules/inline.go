package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	googlesearch "github.com/rocketlaunchr/google-search"
	tb "gopkg.in/tucnak/telebot.v3"
)

func inline_markup(query string) *tb.InlineKeyboardMarkup {
	btns := &tb.InlineKeyboardMarkup{}
	btns.InlineKeyboard = [][]tb.InlineButton{{tb.InlineButton{
		Text:            "Search again",
		InlineQueryChat: query,
	}}}
	return btns
}

func gsearch_inline(c tb.Context) error {
	query := c.Query().Text
	if !strings.HasPrefix(query, "google") {
		return nil
	}
	qarg := strings.SplitN(query, " ", 2)
	if len(qarg) == 1 {
		return nil
	}
	ctx := context.Background()
	search, _ := googlesearch.Search(ctx, qarg[1])
	results := make(tb.Results, len(search))
	for i, r := range search {
		text := fmt.Sprintf("<b><a href='%s'>%s</a></b>\n%s", r.URL, r.Title, r.Description)
		rq := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: inline_markup("google")}, Title: r.Title, Text: text, Description: r.Description, HideURL: true, ThumbURL: "https://te.legra.ph/file/be8c347e07867d4547c6c.jpg"}
		results[i] = rq
		results[i].SetResultID(strconv.Itoa(i))
	}
	c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
	return nil
}

func ud_inline(c tb.Context) error {
	query := c.Query().Text
	if !strings.HasPrefix(query, "ud") {
		return nil
	}
	qarg := strings.SplitN(query, " ", 2)
	if len(qarg) == 1 {
		return nil
	}
	api := fmt.Sprint("http://api.urbandictionary.com/v0/define?term=", qarg[1])
	resp, _ := myClient.Get(api)
	var v mapType
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&v)
	res := v["list"].([]interface{})
	results := make(tb.Results, len(res))
	for i, r := range res {
		if hq, ok := r.(map[string]interface{}); ok {
			if defeniton, ok := hq["definition"]; ok {
				if example, ok := hq["example"]; ok {
					if len(results) == 3 {
						break
					}
					text := fmt.Sprintf("<b>%s:</b>\n\n%s\n\n<i>%s</i>", strings.Title(qarg[1]), defeniton, example)
					rq := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: inline_markup("ud")}, Title: "Result " + strconv.Itoa(i), Text: text, Description: defeniton.(string)}
					results[i] = rq
					results[i].SetResultID(strconv.Itoa(i))
				}
			}
		}
	}
	c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
	return nil
}
