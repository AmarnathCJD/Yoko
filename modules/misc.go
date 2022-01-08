package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/StalkR/imdb"
	tg "github.com/TechMinerApps/telegraph"
	"go.mongodb.org/mongo-driver/bson"
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
	resp, err := myClient.Get("https://api.roseloverx.in/username?username=" + url)
	if err != nil {
		fmt.Println("No response from request")
		return nil, err
	}
	defer resp.Body.Close()
	var t mapType
	json.NewDecoder(resp.Body).Decode(&t)
	return t, err
}

func Info(c tb.Context) error {
	m := c.Message()
	if !m.IsReply() && string(m.Payload) == string("") {
		user_obj := m.Sender
		final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> @%s\n<b>User Link:</b> <a href='tg://user?id=%s'>%s</a>\n\n<b>Gbanned:</b> %s", strconv.Itoa(int(user_obj.ID)), user_obj.FirstName, user_obj.LastName, user_obj.Username, strconv.Itoa(int(user_obj.ID)), "link", "No")
		b.Reply(m, final_msg)
	} else {
		x := strings.SplitN(m.Payload, " ", 2)
		if !m.IsReply() && len(x) > 0 && !isInt(x[0]) {
			u, _ := getJson(strings.TrimPrefix(x[0], "@"))
			u_obj := &tb.User{ID: int64(u["id"].(float64)), Username: u["username"].(string), FirstName: u["first_name"].(string), LastName: u["last_name"].(string)}
			c.Reply(fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%d</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> %s\n<b>DC ID:</b> <code>%d</code>\n<b>User Link:</b> <a href='tg://user?id=%d'>%s</a>\n\n<b>Gbanned:</b> %s", u_obj.ID, u_obj.FirstName, u_obj.LastName, u_obj.Username, int(u["dc_id"].(float64)), u_obj.ID, "link", "No"))
			return nil
		}
		user_obj, _ := get_user(m)
		final_msg := fmt.Sprintf("<b>User info</b>\n<b>ID:</b> <code>%s</code>\n<b>First Name:</b> %s\n<b>Last Name:</b> %s\n<b>Username:</b> @%s\n<b>User Link:</b> <a href='tg://user?id=%s'>%s</a>\n\n<b>Gbanned:</b> %s", strconv.Itoa(int(user_obj.ID)), user_obj.FirstName, user_obj.LastName, user_obj.Username, strconv.Itoa(int(user_obj.ID)), "link", "No")
		b.Reply(m, final_msg)
	}
	return nil
}

func ID_info(c tb.Context) error {
	if !c.Message().IsReply() && string(c.Message().Payload) == string("") {
		c.Reply(fmt.Sprintf("This chat's ID is: <code>%d</code>", c.Chat().ID))
		return nil
	} else {
		user_obj, _ := get_user(c.Message())
		if c.Message().IsForwarded() {
			if c.Message().FromChannel() {
				c.Reply(fmt.Sprintf("User %s's ID is <code>%d</code>.\nThe forwarded channel, %s, has an id of <code>%d</code>.", user_obj.FirstName, user_obj.ID, c.Message().OriginalChat.Title, c.Message().OriginalChat.ID))
				return nil
			} else if c.Message().OriginalSender.ID != user_obj.ID {
				err := c.Reply(fmt.Sprintf("User %s's ID is <code>%d</code>.\nThe forwarded user, %s, has an ID of <code>%d</code>", user_obj.FirstName, user_obj.ID, c.Message().OriginalSenderName, c.Message().OriginalSender.ID))
				fmt.Println(err)
				return nil
			}
		}
		err := c.Reply(fmt.Sprintf("<b>User %s's ID is <code>%d</code>", user_obj.FirstName, user_obj.ID))
		fmt.Println(err)
		return nil
	}
}

func IMDb(c tb.Context) error {
	m := c.Message()
	client := http.DefaultClient
	results, _ := imdb.SearchTitle(client, m.Payload)
	fmt.Println(results)
	title, _ := imdb.NewTitle(client, results[0].ID)
	movie := fmt.Sprintf("<b><u>%s</u></b>\n<b>Type:</b> %s\n<b>Year:</b> %s\n<b>AKA:</b> %s\n<b>Duration:</b> %s\n<b>Rating:</b> %s/10\n<b>Genre:</b> %s\n\n<code>%s</code>\n<b>Source ---> IMDb</b>", title.Name, title.Type, strconv.Itoa(title.Year), title.AKA[0], title.Duration, title.Rating, strings.Join(title.Genres, ", "), title.Description)
	menu.Inline(menu.Row(menu.URL("ImDB", fmt.Sprintf("https://m.imdb.com/title/%s/", title.ID))))
	b.Reply(m, &tb.Photo{File: tb.FromURL(title.Poster.URL), Caption: movie}, menu)
	return nil
}

func Crypto(c tb.Context) error {
	resp, err := myClient.Get("https://api.coingecko.com/api/v3/simple/price?ids=bitcoin%2Clitecoin%2Cdogecoin%2Cbabydoge%2Cethereum%2Cxrp&vs_currencies=usd%2Cinr")
	if err != nil {
		c.Reply(err.Error())
		return nil
	}
	defer resp.Body.Close()
	var r mapType
	json.NewDecoder(resp.Body).Decode(&r)
	crypto := fmt.Sprintf("<b>Crypto Prices</b>\n%s: %d$\n%s: %d$\n%s: %f$\n%s: %d$\n%s: %f$", "Bitcoin", int(r["bitcoin"].(map[string]interface{})["usd"].(float64)), "Ethereum", int(r["ethereum"].(map[string]interface{})["usd"].(float64)), "Dogecoin", r["dogecoin"].(map[string]interface{})["usd"].(float64), "Litecoin", int(r["litecoin"].(map[string]interface{})["usd"].(float64)), "Babydoge", r["babydoge"].(map[string]interface{})["usd"].(float64))
	c.Reply(crypto)
	return nil
}

func Translate(c tb.Context) error {
	text, lang := "", "en"
	if !c.Message().IsReply() && c.Message().Payload == string("") {
		c.Reply("Provide the text to be translated!")
		return nil
	} else if c.Message().IsReply() {
		text = c.Message().ReplyTo.Text
		if c.Message().Payload != string("") {
			lang = strings.SplitN(c.Message().Payload, " ", 2)[0]
		}
	} else if c.Message().Payload != string("") {
		args := strings.SplitN(c.Message().Payload, " ", 2)
		if len(args) == 2 && len([]rune(args[0])) == 2 {
			lang, text = args[0], args[1]
		} else {
			text = c.Message().Payload
		}
	}
	url_d := "https://script.google.com/macros/s/AKfycbzFXVfjwX_RB6XkjLpwlMIXl_IVeoqaYnfhRf774xknBAcV00Ef3OPK89uS7TBFppwfVg/exec"
	data := url.Values{"text": {text}, "source": {""}, "target": {lang}}
	rq, err := http.PostForm(url_d, data)
	if err != nil {
		c.Reply(err.Error())
		return nil
	}
	defer rq.Body.Close()
	var r mapType
	json.NewDecoder(rq.Body).Decode(&r)
	translated := fmt.Sprintf("<b>translated to %s:</b>\n<code>%s</code>", lang, r["result"].(string))
	c.Reply(translated)
	return nil
}

func Ud(c tb.Context) error {
	api := fmt.Sprint("http://api.urbandictionary.com/v0/define?term=", c.Message().Payload)
	resp, _ := myClient.Get(api)
	var v mapType
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&v)
	res := v["list"].([]interface{})
	fmt.Println(res)
	if len(res) == 0 {
		b.Reply(c.Message(), "No results found.")
		return nil
	}
	b.Reply(c.Message(), fmt.Sprintf("<b>%s:</b>\n\n%s\n\n<i>%s</i>", c.Message().Payload, res[0].(map[string]interface{})["definition"], res[0].(map[string]interface{})["example"]))
	return nil
}

func Bin_check(c tb.Context) error {
	bin := c.Message().Payload
	url := "https://lookup.binlist.net/%s"
	resp, _ := http.Get(fmt.Sprintf(url, bin))
	var v bson.M
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&v)
	country := v["country"].(map[string]interface{})
	bank := v["bank"].(map[string]interface{})
	out_str := fmt.Sprintf("<b>BIN/IIN:</b> <code>%s</code> %s", bin, country["emoji"])
	if scheme, f := v["scheme"]; f {
		out_str += fmt.Sprintf("\n<b>Card Brand:</b> %s", strings.Title(scheme.(string)))
	}
	if ctype, f := v["type"]; f {
		out_str += fmt.Sprintf("\n<b>Card Type:</b> %s", strings.Title(ctype.(string)))
	}
	if brand, f := v["brand"]; f {
		out_str += fmt.Sprintf("\n<b>Card Level:</b> %s", strings.Title(brand.(string)))
	}
	if prepaid, f := v["prepaid"]; f {
		out_str += fmt.Sprintf("\n<b>Prepaid:</b> %s", strings.Title(strconv.FormatBool(prepaid.(bool))))
	}
	if name, f := bank["name"]; f {
		out_str += fmt.Sprintf("\n<b>Bank:</b> %s", strings.Title(name.(string)))
	}
	if ctry, f := country["name"]; f {
		out_str += fmt.Sprintf("\n<b>Country:</b> %s - %s - $%s", strings.Title(ctry.(string)), country["alpha2"], country["currency"])
	}
	if url, f := bank["url"]; f {
		out_str += fmt.Sprintf("\n<b>Website:</b> <code>%s</code>", url)
	}
	out_str += "\n<b>━━━━━━━━━━━━━</b>"
	out_str += fmt.Sprintf("\nChecked by <a href='tg://user?id=%s'>%s</a>", strconv.Itoa(int(c.Message().Sender.ID)), c.Message().Sender.FirstName)
	c.Reply(out_str)
	return nil
}

func telegraph(c tb.Context) error {
	text := c.Message().Payload
	title := time.Now().Format("01-02-2006 15:04:05 Monday")
	if c.Message().IsReply() {
		if c.Message().ReplyTo.Text != string("") {
			text = c.Message().ReplyTo.Text
			if c.Message().Payload != string("") {
				title = c.Message().Payload
			}
		} else if c.Message().ReplyTo.Document != nil {
			c.Bot().Download(&c.Message().ReplyTo.Document.File, "doc.txt")
			data, err := ioutil.ReadFile("doc.txt")
			if err != nil {
				c.Reply(err.Error())
				return nil
			} else {
				text = string(data)
				if c.Message().Payload != string("") {
					title = c.Message().Payload
				}
			}
			os.Remove("doc.txt")
		}
	}
	if text == string("") {
		c.Reply("No text was provided!")
		return nil
	}
	client := tg.NewClient()
	client.CreateAccount(tg.Account{ShortName: "Yoko", AuthorName: c.Sender().FirstName})
	content, _ := client.ContentFormat(text)
	pageData := tg.Page{
		Title:   title,
		Content: content,
	}
	page, err := client.CreatePage(pageData, true)
	fmt.Println(err)
	menu.Inline(menu.Row(menu.URL("Telegraph", page.URL)))
	c.Reply(fmt.Sprintf("Pasted to <a href='%s'>Tele.graph.org</a>!", page.URL), &tb.SendOptions{DisableWebPagePreview: true, ReplyMarkup: menu})
	return nil
}

func Math(c tb.Context) error {
	query := c.Message().Payload
	if query == string("") {
		c.Reply("Please provide the Mathamatical Equation.")
		return nil
	} else {
		endpoint := "https://evaluate-expression.p.rapidapi.com/expression"
		req, _ := http.NewRequest("GET", endpoint, nil)
		req.Header.Add("x-rapidapi-key", "fef481fee3mshf99983bfc650decp104100jsnbad6ddb2c846")
		req.Header.Add("x-rapidapi-host", "evaluate-expression.p.rapidapi.com")
		q := req.URL.Query()
		q.Add("expression", query)
		req.URL.RawQuery = q.Encode()
		fmt.Println(req.URL.String())
		r, err := myClient.Do(req)
		if err != nil {
			c.Reply(err.Error())
		}
		defer r.Body.Close()
                var body mapType
                json.NewDecoder(r.Body).Decode(&body)
		fmt.Println(body)
	}
	return nil
}

func Paste(c tb.Context) error {
	uri := "https://api.roseloverx.in/paste"
	req, _ := http.NewRequest("GET", uri, nil)
	q := req.URL.Query()
	q.Add("text", c.Message().Payload)
	req.URL.RawQuery = q.Encode()
        fmt.Println(req.URL.String())
	resp, err := myClient.Do(req)
	if err != nil {
		c.Reply(err.Error())
		return nil
	}
	var body mapType
	json.NewDecoder(resp.Body).Decode(&body)
        sel.Inline(sel.Row(sel.URL("View Paste", fmt.Sprintf("https://nekobin.com/%s", body["result"].(map[string]interface{})["key"].(string)))))
	c.Reply(fmt.Sprintf("Pasted to <b><a href='https://nekobin.com/%s'>NekoBin</a></b>.", body["result"].(map[string]interface{})["key"].(string)), &tb.SendOptions{DisableWebPagePreview: true, ReplyMarkup: sel})
	return nil
}
