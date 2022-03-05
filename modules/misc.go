package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fogleman/gg"

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
		if c.Sender().ID == 136817688 && c.Message().SenderChat != nil {
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
	if u.ID == 0 {
		return nil
	}
	if u.ID == 1087968824 {
		Chat := c.Chat()
		u = User{
			ID:       Chat.ID,
			First:    EscapeHTML(Chat.Title),
			Last:     "",
			Username: "@" + Chat.Username,
			Mention:  "",
			DC:       0,
			Type:     "chat",
		}
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
		Info += "\n\n<b>Gbanned:</b> No"
	}
	return c.Reply(Info)

}

func GetID(c tb.Context) error {
	var u User
	if !c.Message().IsReply() && c.Message().Payload == string("") {
		if c.Sender().ID == 136817688 {
			u = User{ID: c.Message().SenderChat.ID, First: c.Message().SenderChat.FirstName, Type: "user"}
		} else {
			return c.Reply(fmt.Sprintf("<b>User ID:</b> <code>%d</code>,\n<b>Chat ID:</b> <code>%d</code>.", c.Sender().ID, c.Chat().ID))
		}

	} else {
		u, _ = GetUser(c)
	}
	if u.ID == 0 {
		return nil
	}
	if c.Message().IsReply() && c.Message().ReplyTo.IsForwarded() {
		ID, FirstName, Type := GetForwardID(c)
		user := User{ID: ID, First: FirstName, Type: Type}
		return c.Reply(fmt.Sprintf("User %s's ID is <code>%d</code>.\nThe forwarded %s, %s, has an ID of <code>%d</code>", u.First, u.ID, strings.Title(user.Type), user.First, user.ID))
	}
	return c.Reply(fmt.Sprintf("User %s's ID is <code>%d</code>.", u.First, u.ID))
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
		if chat.StickerSet != "" {
			msg += fmt.Sprintf("\n<b>Sticker Set Name:</b> %s", chat.StickerSet)
		}
		if chat.SlowMode != 0 {
			msg += fmt.Sprintf("\n<b>Slow Mode Delay:</b> %d", chat.SlowMode)
		}
		c.Reply(msg, &tb.SendOptions{DisableWebPagePreview: true})
		return nil
	} else {
		c.Reply("Invalid chat")
		return nil
	}
}

func WikiPedia(c tb.Context) error {
	Q := GetArgs(c)
	WikiMedia := "https://en.wikipedia.org/w/api.php?format=json&action=query&prop=extracts|pageimages&exintro&explaintext&generator=search&gsrsearch=intitle:" + url.QueryEscape(Q) + "&gsrlimit=1&redirects=1"
	res, err := Client.Get(WikiMedia)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)
	if data["query"] == nil {
		return c.Reply("No results found.")
	}
	pages := data["query"].(map[string]interface{})["pages"].(map[string]interface{})
	var page map[string]interface{}
	for _, v := range pages {
		page = v.(map[string]interface{})
	}
	Wiki := fmt.Sprintf("<b><u>%s</u></b>", page["title"].(string))
	var Description string
	if len(page["extract"].(string)) >= 800 {
		Extract := page["extract"].(string)[:800]
		chunks := strings.Split(Extract, ".")
		Description = strings.ReplaceAll(Extract, chunks[len(chunks)-1], "")
	} else {
		Description = page["extract"].(string)
	}
	Wiki += "\n<i>" + Description + "</i>\n -WikiPedia"
	c.Reply(Wiki)
	return nil
}

func FakeGen(c tb.Context) error {
	Args := GetArgs(c)
	if Args == "" {
		Args = "US"
	} else {
		Args = ParseCountry(Args)
	}
	res, err := Client.Get("https://randomuser.me/api?results=1&gender=&password=upper,lower,12&exc=register,picture,id&nat=" + Args)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	var FakeD FakeID
	json.NewDecoder(res.Body).Decode(&FakeD)
	FakeString := fmt.Sprintf("<b>Fake Generated (%s)</b>\n", Args)
	Fake := FakeD.Results[0]
	if Fake.Name.Title != "" {
		FakeString += "<b>First Name:</b> <code>" + Fake.Name.Title + " " + Fake.Name.First + "</code>\n"
		FakeString += "<b>Last Name:</b> <code>" + Fake.Name.Last + "</code>\n"
	}
	FakeString += "<b>Gender:</b> <code>" + Fake.Gender + "</code>\n"
	FakeString += "<b>Street:</b> <code>" + fmt.Sprint(Fake.Location.Street.Number) + ", " + Fake.Location.Street.Name + "</code>\n"
	FakeString += "<b>City:</b> <code>" + Fake.Location.City + "</code>\n"
	FakeString += "<b>State:</b> <code>" + Fake.Location.State + "</code>\n"
	FakeString += "<b>Zip:</b> <code>" + fmt.Sprint(Fake.Location.Postcode) + "</code>\n"
	FakeString += "<b>Email:</b> <code>" + Fake.Email + "</code>\n"
	FakeString += "<b>Phone:</b> <code>" + Fake.Phone + "</code>\n"
	FakeString += "<b>Cell:</b> <code>" + Fake.Cell + "</code>\n"
	FakeString += "<b>Age:</b> <code>" + fmt.Sprint(Fake.Dob.Age) + "</code>\n"
	FakeString += "<b>Birthday:</b> <code>" + Fake.Dob.Date.String() + "</code>\n"
	FakeString += "<b>Nat:</b> <code>" + Fake.Nat + "</code>\n"
	return c.Reply(FakeString)
}

func GroupStat(c tb.Context) error {
	return c.Reply(fmt.Sprintf("<b>Total Messages in %s:</b> <code>%d</code>", c.Chat().Title, c.Message().ID))
}

func Imdb(c tb.Context) error {
	Args := GetArgs(c)
	results, err := imdb.SearchTitle(&Client, Args)
	if err != nil {
		log.Print(err)
		return c.Reply("No results found.")
	}
	if len(results) == 0 {
		return c.Reply("No results found.")
	}
	Title, err := imdb.NewTitle(&Client, results[0].ID)
	if err != nil {
		log.Print(err)
		return c.Reply("No results found.")
	}
	var Movie string
	if Title.Name != "" {
		Movie = "<b>" + Title.Name + "</b>\n"
	}
	var Genres string
	var Directors string
	var AKA string
	var Actors string
	for _, x := range Title.Genres {
		Genres += x + ", "
	}
	for _, x := range Title.Directors {
		Directors += x.FullName + ", "
	}
	q := 1
	for _, x := range Title.AKA {
		q++
		if q > 4 {
			break
		}
		AKA += x + ", "
	}
	for _, x := range Title.Actors {
		Actors += x.FullName + ", "
	}
	Genres = strings.TrimSuffix(Genres, ", ")
	Directors = strings.TrimSuffix(Directors, ", ")
	AKA = strings.TrimSuffix(AKA, ", ")
	Actors = strings.TrimSuffix(Actors, ", ")
	Movie += "<b>Year:</b> <code>" + fmt.Sprint(Title.Year) + "</code>\n"
	Movie += "<b>Rating:</b> <code>" + fmt.Sprint(Title.Rating) + "</code>\n"
	Movie += "<b>Genre:</b> " + Genres + "\n"
	Movie += "<b>Runtime:</b> <code>" + fmt.Sprint(Title.Duration) + "</code>\n"
	Movie += "<b>Actors:</b> " + Actors + "\n"
	Movie += "<b>Directors:</b> " + Directors + "\n"
	Movie += "<b>Plot:</b> <i>" + Title.Description + "</i>\n"
	Movie += "<b>AKA:</b> " + AKA + "\n"
	if Title.Poster.URL != "" {
		return c.Reply(&tb.Photo{File: tb.FromURL(Title.Poster.URL), Caption: Movie})
	}
	return c.Reply(Movie)
}

func InstaCSearch(c tb.Context) error {
	Username := GetArgs(c)
	ApiUrl := `https://www.instagram.com/` + Username + `/?__a=1`
	req, _ := http.NewRequest("GET", ApiUrl, nil)
	req.Header.Add("cookie", InstagramCookies)
	res, err := Client.Do(req)
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	var d map[string]interface{}
	json.NewDecoder(res.Body).Decode(&d)
	if _, ok := d["graphql"]; !ok {
		return c.Reply("No such username found in Instagram.")
	}
	GraphQL := d["graphql"].(map[string]interface{})["user"].(map[string]interface{})
	var U = ""
	if ID, ok := GraphQL["id"]; ok {
		U += "<b>ID:</b> <code>" + ID.(string) + "</code>\n"
	}
	if name, ok := GraphQL["full_name"]; ok {
		U += "<b>FullName:</b> " + EscapeHTML(name.(string)) + "\n"
	}
	if uname, ok := GraphQL["username"]; ok {
		U += "<b>Username:</b> <a href='https://instagram.com/" + uname.(string) + "'>" + strings.Title(uname.(string)) + "</a>\n"
	}
	if site, ok := GraphQL["external_url"]; ok && site != nil {
		U += "<b>Website:</b> <code>" + EscapeHTML(site.(string)) + "</code>\n"
	}
	Followers := GraphQL["edge_followed_by"].(map[string]interface{})["count"].(float64)
	Following := GraphQL["edge_follow"].(map[string]interface{})["count"].(float64)
	if bio, ok := GraphQL["biography"]; ok && bio.(string) != "" {
		U += "<b>Bio:</b> " + EscapeHTML(bio.(string)) + "\n"
	}
	if Vf, ok := GraphQL["is_verified"]; ok {
		U += "<b>Verified:</b> " + fmt.Sprint(Vf) + "\n"
	}
	U += "<b>Following:</b> <code>" + fmt.Sprint(int(Following)) + "</code>\n"
	U += "<b>Followers:</b> <code>" + fmt.Sprint(int(Followers)) + "</code>"
	sel.Inline(sel.Row(sel.URL(GraphQL["username"].(string), "https://instagram.com/"+GraphQL["username"].(string))))
	if pfp, ok := GraphQL["profile_pic_url_hd"]; ok {
		return c.Reply(&tb.Photo{File: tb.FromURL(pfp.(string)), Caption: U}, &tb.SendOptions{DisableWebPagePreview: true, ReplyMarkup: sel})
	}
	return c.Reply(U, &tb.SendOptions{DisableWebPagePreview: true, ReplyMarkup: sel})
}

func Roll(c tb.Context) error {
	return c.Reply(&tb.Dice{Type: "🎲", Value: rand.Intn(6)})
}

func UDict(c tb.Context) error {
	Args := GetArgs(c)
	if Args == "" {
		return c.Reply("Please specify a word.")
	}
	ApiUrl := "https://api.urbandictionary.com/v0/define?term=" + url.QueryEscape(Args)
	req, _ := http.NewRequest("GET", ApiUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	res, err := Client.Do(req)
	check(err)
	defer res.Body.Close()
	var d UrbanDict
	json.NewDecoder(res.Body).Decode(&d)
	if d.List == nil {
		return c.Reply("No such word found.")
	}
	var U = ""
	U += "<b>Word:</b> " + d.List[0].Word + "\n"
	U += "<b>Definition:</b> " + d.List[0].Definition + "\n\n"
	U += "<b>Example:</b> " + d.List[0].Example + "\n"
	U += "\n<b>Author:</b> " + d.List[0].Author
	sel.Inline(sel.Row(sel.Data("👍"+fmt.Sprint(d.List[0].ThumbsUp), "thumbs_up"), sel.Data("👎"+fmt.Sprint(d.List[0].ThumbsDown), "thumbs_down")))
	return c.Reply(U, sel)
}

func VotingCB(c tb.Context) error {
	return c.Respond(&tb.CallbackResponse{Text: "soon!", ShowAlert: true})
}

func BinCheck(c tb.Context) error {
	Args := GetArgs(c)
	if Args == "" {
		return c.Reply("Please specify a bin.")
	}
	ApiUrl := "https://lookup.binlist.net/" + Args
	req, _ := http.NewRequest("GET", ApiUrl, nil)
	resp, err := Client.Do(req)
	if err != nil {
		return c.Reply("No such bin found.")
	}
	defer resp.Body.Close()
	var d Bin
	json.NewDecoder(resp.Body).Decode(&d)
	var U = ""
	U += "<b>Bin:</b> " + Args + d.Country.Emoji + "\n"
	U += "<b>Bank:</b> " + d.Bank.Name + "\n"
	U += "<b>Country:</b> " + d.Country.Name + "(" + d.Country.Alpha2 + " - $" + d.Country.Currency + ")" + "\n"
	U += "<b>Type:</b> " + d.Type + "\n"
	U += "<b>Scheme:</b> " + d.Scheme + "\n"
	U += "<b>Brand:</b> " + d.Brand + "\n"
	U += "<b>Phone:</b> " + d.Bank.Phone + "\n"
	return c.Reply(U)
}

func Telegraph(c tb.Context) error {
	var ToUpload interface{}
	Args := GetArgs(c)
	if Args == "" {
		if c.Message().ReplyTo.Photo != nil {
			ToUpload = c.Message().ReplyTo.Photo.File
		} else if c.Message().ReplyTo.Document != nil {
			ToUpload = c.Message().ReplyTo.Document.File
		} else if c.Message().ReplyTo.Sticker != nil {
			ToUpload = c.Message().ReplyTo.Sticker.File
		} else if c.Message().ReplyTo.Video != nil {
			ToUpload = c.Message().ReplyTo.Video.File
		} else {
			return c.Reply("Please specify a file to upload.")
		}
	} else {
		ToUpload = Args
	}
	switch ToUpload.(type) {
	case string:
		req, err := http.NewRequest("GET", "https://api.telegra.ph/createPage?access_token="+TelegraphToken+"&title=Sample+Page&author_name=Anonymous&content="+ToUpload.(string)+"&return_content=true", nil)
		req.Header.Add("Content-Type", "application/json")
		resp, err := Client.Do(req)
		check(err)
		defer resp.Body.Close()
		var d mapType
		json.NewDecoder(resp.Body).Decode(&d)
		fmt.Println(d)
	case tb.File:
		return c.Reply(ToUpload.(tb.File).FileID)
	}
	return nil
}

func AuddIO(c tb.Context) error {
	if !c.Message().IsReply() {
		return c.Reply("reply to an audio message!")
	} else if c.Message().ReplyTo.Audio == nil && c.Message().ReplyTo.Video == nil {
		return c.Reply("reply to an audio/video message!")
	}
	if c.Message().ReplyTo.Audio != nil {
		c.Bot().Download(&c.Message().ReplyTo.Audio.File, "audio.mp3")
	} else if c.Message().ReplyTo.Video != nil {
		c.Bot().Download(&c.Message().ReplyTo.Video.File, "audio.mp3")
	}
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	go func() {
		defer pipeWriter.Close()
		if err := addFileToWriter(writer, "audio.mp3", "file", "audio.mp3"); err != nil {
			pipeWriter.CloseWithError(err)
		}
		writer.WriteField("api_token", "test")
		writer.WriteField("return", "apple_music,spotify")
		writer.Close()
	}()
	resp, err := Client.Post("https://api.audd.io", writer.FormDataContentType(), pipeReader)
	if err != nil {
		return c.Reply(err.Error())
	}
	defer resp.Body.Close()
	var d AuddApi
	json.NewDecoder(resp.Body).Decode(&d)
	if d.Result.Title == "" {
		return c.Reply("No results found.")
	}
	var U = ""
	U += "<b>Title:</b> " + d.Result.Title + "\n"
	U += "<b>Artist:</b> " + d.Result.Artist + "\n"
	U += "<b>Album:</b> " + d.Result.Album + "\n"
	U += "<b>Release Date:</b> " + d.Result.ReleaseDate + "\n"
	U += "<b>Label:</b> " + d.Result.Label + "\n"
	if d.Result.Spotify.ExternalUrls.Spotify != "" {
		sel.Inline(sel.Row(sel.URL("🎶 Open in Spotify", d.Result.Spotify.ExternalUrls.Spotify)))
	}
	if d.Result.Spotify.Album.Name != "" {
		return c.Reply(&tb.Photo{File: tb.FromURL(d.Result.Spotify.Album.Images[0].URL), Caption: U}, sel)
	} else if d.Result.AppleMusic.Artwork.URL != "" {
		Url := d.Result.AppleMusic.Artwork.URL
		Url = strings.Replace(Url, "{w}", fmt.Sprint(d.Result.AppleMusic.Artwork.Width), 1)
		Url = strings.Replace(Url, "{h}", fmt.Sprint(d.Result.AppleMusic.Artwork.Height), 1)
		return c.Reply(&tb.Photo{File: tb.FromURL(Url), Caption: U}, sel)
	} else {
		return c.Reply(U, sel)
	}

}

func SongDownload(c tb.Context) error {
	Args := GetArgs(c)
	search, err := YoutubeSearch(Args, 1)
	if err != nil {
		return c.Reply(err.Error())
	}
	if len(search.Items) == 0 {
		return c.Reply("No results found.")
	}
	Result := search.Items[0]
	youtube := yt.Client{}
	video, err := youtube.GetVideo("https://www.youtube.com/watch?v=" + Result.Id.VideoId)
	check(err)
	stream, _, err := youtube.GetStream(video, video.Formats.FindByItag(140))
	defer stream.Close()
	outFile, _ := os.Create("song.mp3")
	defer outFile.Close()
	io.Copy(outFile, stream)
	duration, _ := time.ParseDuration(video.Duration.String())
	c.Bot().Notify(c.Chat(), "upload_voice")
	return c.Reply(&tb.Audio{File: tb.File{FileLocal: "song.mp3"}, Title: video.Title, Duration: int(duration), FileName: video.Title, Performer: video.Author, Caption: video.Title})
}

////////////////////////////////// OLD-NEW /////////////////////////////////////////////

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

type mapType map[string]interface{}

func telesgraph(c tb.Context) error {
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
		r, err := Client.Do(req)
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

func DogeSticker(c tb.Context) error {
	Args := GetArgs(c)
	im, err := gg.LoadImage("./modules/assets/IMG_20220227_202434_649_cleanup.jpg")
	check(err)
	dc := gg.NewContext(461, 512)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("./modules/assets/Swiss 721 Black Extended BT.ttf", GetDimension(len(Args))); err != nil {
		check(err)
	}
	dc.DrawStringAnchored(Args, (461/2)-44, (512/3*3/4)-20, 0.5, 0.5)
	dc.DrawRoundedRectangle(0, 0, 461, 512, 0)
	dc.DrawImage(im, 0, 0)
	dc.DrawStringAnchored(Args, (461/2)-44, (512/3*3/4)-20, 0.5, 0.5)
	dc.Clip()
	dc.SavePNG("out.webp")
	return c.Reply(&tb.Sticker{File: tb.File{FileLocal: "out.webp"}})
}

func GetDimension(i int) float64 {
	if i <= 5 {
		return float64(85)
	} else {
		p := i + 80
		for x := 0; x <= i; x++ {
			p = p - 4
		}
		if p < 50 {
			p = 48
		}
		return float64(p)
	}
}
