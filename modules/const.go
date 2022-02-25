package modules

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// constants
const OWNER_ID = int64(1833850637)
const YOUTUBE_API_KEY = "AIzaSyAEz0eRkbsEE7TrHGKmd_iXh4AmYJlMKDs"
const ResolveURL = "https://707c80624779.up.railway.app/username?username="

// vars
var BOT_USERNAME = b.Me.Username
var BOT_NAME = b.Me.FirstName
var BOT_ID = b.Me.ID
var StartTime = time.Now()
var Client = http.Client{Timeout: time.Second * 10}

//types
type FakeID struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street struct {
				Number int    `json:"number"`
				Name   string `json:"name"`
			} `json:"street"`
			City     string `json:"city"`
			State    string `json:"state"`
			Country  string `json:"country"`
			Postcode int    `json:"postcode"`
		} `json:"location"`
		Email string `json:"email"`
		Dob   struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"dob"`
		Phone string `json:"phone"`
		Cell  string `json:"cell"`
		Nat   string `json:"nat"`
	} `json:"results"`
}

var notes_help = "✨ Here is the help for **Notes:**\n**Command for Members**\n**->** `/get notename`: get the note with this notename\n**-** #notename: same as /get\n**->** `/notes`: list all saved notes in this chat\n**Command for Admins**\n**->** `/save notename notedata`: saves notedata as a note with name notename, reply to a message or document to save it\n**->** `/clear notename`: clear note with this name\n**->** `/privatenote on/yes/off/no`: whether or not to send the note in PM. Write del besides on/off to delete hashtag message on group.\n**Note**\n **-** Only admins can use This module\n **-** To save a document (like photo, audio, etc.), reply to a document or media then type /save\n **-** Need help for parsing text? Check /markdownhelp\nSave data for future users with notes!\nNotes are great to save random tidbits of information; a phone number, a nice gif, a funny picture - anything!\nAlso you can save a text/document with buttons, you can even save it in here."

var help = bson.M{"notes": notes_help}

var PLUGIN_LIST = []string{"admin", "bans", "chatbot", "feds", "greetings", "inline", "lock", "misc", "notes", "pin", "stickers", "warns"}

var help_caption = `
Hey!, My name is Mika.
I am a group management bot, here to help you get around and keep the order in your groups!
I have lots of handy features.
So what are you waiting for?
Add me in your groups and give me full rights to make me function well.`

var COUNTRY_CODES = bson.M{"Australia": "AU", "Brazil": "BR", "Canada": "CA", "Switzerland": "CH", "Germany": "DE", "France": "FR", "Netherlands": "NL", "Russia": "RU", "Spain": "ES", "Turkey": "TR", "United Kingdom": "GB", "United States": "US", "SK": "Sweden"}
var CODE_C = []string{"AU", "BR", "CA", "CH", "DE", "DK", "ES", "FI", "FR", "GB", "IE", "IR", "NO", "NL", "NZ", "TR", "US"}
var AFK_STR = []string{
	"<b>%s</b> is here!",
	"<b>%s</b> is back!",
	"<b>%s</b> is now in the chat!",
	"<b>%s</b> is awake!",
	"<b>%s</b> is back online!",
	"<b>%s</b> is finally here!",
	"Welcome back! <b>%s</b>",
	"Where is <b>%s</b>?\nIn the chat!",
	"Pro <b>%s</b>, is back alive!",
}

var stripe_1 = `
<b>⌥ Gateway ✑ %s</b>
<b>CC ✑</b> <code>%s|%s|%s|%s</code>
<b>⌥ Status ✑ %s</b> %s %s
<b>⌥ Response ✑</b> %s

<b>⎋ Card Details: %s</b>
<b>⎋ Time: %ds</b>
<b>✁Checked by</b> <b>%s</b> [%s]
`
var CNT = []string{"locks", "flood", "filters", "get", "notes", "saved", "adminlist", "info", "warns", "rules", "approval"}
