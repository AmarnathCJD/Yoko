package modules

import (
	"bytes"
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
        "fmt"
	"github.com/StalkR/imdb"
	tg "github.com/TechMinerApps/telegraph"
	"github.com/anaskhan96/soup"
	yt "github.com/kkdai/youtube/v2"
	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/telebot.v3"
)

func UserInfo(c tb.Context) error {
	var u User
	if !c.Message().IsReply() && c.Message().Payload == string("") {
		if c.Sender().ID == 136817688 {
			SenderChat := c.Message().SenderChat
			u = User{
				ID:       SenderChat.ID,
				First:    EscapeHTML(SenderChat.Title),
				Last:     "",
				Username: "@" + SenderChat.Username,
				Mention:  "",
				DC:       0,
				Type:     "chat",
			}
		} else {
			Sender := c.Sender()
			u = User{
				ID:       Sender.ID,
				First:    EscapeHTML(Sender.FirstName),
				Last:     EscapeHTML(Sender.LastName),
				Type:     "user",
				Username: "@" + Sender.Username,
				Mention:  GetMention(Sender.ID, Sender.FirstName),
				DC:       0,
			}

		}
	} else {
		u, _ = GetUser(c)
	}
        fmt.Println(u)
	if u.ID == 0 {
		return nil
	}
	Info := ""
	if u.Type == "chat" {
		Info += "<b>Channel Info</b>"
	} else {
		Info += "<b>User Info</b>"
	}
	Info += fmt.Sprintf("\n<b>ID:</b> <code>%d</code>", u.ID)
	if u.First != string("") {
		if u.Type == "chat" {
			Info += fmt.Sprintf("\n<b>Title:</b> %s", u.First)
		} else {
			Info += fmt.Sprintf("\n<b>FirstName:</b> %s", u.First)
		}
	}
	if u.Last != string("") {
		Info += fmt.Sprintf("\n<b>LastName:</b> %s", u.Last)
	}
	if u.Username != string("") {
		Info += fmt.Sprintf("\n<b>Username:</b> %s", u.Username)
	}
	if u.DC != 0 {
		Info += fmt.Sprintf("\n<b>DC ID:</b> <code>%d</code>", u.DC)
	}
	if u.Type != "chat" {
		Info += fmt.Sprintf("\n<b>User Link:</b> %s", u.Mention)
	}
	Info += "\n\n<b>Gbanned:</b> No"
	return c.Reply(Info)

}

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
			if e, ok := u["error"]; ok {
				b.Reply(m, e.(string))
				return nil, ""
			}
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
	resp, err := myClient.Get("https://polar-refuge-17864.herokuapp.com/username?username=" + url)
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
			if u == nil {
				return nil
			}
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
		err := c.Reply(fmt.Sprintf("User %s's ID is <code>%d</code>", user_obj.FirstName, user_obj.ID))
		fmt.Println(err)
		return nil
	}
}

func ChatInfo(c tb.Context) error {
	var chat *tb.Chat
	if c.Message().IsReply() && c.Message().ReplyTo.FromChannel() {
		chat_id := c.Message().ReplyTo.SenderChat.ID
		chat, _ = c.Bot().ChatByID(chat_id)
	} else if c.Message().Payload != string("") {
		if isInt(c.Message().Payload) {
			chat_, _ := strconv.Atoi(c.Message().Payload)
			chat, _ = c.Bot().ChatByID(int64(chat_))
		} else {
			chat, _ = c.Bot().ChatByUsername(c.Message().Payload)
		}
	} else {
		chat, _ = c.Bot().ChatByID(c.Chat().ID)
	}
	if chat != nil {
		msg := fmt.Sprintf("<b>Chat info</b>\n<b>ID:</b> <code>%d</code>\n<b>Title:</b> %s", chat.ID, chat.Title)
		if chat.Username != "" {
			msg += fmt.Sprintf("\n<b>Username:</b> @%s", chat.Username)
		}
		msg += fmt.Sprintf("\n<b>Link:</b> <a href='tg://resolve?domain=%s'>%s</a>", chat.Username, "link")
		if chat.Description != "" {
			msg += fmt.Sprintf("\n<b>Description:</b> <code>%s</code>", chat.Description)
		}
		if chat.LinkedChatID != 0 {
			msg += fmt.Sprintf("\n<b>Linked Chat ID:</b> %d", chat.LinkedChatID)
		}
		if chat.InviteLink != "" {
			msg += fmt.Sprintf("\n<b>Invite Link:</b> <a href='%s'>%s</a>", chat.InviteLink, "link")
		}
		c.Reply(msg, &tb.SendOptions{DisableWebPagePreview: true})
		return nil
	} else {
		c.Reply("Invalid chat")
		return nil
	}
}

func IMDb(c tb.Context) error {
	client := http.DefaultClient
	results, _ := imdb.SearchTitle(client, c.Message().Payload)
	if len(results) == 0 {
		return c.Reply("No results found!")
	}
	title, _ := imdb.NewTitle(client, results[0].ID)
	movie := fmt.Sprintf("<b><u>%s</u></b>\n<b>Type:</b> %s\n<b>Year:</b> %s\n<b>AKA:</b> %s\n<b>Duration:</b> %s\n<b>Rating:</b> %s/10\n<b>Genre:</b> %s\n\n<code>%s</code>\n<b>Source ---> IMDb</b>", title.Name, title.Type, strconv.Itoa(title.Year), title.AKA[0], title.Duration, title.Rating, strings.Join(title.Genres, ", "), title.Description)
	menu.Inline(menu.Row(menu.URL("ImDB", fmt.Sprintf("https://m.imdb.com/title/%s/", title.ID))))
	return c.Reply(&tb.Photo{File: tb.FromURL(title.Poster.URL), Caption: movie}, menu)
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
	client.CreateAccount(tg.Account{ShortName: "mika", AuthorName: c.Sender().FirstName})
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
		url := "https://evaluate-expression.p.rapidapi.com"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-rapidapi-host", "evaluate-expression.p.rapidapi.com")
		req.Header.Add("x-rapidapi-key", "cf9e67ea99mshecc7e1ddb8e93d1p1b9e04jsn3f1bb9103c3f")
		q := req.URL.Query()
		q.Add("expression", c.Message().Payload)
		req.URL.RawQuery = q.Encode()
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			c.Reply(err.Error())
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) != string("") {
			c.Reply(string(body))
		}
	}
	return nil
}

func Paste(c tb.Context) error {
	text := c.Message().Payload
	if c.Message().IsReply() {
		if c.Message().ReplyTo.Text != string("") {
			text = c.Message().ReplyTo.Text
		} else if c.Message().ReplyTo.Document != nil {
			c.Bot().Download(&c.Message().ReplyTo.Document.File, "doc.txt")
			data, err := ioutil.ReadFile("doc.txt")
			if err != nil {
				c.Reply(err.Error())
				return nil
			} else {
				text = string(data)
			}
			os.Remove("doc.txt")
		}
	}
	if text == string("") {
		c.Reply("Give some text to paste it!")
		return nil
	}
	if strings.Contains(c.Message().Payload, "-h") {
		uri := "https://www.toptal.com/developers/hastebin/documents"
		req, _ := http.NewRequest("POST", uri, bytes.NewBufferString(strings.ReplaceAll(text, "-h", "")))
		r, err := myClient.Do(req)
		check(err)
		defer r.Body.Close()
		var bd bson.M
		json.NewDecoder(r.Body).Decode(&bd)
		key, sucess := bd["key"]
		if !sucess {
			c.Reply("HasteBin Down.")
			return nil
		} else {
			key = key.(string)
		}
		URL := fmt.Sprintf("https://www.toptal.com/developers/hastebin/%s", key)
		sel.Inline(sel.Row(sel.URL("View Paste", URL)))
		c.Reply(fmt.Sprintf("<b>Pasted to <a href='%s'>HasteBin</a></b>", URL), sel)
		return nil
	}
	postBody, _ := json.Marshal(map[string]string{
		"content": text,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("https://warm-anchorage-15807.herokuapp.com/api/documents", "application/json", responseBody)
	check(err)
	defer resp.Body.Close()
	var body mapType
	json.NewDecoder(resp.Body).Decode(&body)
	sel.Inline(sel.Row(sel.URL("View Paste", fmt.Sprintf("https://warm-anchorage-15807.herokuapp.com/%s", body["result"].(map[string]interface{})["key"].(string)))))
	c.Reply(fmt.Sprintf("Pasted to <b><a href='https://warm-anchorage-15807.herokuapp.com/%s'>NekoBin</a></b>.", body["result"].(map[string]interface{})["key"].(string)), &tb.SendOptions{DisableWebPagePreview: true, ReplyMarkup: sel})
	return nil
}

func Fake_gen(c tb.Context) error {
	cq := Parse_country(strings.TrimSpace(c.Message().Payload))
	if cq == string("") {
		return nil
	}
	url := "https://www.fakexy.com/fake-address-generator-" + cq
	resp, err := soup.Get(url)
	if err != nil {
		return nil
	}
	doc := soup.HTMLParse(resp)
	tbody := doc.FindAll("tbody")
	fname, bday, street, city, country, state, zip, num, gen := "", "", "", "", "", "", "", "", ""
	for _, x := range tbody {
		for _, y := range x.FindAll("tr") {
			d := y.FindAll("td")
			if d[0].Text() == "Street" {
				street = d[1].Text()
			} else if d[0].Text() == "City/Town" {
				city = d[1].Text()
			} else if d[0].Text() == "State/Province/Region" {
				state = d[1].Text()
			} else if d[0].Text() == "Zip/Postal Code" {
				zip = d[1].Text()
			} else if d[0].Text() == "Phone Number" {
				num = d[1].Text()
			} else if d[0].Text() == "Birthday" {
				bday = d[1].Text()
			} else if d[0].Text() == "Gender" {
				gen = d[1].Text()
			} else if d[0].Text() == "Full Name" {
				fname = d[1].Text()
			} else if d[0].Text() == "Country" {
				country = d[1].Text()
			}
		}
	}
	msg := fmt.Sprintf("<b>Fake Address Gen</b>\n<b>Full Name:</b> <code>%s</code>\n<b>Gender:</b> <code>%s</code>\n<b>DOB:</b> <code>%s</code>\n<b>Street:</b> <code>%s</code>\n<b>City/Town:</b> <code>%s</code>\n<b>State:</b> <code>%s</code>\n<b>Zip:</b> <code>%s</code>\n<b>Country:</b> <code>%s</code>\n<b>Phone Number:</b> <code>%s</code>", fname, gen, bday, street, city, state, zip, country, num)
	c.Reply(msg)
	return nil
}

func YT_search(c tb.Context) error {
	return nil
}

func GroupStat(c tb.Context) error {
	return c.Reply(fmt.Sprintf("<b>Total Messages in %s:</b> <code>%d</code>", c.Chat().Title, c.Message().ID))
}

func WebSS(c tb.Context) error {
	query := c.Message().Payload
	body := strings.NewReader(fmt.Sprintf("url=%s&cookies=0&proxy=0&delay=0&captchaToken=false&device=1&platform=1&browser=1&fFormat=1&width=1280&height=800&uid=NaN", url.QueryEscape(query)))
	req, err := http.NewRequest("POST", "https://onlinescreenshot.com/", body)
	check(err)
	req.Header.Set("Authority", "onlinescreenshot.com")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"96\", \"Opera\";v=\"82\", \";Not A Brand\";v=\"99\"")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 OPR/82.0.4227.58")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
	req.Header.Set("Origin", "https://onlinescreenshot.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://onlinescreenshot.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	resp, err := http.DefaultClient.Do(req)
	check(err)
	defer resp.Body.Close()
	var res mapType
	json.NewDecoder(resp.Body).Decode(&res)
	if img_url, ok := res["imgUrl"]; ok {
		if _, ok := img_url.(bool); ok {
			c.Reply(res["msg"].(string))
		}
		return c.Reply(&tb.Photo{File: tb.FromURL(img_url.(string))})
	}
	return nil
}

func Tr2(c tb.Context) error {
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
		args := strings.SplitN(c.Message().Text, " ", 3)
		if len([]rune(args[1])) == 2 {
			lang = args[1]
			text = args[2]
		} else {
			text = args[1] + " " + args[2]
		}
	}
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`async=translate,sl:auto,tl:%s,st:%s,id:1643102010421,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc`, lang, url.QueryEscape(text)))
	req, _ := http.NewRequest("POST", "https://www.google.com/async/translate?vet=12ahUKEwiM3pvpx8z1AhV_SmwGHRb5C5MQqDh6BAgDECY..i&ei=EL_vYYyWFP-UseMPlvKvmAk&client=opera&yv=3", data)
	req.Header.Set("authority", "www.google.com")
	req.Header.Set("sec-ch-ua", `"Opera";v="83", "Chromium";v="97", ";Not A Brand";v="99"`)
	req.Header.Set("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36 OPR/83.0.4254.19")
	req.Header.Set("sec-ch-ua-arch", `"x86"`)
	req.Header.Set("sec-ch-ua-full-version", `"97.0.4692.71"`)
	req.Header.Set("sec-ch-ua-platform-version", `"5.13.0"`)
	req.Header.Set("sec-ch-ua-bitness", `"64"`)
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://www.google.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://www.google.com/")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "SEARCH_SAMESITE=CgQIuJQB; OTZ=6326456_34_34__34_; HSID=ATa13Uw3JpMJmWA3t; SSID=AOEkIbxbQxhvi1FY3; APISID=rdsFU1YTbgq0B3E-/AjkLEBu-qaec_yvgN; SAPISID=-fU4gGX9wHh-Plxb/A_gvZWiONzjK_xLc6; __Secure-1PAPISID=-fU4gGX9wHh-Plxb/A_gvZWiONzjK_xLc6; __Secure-3PAPISID=-fU4gGX9wHh-Plxb/A_gvZWiONzjK_xLc6; SID=GAjUGBrrRyEllUAh04TJFwG4UKCvWjg7c9IZNv-jwJUf6MGArEHHWkJnI71PGYs6d60-Tg.; __Secure-1PSID=GAjUGBrrRyEllUAh04TJFwG4UKCvWjg7c9IZNv-jwJUf6MGAfVw1akWNyBXiDczCq91ttQ.; __Secure-3PSID=GAjUGBrrRyEllUAh04TJFwG4UKCvWjg7c9IZNv-jwJUf6MGAvc-8vuvdO7JYDf0vkP95zg.; 1P_JAR=2022-01-25-09; DV=Y7jy1785Mz1PUOUcLYCoi47rniUI6RcvSfGgakoo6QAAAGCqGssnH9E8zQAAAPg2dSf3vHJGVwAAAA; NID=511=m_HvcK6BB_kHXAzPUuyjqfb0UwSZwalTj5paM9hr2P2EkonwyUIGZSQA7ConYzeH9J4YFCI-nkCZgSMnwv7XTUrcnI8Y4yRx8L65nX7vtL-1fGk_6xl5s5iTgWABhH45EDx42PKUBT1WkL3MeYqcx45-KOMff3brrvu2aYVr3litCGralFYl6lL12MepW9Rd-o-vgGZc_991llxxl3T9Nfs1iteD2w1vg8Ccaha9e2I8Sw7DVGSfuis2YyOact5jD9kf3kvGvjSlT6bMkM7s1s_QvGMeMePiVXvGxzmYoYd5IFhhdHTiJV4PLUxW2K-Nw7Bd-6Il; SIDCC=AJi4QfGW8KIy7dxF647EtoaG4uvUHqFYuyzg1zxB5tueO2ecYsmURGkxgMx6-AOBAUY8WZ8dWw; __Secure-3PSIDCC=AJi4QfFBdEFXcQFlKqAhaj5Ev2D0su31YpK9y1sJRYAiDUkZhsAy6GJ4IQYaz9aSQQMzEDT4R7o")
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	bodyText, _ := ioutil.ReadAll(resp.Body)
	x := soup.HTMLParse(string(bodyText))
	g := x.Find("span", "id", "tw-answ-target-text")
	c.Reply(fmt.Sprint(g.Text()))
	return nil
}

func Music(c tb.Context) error {
	r, _ := SearchYT(c.Message().Payload, 2)
	fmt.Println(r.Items[0])
	ID := r.Items[0].Id.VideoId
	y := yt.Client{HTTPClient: myClient}
	vid, err := y.GetVideo("https://www.youtube.com/watch?v=" + ID)
	format := vid.Formats.FindByQuality("tiny")
	stream, _, _ := y.GetStream(vid, format)
	b, _ := ioutil.ReadAll(stream)
	ioutil.WriteFile("t.mp3", b, 0666)
	stream.Close()
	check(err)
	duration, _ := time.ParseDuration(vid.Duration.String())
	c.Bot().Notify(c.Chat(), "upload_voice")
	return c.Reply(&tb.Audio{
		File:      tb.File{FileLocal: "t.mp3"},
		Title:     vid.Title,
		Performer: vid.Author,
		FileName:  vid.Title,
		Duration:  int(duration.Seconds()),
	})
}
