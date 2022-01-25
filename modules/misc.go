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

	"github.com/StalkR/imdb"
	tg "github.com/TechMinerApps/telegraph"
	"github.com/anaskhan96/soup"
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
		err := c.Reply(fmt.Sprintf("User %s's ID is <code>%d</code>", user_obj.FirstName, user_obj.ID))
		fmt.Println(err)
		return nil
	}
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
	resp, err := http.Post("https://neko.roseloverx.in/api/documents", "application/json", responseBody)
	check(err)
	defer resp.Body.Close()
	var body mapType
	json.NewDecoder(resp.Body).Decode(&body)
	sel.Inline(sel.Row(sel.URL("View Paste", fmt.Sprintf("https://neko.roseloverx.in/%s", body["result"].(map[string]interface{})["key"].(string)))))
	c.Reply(fmt.Sprintf("Pasted to <b><a href='https://neko.roseloverx.in/%s'>NekoBin</a></b>.", body["result"].(map[string]interface{})["key"].(string)), &tb.SendOptions{DisableWebPagePreview: true, ReplyMarkup: sel})
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
	msg := fmt.Sprintf("<b>Fake Address Gen</b>\n<b>Full Name:</b> %s\n<b>Gender:</b> %s\n<b>DOB:</b> %s\n<b>Street:</b> %s\n<b>City/Town:</b> %s\n<b>State:</b> %s\n<b>Zip:</b> %s\n<b>Country:</b> %s\n<b>Phone Number:</b> %s", fname, gen, bday, street, city, state, zip, country, num)
	c.Reply(msg)
	return nil
}

func YT_search(c tb.Context) error {
	return nil
}

func StripeCharge(c tb.Context) error {
	d := strings.TrimSpace(c.Message().Payload)
	cc, year, month, cvc := "", "", "", ""
	for i, x := range strings.SplitN(d, "|", -1) {
		if i == 0 {
			cc = x
		} else if i == 1 {
			month = x
		} else if i == 2 {
			year = x
		} else if i == 3 {
			cvc = x
		}
	}
	for _, x := range []string{cc, year, month, cvc} {
		if x == string("") {
			c.Reply("Invalid format, please send as <code>/st cc|mm|yy|cvv</code>")
			return nil
		}
	}
	client := &http.Client{}
	current_time := time.Now()
	var postdata = strings.NewReader(fmt.Sprintf(`card[number]=%s&card[cvc]=%s&card[exp_month]=%s&card[exp_year]=%s`, cc, cvc, month, year) + `&guid=73f5a9dc-4f9d-4102-9f39-112d2bc87189f08893&muid=48d7d431-4532-4771-ad19-b3c4f6f9a71fb97291&sid=d61644dc-b7f0-456b-8669-e4aeff5cae0f629bec&payment_user_agent=stripe.js%2F558e252d7%3B+stripe-js-v3%2F558e252d7&time_on_page=1319401&key=pk_live_O98c9ngrjsN9aCgHLae6hqHU&pasted_fields=number`)
	requ, _ := http.NewRequest("POST", "https://api.stripe.com/v1/tokens", postdata)
	requ.Header.Set("authority", "api.stripe.com")
	requ.Header.Set("sec-ch-ua", `"Chromium";v="96", "Opera";v="82", ";Not A Brand";v="99"`)
	requ.Header.Set("accept", "application/json")
	requ.Header.Set("content-type", "application/x-www-form-urlencoded")
	requ.Header.Set("sec-ch-ua-mobile", "?0")
	requ.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 OPR/82.0.4227.58")
	requ.Header.Set("sec-ch-ua-platform", `"Linux"`)
	requ.Header.Set("origin", "https://js.stripe.com")
	requ.Header.Set("sec-fetch-site", "same-site")
	requ.Header.Set("sec-fetch-mode", "cors")
	requ.Header.Set("sec-fetch-dest", "empty")
	requ.Header.Set("referer", "https://js.stripe.com/")
	requ.Header.Set("accept-language", "en-US,en;q=0.9")
	res, err := client.Do(requ)
	check(err)
	defer res.Body.Close()
	var stripe bson.M
	json.NewDecoder(res.Body).Decode(&stripe)
	card_token, ok := stripe["id"].(string)
	if !ok {
		return nil
	}
	var data = strings.NewReader(fmt.Sprintf(`{"payment_type":"stripe","token":"%s","coupon":null,"save_card":true,"credit_card_id":null,"regionInfo":{"clientTcpRtt":16,"longitude":"75.77210","latitude":"11.24480","tlsCipher":"AEAD-AES128-GCM-SHA256","continent":"AS","asn":138754,"clientAcceptEncoding":"gzip, deflate, br","country":"IN","tlsClientAuth":{"certIssuerDNLegacy":"","certIssuerSKI":"","certSubjectDNRFC2253":"","certSubjectDNLegacy":"","certFingerprintSHA256":"","certNotBefore":"","certSKI":"","certSerial":"","certIssuerDN":"","certVerified":"NONE","certNotAfter":"","certSubjectDN":"","certPresented":"0","certRevoked":"0","certIssuerSerial":"","certIssuerDNRFC2253":"","certFingerprintSHA1":""},"tlsExportedAuthenticator":{"clientFinished":"8255604b0f38040d2f0559f28bbfdab14a03556da66f9de6cdc818e9ab2da9c6","clientHandshake":"a0935c8ad4a89d3f13dd1eccb660dabbee5b298c712ae658228ba641289a41e4","serverHandshake":"0be4363d7771a91e6407ced19c8662ce49510ac46dc3333302691c58a56a7ce8","serverFinished":"5e078e7242993644c6b99c1d5b643932aede4adcbb76372616e38148ec47ba3c"},"tlsVersion":"TLSv1.3","colo":"MAA","timezone":"Asia/Kolkata","city":"Kozhikode","httpProtocol":"HTTP/2","edgeRequestKeepAliveStatus":1,"requestPriority":"weight=220;exclusive=1","botManagement":{"ja3Hash":"cd08e31494f9531f560d64c695473da9","staticResource":false,"verifiedBot":false,"score":34},"clientTrustScore":34,"region":"Kerala","regionCode":"MH","asOrganization":"Kerala Vision Broad Band Private Limited","postalCode":"673004"},"email":"camarnath563@gmail.com","products":[304]}`, card_token))
	req, err := http.NewRequest("POST", "https://www.masterclass.com/api/v2/orders", data)
	check(err)
	req.Header.Set("authority", "www.masterclass.com")
	req.Header.Set("sec-ch-ua", `"Chromium";v="96", "Opera";v="82", ";Not A Brand";v="99""`)
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("x-csrf-token", "P04cdV9CTN6EiDZdBu7qou4gOQgSpVBzUfQ1jxkQxHPYDdhRR6KS37H15c/fYevAowoL+KglNj2a+c8q3BEZzg==")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 OPR/82.0.4227.58")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://www.masterclass.com")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://www.masterclass.com/plans")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "ajs_anonymous_id=%22rails-gen-95edfa57-b35a-4123-af07-a7059b0d8d81%22; first_visit=1; first_session=1; first_visit_on=2022-01-22; _session_id=80fa49d7db9128d509bd2ef301356c42; track_visit_session_id=f36cd6b3-b235-49ba-affc-be4298be1172; gdpr_checked=true; checkout_product_tiers_plan={%22name%22:%22Standard%22%2C%22price%22:{%22annual%22:%22%E2%82%B915%2C540%22%2C%22monthly%22:%22%E2%82%B91%2C295%22%2C%22flatRate%22:1554000%2C%22id%22:304}%2C%22id%22:304}; mc_cart=%5B%7B%22id%22%3A304%2C%22cart_exclusive%22%3Atrue%2C%22open_as_gift%22%3Afalse%7D%5D; cart_email=%22camarnath563%40gmail.com%22; __stripe_mid=48d7d431-4532-4771-ad19-b3c4f6f9a71fb97291; __stripe_sid=d61644dc-b7f0-456b-8669-e4aeff5cae0f629bec; _mcl=1; splitter_subject_id=8097400; split=%7B%22falcon_personalization_standalone%22%3A%22variation%22%2C%22post_pw_tv_my_progress%22%3A%22variation%22%2C%22checkout_modal_recently_viewed_instructor_avatars_disco_7_v1%22%3A%22variant_1%22%2C%22boso21_bipartisan_approach_cm%22%3A%22variant_2%22%2C%22announcement_tile_course_duration_disco_213_v1%22%3A%22variant_2%22%2C%22boso21_gift_hero_revamp%22%3A%22variant_1%22%2C%22course_impact_survey%22%3A%22variant_1%22%7D; __cf_bm=nermHfFb5_OdQdSr4IxzlHvbSuU7ywz1cPZuwtckv0E-1642827498-0-AdfMcVDHwL/WYIJ/vP+UBGEx8mmcLVctoEFIXU4sA04jL5x2QqW5qYAmdsVteuci1NJzDMvQHtwcD8KwU/atj30=")
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	var r bson.M
	json.NewDecoder(resp.Body).Decode(&r)
	bin := GetBin(cc)
	total_time := time.Now().Unix() - current_time.Unix()
	status := "Free User"
	if c.Sender().ID == int64(1833850637) {
		status = "Master"
	}
	if r["error"] != nil {
		if r["error"].(map[string]interface{})["message"].(string) == "Your card has insufficient funds." {
			c.Reply(fmt.Sprintf(insuf_funds, cc, month, year, cvc, r["error"].(map[string]interface{})["code"].(string), bin, total_time, c.Sender().ID, c.Sender().FirstName, status))
		} else if r["error"].(map[string]interface{})["code"].(string) == "card_declined" {
			c.Reply(fmt.Sprintf(dead_cc, cc, month, year, cvc, r["error"].(map[string]interface{})["message"].(string), bin, total_time, c.Sender().ID, c.Sender().FirstName, status))
		} else if r["error"].(map[string]interface{})["code"].(string) == "incorrect_cvc" {
			c.Reply(fmt.Sprintf(ccn_cc, cc, month, year, cvc, r["error"].(map[string]interface{})["message"].(string), bin, total_time, c.Sender().ID, c.Sender().FirstName, status))
		} else {
			c.Reply(fmt.Sprint(r))
		}
	} else {
		c.Reply(fmt.Sprint(r))
	}
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
		args := strings.SplitN(c.Message().Payload, " ", 2)
		if len(args) == 2 && len([]rune(args[0])) == 2 {
			lang, text = args[0], args[1]
		} else {
			text = c.Message().Payload
		}
	}
	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`async=translate,sl:auto,tl:%s,st:%s,id:1643102010421,qc:true,ac:true,_id:tw-async-translate,_pms:s,_fmt:pc`, lang, text))
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
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bodyText, _ := ioutil.ReadAll(resp.Body)
	x := soup.HTMLParse(string(bodyText))
	g := x.Find("span", "id", "tw-answ-target-text")
	c.Reply(fmt.Sprint(g.Text()))
	fmt.Println(string(bodyText))
	return nil
}
