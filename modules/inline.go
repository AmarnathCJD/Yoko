package modules

import (
	"context"
	"fmt"
	"io"
	"net/http"
        "strings"
        "strconv"
	googlesearch "github.com/rocketlaunchr/google-search"

	"github.com/anaskhan96/soup"
	tb "gopkg.in/tucnak/telebot.v3"
)

func Gg(c tb.Context) error {
	client := &http.Client{}
	url := "https://www.google.com/search?&q=cheems&num=8"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	x, _ := client.Do(request)
	defer x.Body.Close()
	b, err := io.ReadAll(x.Body)
	if err != nil {
		fmt.Println(err)
	}
	doc := soup.HTMLParse(string(b))
	fmt.Println(doc.Find("div", "class"))
	return nil
}

func Test(c tb.Context) error {
	ctx := context.Background()
	fmt.Println(googlesearch.Search(ctx, "cars for sale in Toronto, Canada"))
	return nil
}



func Gsearch_inline(c tb.Context) error {
        query := c.Query().Text
        if !strings.HasPrefix(query, "google"){
         return nil
        }
        qarg := strings.SplitN(query, " ", 2)
        if len(qarg) == 1{
            return nil
        }
        ctx := context.Background()
	search, _ := googlesearch.Search(ctx, qarg[1])
        results := make(tele.Results, len(urls))
        for i, r := range search {
           rq := &tb.ArticleResult{Title: r.Title, Text: "smd", Description: r.Description}
           results[i] = rq
           results[i].SetResultID(strconv.Itoa(i))
        err := c.Bot().Answer(c.Query(), &tb.QueryResponse{
		Results:   results,
		CacheTime: 60
	})
        fmt.Println(err)
	return nil
}
