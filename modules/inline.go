package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/StalkR/imdb"
	googlesearch "github.com/rocketlaunchr/google-search"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	tb "gopkg.in/telebot.v3"
)

var imdb_btn = sel.Data("imdb_in", "imdb_inline")

func inline_markup(query string) *tb.ReplyMarkup {
	btns := &tb.ReplyMarkup{}
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
	} else if strings.HasPrefix(query, "yt") {
		yt_search(c)
		return nil
	} else if strings.HasPrefix(query, "insta") {
		return InstaGramSearch(c)
	}
	return nil
}

func InlineMainMenu(c tb.Context) {
	text := "Inline Query Help Menu"
	btns := &tb.ReplyMarkup{}
	btns.InlineKeyboard = [][]tb.InlineButton{{tb.InlineButton{
		Text:            "Google Search",
		InlineQueryChat: "google ",
	}, tb.InlineButton{Text: "UD Search", InlineQueryChat: "ud "}}, {tb.InlineButton{Text: "IMDb Search", InlineQueryChat: "imdb "}, tb.InlineButton{Text: "Youtube Search", InlineQueryChat: "yt "}}}
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
		Results:           results,
		CacheTime:         60,
		SwitchPMText:      "Mika",
		SwitchPMParameter: "start",
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
		btns := &tb.ReplyMarkup{}
		btns.InlineKeyboard = [][]tb.InlineButton{{tb.InlineButton{Text: result.Name, Data: fmt.Sprintf("imdb_inline_%s", result.ID), Unique: "imdb_inline"}}, {tb.InlineButton{
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

func ImdbCB(c tb.Context) error {
	d := strings.Split(c.Callback().Data, "_")
	title, err := imdb.NewTitle(myClient, d[2])
	check(err)
	movie := fmt.Sprintf("<b><u>%s</u></b>\n<b>Type:</b> %s\n<b>Year:</b> %s\n<b>AKA:</b> %s\n<b>Duration:</b> %s\n<b>Rating:</b> %s/10\n<b>Genre:</b> %s\n\n<code>%s</code>\n<b>Source ---> IMDb</b>", title.Name, title.Type, strconv.Itoa(title.Year), title.AKA[0], title.Duration, title.Rating, strings.Join(title.Genres, ", "), title.Description)
	sel.Inline(sel.Row(sel.QueryChat("Search again", "imdb ")))
	return c.Edit(movie, &tb.SendOptions{ParseMode: tb.ModeHTML, ReplyMarkup: sel})
}

func yt_search(c tb.Context) {
	q := c.Query().Text
	args := strings.SplitN(q, " ", 2)
	if len(args) == 1 {
		return
	}
	arg := args[1]
	client, _ := youtube.NewService(context.TODO(), option.WithAPIKey(YOUTUBE_API_KEY))
	call := client.Search.List([]string{"snippet"}).Q(arg).MaxResults(10)
	resp, err := call.Do()
	check(err)
	results := make(tb.Results, 10)
	for i, x := range resp.Items {
		text := fmt.Sprintf("<b><a href='https://www.youtube.com/watch?v=%s'>%s</a></b>\n<b>%s</b>\n<b>Published:</b> %s\n\n <i>%s</i>", x.Id.VideoId, x.Snippet.Title, x.Snippet.ChannelTitle, x.Snippet.PublishedAt, x.Snippet.Description)
		r := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: inline_markup("yt"), Content: &tb.InputTextMessageContent{Text: text, DisablePreview: false}}, Title: x.Snippet.Title, Description: x.Snippet.ChannelTitle, ThumbURL: x.Snippet.Thumbnails.Medium.Url}
		results[i] = r
		results[i].SetResultID(strconv.Itoa(i))
	}
	err = c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
	check(err)
}

func InstaGramSearch(c tb.Context) error {
	q := c.Query().Text
	args := strings.SplitN(q, " ", 2)
	if len(args) == 1 {
		return nil
	}
	arg := args[1]
	InstaUrl := `https://www.instagram.com/web/search/topsearch/?context=blended&query=` + arg + `&rank_token=0.29359664635546645&include_reel=true`
	req, _ := http.NewRequest("GET", InstaUrl, nil)
	req.Header.Add("cookie", InstagramCookies)
	resp, err := Client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	var data InstSearch
	json.NewDecoder(resp.Body).Decode(&data)
	Len := len(data.Users)
	if Len > 10 {
		Len = 10
	}
	results := make(tb.Results, Len)
	qd := 0
	for i, x := range data.Users {
		if qd > 10 {
			break
		}
		qd++
		text := ""
		User := x.User
		if User.Username != "" {
			text += fmt.Sprintf("<b>@%s</b>\n", User.Username)
		}
		if User.FullName != "" {
			text += fmt.Sprintf("<b>FullName</b>%s\n", User.FullName)
		}
		text += fmt.Sprintf("<b>Pk</b>%s\n", User.Pk)
		text += fmt.Sprintf("<b>Verified</b>%v\n", User.IsVerified)
		text += fmt.Sprintf("<b>Private</b>%v\n", User.IsPrivate)
		text += fmt.Sprintf("<b>HasAnonymousProfilePicture</b>%v\n", User.HasAnonymousProfilePicture)
		text += fmt.Sprintf("<b>Seen</b>%d\n", User.Seen)
		text += fmt.Sprintf("<b>HasHighlightReels</b>%v\n", User.HasHighlightReels)
		r := &tb.ArticleResult{ResultBase: tb.ResultBase{ReplyMarkup: inline_markup("insta"), Content: &tb.InputTextMessageContent{Text: text, DisablePreview: false}}, Title: User.Username, Description: User.FullName}
		results[i] = r
		results[i].SetResultID(strconv.Itoa(i))
	}
	return c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60,
	})
}
