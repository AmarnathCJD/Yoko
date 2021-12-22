package modules

import (
	"fmt"
	"io"
	"net/http"
        tb "gopkg.in/tucnak/telebot.v3"
	"github.com/anaskhan96/soup"
)

func Gg(c tb.Context) error {
	client := &http.Client{}
	url := "https://www.google.com/search?&q=cheems&num=8"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "[Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)] [Chrome/61.0.3163.100 Safari/537.36]")
	x, _ := client.Do(request)
	defer x.Body.Close()
	b, err := io.ReadAll(x.Body)
	if err != nil {
		fmt.Println(err)
	}
	doc := soup.HTMLParse(string(b))
	fmt.Println(doc.Find("div", "class", "g"))
        return nil
}
