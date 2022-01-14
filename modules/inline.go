package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/StalkR/imdb"
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

func InlineQueryHandler(c tb.Context) error {
	query := c.Query().Text
	if query == string("") {
		InlineMainMenu(c)
		return nil
	} else if strings.HasPrefix(query, "google") {
		gsearch_inline(c)
		return nil
	} else if strings.HasPrefix(query, "ud") {
		ud_inline(c)
		return nil
	} else if strings.HasPrefix(query, "imdb") {
		imdb_inline(c)
		return nil
	}
	return nil
}

func InlineMainMenu(c tb.Context) {
	text := "Inline Query Help Menu"
	btns := &tb.InlineKeyboardMarkup{}
	btns.InlineKeyboard = [][]tb.InlineButton{{tb.InlineButton{
		Text:            "Google Search",
		InlineQueryChat: "google ",
	}, tb.InlineButton{Text: "UD Search", InlineQueryChat: "ud "}, tb.InlineButton{Text: "IMDb Search", InlineQueryChat: "imdb "}}}
	result := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: btns}, Title: text, Description: "Here is the inline help menu", Text: text}
	results := make(tb.Results, 1)
	results[0] = result
	results[0].SetResultID("0")
	c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
}

func gsearch_inline(c tb.Context) {
	query := c.Query().Text
	qarg := strings.SplitN(query, " ", 2)
	if len(qarg) == 1 {
		return
	}
	ctx := context.Background()
	search, _ := googlesearch.Search(ctx, qarg[1])
	results := make(tb.Results, len(search))
	for i, r := range search {
		if r.Title != "" {
			text := fmt.Sprintf("<b><a href='%s'>%s</a></b>\n%s", r.URL, r.Title, r.Description)
			rq := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: inline_markup("google"), Content: &tb.InputTextMessageContent{Text: text, DisablePreview: true}}, Title: r.Title, Description: r.Description, ThumbURL: "https://te.legra.ph/file/be8c347e07867d4547c6c.jpg"}
			results[i] = rq
			results[i].SetResultID(strconv.Itoa(i))
		}
	}
	c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
}

func ud_inline(c tb.Context) {
	query := c.Query().Text
	qarg := strings.SplitN(query, " ", 2)
	if len(qarg) == 1 {
		return
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
					text := fmt.Sprintf("<b>Definition for %s:</b>\n\n%s\n\n%s", strings.Title(qarg[1]), defeniton, example)
					rq := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: inline_markup("ud")}, Title: "Defenition " + strconv.Itoa(i+1), Text: text, Description: defeniton.(string), ThumbURL: "https://te.legra.ph/file/658c83f2622fb2237fd82.jpg"}
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
}

func imdb_inline(c tb.Context) {
	q := c.Query().Text
	args := strings.SplitN(q, " ", 2)
	if len(args) == 1 {
		return
	}
	arg := args[1]
	search, _ := imdb.SearchTitle(myClient, arg)
	results := make(tb.Results, 10)
	qb := 0
	for i, result := range search {
		if qb >= 10 {
			break
		}
		btns := &tb.InlineKeyboardMarkup{}
		btns.InlineKeyboard = [][]tb.InlineButton{{tb.InlineButton{Text: result.Name, Data: fmt.Sprintf("imdb_inline_%d", &result.ID)}}, {tb.InlineButton{
			Text:            "Search again",
			InlineQueryChat: "imdb ",
		}}}
		r := &tb.ArticleResult{
			ResultBase:  tb.ResultBase{ReplyMarkup: btns},
			Title:       result.Name,
			Text:        fmt.Sprintf("Click here to view about <b>%s</b>\n<b>Year:</b>%s", result.Name, strconv.Itoa(result.Year)),
			Description: strconv.Itoa(result.Year),
		}
		results[i] = r
		results[i].SetResultID(strconv.Itoa(i))
		qb++
	}
	err := c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
	if err != nil {
		fmt.Println(err)
	}
}
