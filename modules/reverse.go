package modules

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	tb "gopkg.in/telebot.v3"
)

type Imgdata struct {
	OU      string `json:"ou"`
	WebPage bool
}

type requestParams struct {
	Method      string
	URL         string
	Contenttype string
	Data        io.Reader
	Client      *http.Client
}

const (
	baseurl = "https://www.google.com"
)

func ReverseSearch(c tb.Context) error {
	im := Imgdata{}
	if !c.Message().IsReply() {
		return c.Reply("Reply to an image to reverse search it on Google")
	}
	var image string
	if c.Message().ReplyTo.Photo != nil {
		c.Bot().Download(&c.Message().ReplyTo.Photo.File, "rev_image.jpg")
		image = "rev_image.jpg"
	} else if c.Message().ReplyTo.Sticker != nil {
		if c.Message().ReplyTo.Sticker.Animated {
			return c.Reply("Animated stickers aren't supported")
		}
		c.Bot().Download(&c.Message().ReplyTo.Sticker.File, "rev_image.webp")
		image = "rev_image.webp"
	}
	ar, q, err := im.ImgFromFile(image)
	if err != nil {
		return err
	}
	if len(ar) == 0 || q == "" {
		return fmt.Errorf("no results found")
	}
	return c.Reply(fmt.Sprintf("Search results for %s", q))
}

func (im *Imgdata) ImgFromFile(file string) ([]string, string, error) {
	var err error
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fs, err := os.Open(file)
	if err != nil {
		return nil, "", err
	}
	defer fs.Close()
	data, err := w.CreateFormFile("encoded_image", file)
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(data, fs); err != nil {
		return nil, "", err
	}
	w.Close()
	r := &requestParams{
		Method: "POST",
		URL:    baseurl + "/searchbyimage/upload",
		Data:   &b,
		Client: &http.Client{
			Timeout:       time.Duration(10) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return fmt.Errorf("redirect") },
		},
		Contenttype: w.FormDataContentType(),
	}
	var res *http.Response
	for {
		res, err = r.fetchURL()
		if err != nil {
			return nil, "", err
		}
		if res.StatusCode == 200 {
			break
		}
		reurl, _ := res.Location()
		r.URL = reurl.String()
		r.Method = "GET"
		r.Data = nil
		r.Contenttype = ""
	}
	ar, q, err := r.getURLs(res, im.WebPage)
	if err != nil {
		return nil, q, err
	}
	return ar, q, nil
}

func (r *requestParams) fetchURL() (*http.Response, error) {
	req, err := http.NewRequest(
		r.Method,
		r.URL,
		r.Data,
	)
	if err != nil {
		return nil, err
	}
	if len(r.Contenttype) > 0 {
		req.Header.Set("Content-Type", r.Contenttype)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	res, _ := r.Client.Do(req)
	return res, nil
}

func (r *requestParams) getURLs(res *http.Response, imWebPage bool) ([]string, string, error) {
	var url string
	var match string
	var chk bool
	var ar []string
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}
	doc.Find("g-section-with-header").Each(func(_ int, s *goquery.Selection) {
		url, chk = s.Find("div").Find("title-with-lhs-icon").Find("a").Attr("href")
		if !chk {
			return
		}
	})
	r.URL = baseurl + url
	r.Client = &http.Client{Timeout: time.Duration(10) * time.Second}
	res, err = r.fetchURL()
	if err != nil {
		return nil, "", err
	}
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", err
	}
	reg1 := regexp.MustCompile("key: 'ds:1'")
	reg2 := regexp.MustCompile("\"(https?:\\/\\/.+?)\",\\d+,\\d+")
	reg3 := regexp.MustCompile(`https:\/\/encrypted\-tbn0`)
	reg4 := regexp.MustCompile("\"(https?:\\/\\/.+?)\"")
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		if reg1.MatchString(s.Text()) {
			var urls [][]string
			if imWebPage {
				strInURL := reg2.ReplaceAllString(s.Text(), "")
				urls = reg4.FindAllStringSubmatch(strInURL, -1)
			} else {
				urls = reg2.FindAllStringSubmatch(s.Text(), -1)
			}
			for _, u := range urls {
				if !reg3.MatchString(u[1]) {
					ss, err := strconv.Unquote(`"` + u[1] + `"`)
					if err == nil {
						ar = append(ar, ss)
					}
				}
			}
		}
	})
	doc.Find("input").Each(func(_ int, s *goquery.Selection) {
		if s.HasClass("og3lId") {
			match = s.AttrOr("value", "NO MATCH")
		}
	})
	if len(ar) == 0 {
		return nil, "", fmt.Errorf("no image found")
	}
	return ar, match, nil
}
