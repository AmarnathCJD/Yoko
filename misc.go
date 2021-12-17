package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/StalkR/imdb"
	tb "gopkg.in/tucnak/telebot.v3"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func get_user(m *tb.Message) (*tb.User, string) {
	if m.IsReply() {
		user_obj := m.ReplyTo.Sender
		if len(m.Payload) != 0 {
			return user_obj, m.Payload
		} else {
			return user_obj, ""
		}
	} else if len(m.Payload) != 0 {
		x := strings.SplitN(m.Payload, " ", 2)
		if isInt(x[0]) {
			user_id, _ := strconv.ParseInt(x[0], 10, 64)
			user_obj, err := b.ChatByID(user_id)
			if err != nil {
				b.Reply(m, "Looks like I don't have control over that user, or the ID isn't a valid one. If you reply to one of their messages, I'll be able to interact with them.")
				return nil, ""
			}
			user := &tb.User{ID: int64((user_obj.ID)), FirstName: user_obj.FirstName, LastName: user_obj.LastName, Username: user_obj.Username}
			if len(x) > 1 {
				return user, x[1]
			} else {
				return user, ""
			}
		} else {
			u, err := getJson(strings.TrimPrefix(x[0], "@"))
			if err != nil {
				b.Reply(m, fmt.Sprint(err.Error()))
				return nil, ""
			}
			user_obj := &tb.User{ID: int64(u["id"].(float64)), Username: u["username"].(string), FirstName: u["first_name"].(string), LastName: u["last_name"].(string)}
			if len(x) > 1 {
				return user_obj, x[1]
			} else {
				return user_obj, ""
			}
		}
	} else {
		b.Reply(m, "You dont seem to be referring to a user or the ID specified is incorrect..")
		return nil, ""
	}
}

type mapType map[string]interface{}

func getJson(url string) (mapType, error) {
	resp, err := myClient.Get("https://roseflask.herokuapp.com/username?username=" + url)
	if err != nil {
		fmt.Println("No response from request")
		return nil, err
	}
	defer resp.Body.Close()
	var t mapType
	json.NewDecoder(resp.Body).Decode(&t)
	return t, err
}

func info(c tb.Context) error {
	m := c.Message()
	if !m.IsReply() && string(m.Payload) == string("") {
		user_obj := m.Sender
		final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> @%s\n<b>User Link:</b> <a href='tg://user?id=%s'>%s</a>\n\n<b>Gbanned:</b> %s", strconv.Itoa(int(user_obj.ID)), user_obj.FirstName, user_obj.LastName, user_obj.Username, strconv.Itoa(int(user_obj.ID)), "link", "No")
		b.Reply(m, final_msg)
	} else {
		user_obj, _ := get_user(m)
		final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> @%s\n<b>User Link:</b> <a href='tg://user?id=%s'>%s</a>\n\n<b>Gbanned:</b> %s", strconv.Itoa(int(user_obj.ID)), user_obj.FirstName, user_obj.LastName, user_obj.Username, strconv.Itoa(int(user_obj.ID)), "link", "No")
		b.Reply(m, final_msg)
	}
	return nil
}

func gp(m *tb.Message) {
	u, _ := get_user(m)
	x, err := b.ChatMemberOf(m.Chat, u)
	fmt.Println(x.Rights)
	if err != nil {
		b.Reply(m, string(err.Error()))
		return
	}
	b.Reply(m, fmt.Sprint(x.Rights))
}

func IMDb(c tb.Context) error {
	m := c.Message()
	client := http.DefaultClient
	results, _ := imdb.SearchTitle(client, m.Payload)
	title, _ := imdb.NewTitle(client, results[0].ID)
	movie := fmt.Sprintf("<b><u>%s</u></b>\n<b>Type:</b> %s\n<b>Year:</b> %s\n<b>AKA:</b> %s\n<b>Duration:</b> %s\n<b>Rating:</b> %s/10\n<b>Genre:</b> %s\n\n<code>%s</code>\n<b>Source ---> IMDb</b>", title.Name, title.Type, strconv.Itoa(title.Year), title.AKA[0], title.Duration, title.Rating, strings.Join(title.Genres, ", "), title.Description)
	menu.Inline(menu.Row(menu.URL("ImDB", fmt.Sprintf("https://m.imdb.com/title/%s/", title.ID))))
	b.Reply(m, &tb.Photo{File: tb.FromURL(title.Poster.URL), Caption: movie}, menu)
	return nil
}

func Crypto(c tb.Context) error {
	m := c.Message()
	resp, _ := myClient.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin%2Clitecoin%2Cdogecoin%2Cbabydoge%2Cethereum%2Cxrp&vs_currencies=usd%2Cinr")
	defer resp.Body.Close()
	var r mapType
	json.NewDecoder(resp.Body).Decode(&r)
	crypto := fmt.Sprintf("<b>Crypto Prices</b>\n%s: %d$\n%s: %d$\n%s: %f$\n%s: %d$\n%s: %f$", "Bitcoin", int(r["bitcoin"].(map[string]interface{})["usd"].(float64)), "Ethereum", int(r["ethereum"].(map[string]interface{})["usd"].(float64)), "Dogecoin", r["dogecoin"].(map[string]interface{})["usd"].(float64), "Litecoin", int(r["litecoin"].(map[string]interface{})["usd"].(float64)), "Babydoge", r["babydoge"].(map[string]interface{})["usd"].(float64))
	b.Reply(m, crypto)
	return nil
}

func translate(c tb.Context) error {
	m := c.Message()
	text, lang := "", "en"
	if !m.IsReply() && m.Payload == string("") {
		b.Reply(m, "Provide the text to be translated!")
		return nil
	} else if m.IsReply() {
		text = m.ReplyTo.Text
		if m.Payload != string("") {
			lang = strings.SplitN(m.Payload, " ", 2)[0]
		}
	} else if m.Payload != string("") {
		args := strings.SplitN(m.Payload, " ", 2)
		if len(args) == 2 && len([]rune(args[0])) == 2 {
			lang, text = args[0], args[1]
		} else {
			text = m.Payload
		}
	}
	url_d := "https://script.google.com/macros/s/AKfycbzFXVfjwX_RB6XkjLpwlMIXl_IVeoqaYnfhRf774xknBAcV00Ef3OPK89uS7TBFppwfVg/exec"
	data := url.Values{"text": {text}, "source": {""}, "target": {lang}}
	rq, _ := http.PostForm(url_d, data)
	defer rq.Body.Close()
	var r mapType
	json.NewDecoder(rq.Body).Decode(&r)
	translated := fmt.Sprintf("<b>translated to %s:</b>\n<code>%s</code>", lang, r["result"].(string))
	b.Reply(m, translated)
	return nil
}

func uD(c tb.Context) error {
	api := fmt.Sprint("http://api.urbandictionary.com/v0/define?term=%s", c.Message().Payload)
	resp, _ := myClient.Get(api)
	var v mapType
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&v)
	fmt.Println(v)
}
