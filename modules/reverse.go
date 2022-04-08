package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	if len(ar) == 1 || q == "" {
		return fmt.Errorf("no results found")
	}
	c.Reply(fmt.Sprintf("<a href='%s'>Search results for %s</a>", ar[len(ar)-1], q))
	var Images tb.Album
	if len(ar) != 1 {
		Images = append(Images, &tb.Photo{File: tb.FromURL(ar[0])})
		Images = append(Images, &tb.Photo{File: tb.FromURL(ar[1])})
	}
	_, sendErr := c.Bot().SendAlbum(c.Chat(), Images)
	return sendErr
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
	ar = append(ar, r.URL)
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

func SearchSpotify(query string) {
	spotifyURL := "https://open.spotify.com/search/results/" + query
	res, err := http.Get(spotifyURL)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	doc.Find("div").Each(func(_ int, s *goquery.Selection) {
		if s.HasClass("search-result-track") {
			s.Find("a").Each(func(_ int, s *goquery.Selection) {
				if s.HasClass("track-name-link") {
					fmt.Println(s.Text())
				}
			})
		}
	})
}

// Lyrics Finder

func LyricsFind(query string) (string, string, string) {
	var ShazamSearchUri = "https://www.shazam.com/services/search/v4/en-US/IN/web/search?term=" + url.QueryEscape(query) + "&numResults=3&offset=0&types=artists,songs&limit=3"
	resp, err := http.DefaultClient.Get(ShazamSearchUri)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var data MusicSearch
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}
	Key := data.Tracks.Hits[0].Track.Key
	Thumb := data.Tracks.Hits[0].Track.Images.Background
	Lyrics := LyricsFindByKey(Key)
	if Lyrics == "" {
		return "No Lyrics Found", "", ""
	} else {
		return Lyrics, Thumb, data.Tracks.Hits[0].Track.Subtitle
	}
}

func LyricsFindByKey(key string) string {
	var SHzmUri = "https://www.shazam.com/discovery/v5/en-US/IN/web/-/track/" + key + "?shazamapiversion=v3&video=v3"
	fmt.Println(SHzmUri)
	resp, err := http.DefaultClient.Get(SHzmUri)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var data Lyrics
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}
	var lyrics string
	for _, v := range data.Sections {
		if v.Type == "LYRICS" {
			lyrics = strings.Join(v.Text, "\n")
		}
	}
	log.Println(lyrics)
	return lyrics
}

func LyricsFinderHandle(c tb.Context) error {
	query := GetArgs(c)
	lyrics, thumb, name := LyricsFind(query)
	if lyrics == "" {
		return c.Reply("No Lyrics Found")
	}
	LYRICS := "Lyrics for <b>" + name + "</b>\n" + lyrics
	if thumb != "" {
		return c.Reply(&tb.Photo{File: tb.FromURL(thumb), Caption: LYRICS})
	} else {
		return c.Reply(LYRICS)
	}
}

func SearchDictionary(word string) OxfordDict {
	req, _ := http.NewRequest("GET", "https://od-api.oxforddictionaries.com:443/api/v2/entries/"+"en-gb"+"/"+url.QueryEscape(strings.ToLower(word)), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", "b202020b")
	req.Header.Add("app_key", "192992140ee9309602f090659e50eff8")
	resp, err := Client.Do(req)
	if err != nil {
		log.Println(err)
		return OxfordDict{}
	}
	defer resp.Body.Close()
	var data OxfordDict
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
		return OxfordDict{}
	} else {
		return data
	}
}

func DictionaryHandle(c tb.Context) error {
	query := GetArgs(c)
	data := SearchDictionary(query)
	if data.Results == nil {
		return c.Reply("No results found")
	}
	var result string
	result += "<b>Definition for <u>" + strings.Title(query) + "</u></b>\n"
	if data.Results[0].LexicalEntries != nil && data.Results[0].LexicalEntries[0].Entries != nil && data.Results[0].LexicalEntries[0].Entries[0].Senses != nil && data.Results[0].LexicalEntries[0].Entries[0].Senses[0].Definitions != nil {
		result += data.Results[0].LexicalEntries[0].Entries[0].Senses[0].Definitions[0]
	}
	if data.Results[0].LexicalEntries != nil && data.Results[0].LexicalEntries[0].Entries != nil && data.Results[0].LexicalEntries[0].Entries[0].Senses != nil && data.Results[0].LexicalEntries[0].Entries[0].Senses[0].Examples != nil {
		result += "\n\n<b>Example:</b>\n"
		result += data.Results[0].LexicalEntries[0].Entries[0].Senses[0].Examples[0].Text
	}
	result += "\n              <b><i>-Oxford Dictionary</i></b>"
	if data.Results[0].LexicalEntries != nil && data.Results[0].LexicalEntries[0].Entries != nil && data.Results[0].LexicalEntries[0].Entries[0].Pronunciations != nil {
		aud := data.Results[0].LexicalEntries[0].Entries[0].Pronunciations[0].AudioFile
		if aud != "" {
			resp, err := http.Get(aud)
			if err != nil {
				log.Println(err)
				return c.Reply(result)
			}
			defer resp.Body.Close()
			var audio []byte
			if err := json.NewDecoder(resp.Body).Decode(&audio); err != nil {
				log.Println(err)
				return c.Reply(result)
			}
			os.WriteFile("pronunciation.mp3", audio, 0644)
			return c.Reply(&tb.Audio{File: tb.FromDisk("pronunciation.mp3"), Title: "Pronunciation"})
		}
	}
	return c.Reply(result)
}
