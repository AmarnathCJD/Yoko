package modules

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// constants
const OWNER_ID = int64(5112408670)
const ResolveURL = "https://707c80624779.up.railway.app/username?username="
const TelegraphToken = "11b5d2394a13e4dc286ab2c29a7df8a2b02844c76e7cf7a7300d6e5420fd"
const AuddKey = "ad8dbc01e3f943dbc07e1a3abf9b9ce3"

// vars
var BOT_USERNAME = b.Me.Username
var BOT_NAME = b.Me.FirstName
var BOT_ID = b.Me.ID
var StartTime = time.Now()
var Client = http.Client{Timeout: time.Second * 10}

//cookies
var InstagramCookies = `; ; ig_nrcb=1; ; ; ; ; ; "; "; ;`
var PinterestCookies = `_pinterest_referrer=https://www.google.com/; csrftoken=252cb8f31a7abce1ad3bd54bdb09e212; _routing_id="8eb4877c-bd2c-42ec-86ee-cbf941dc3e7a"; sessionFunnelEventLogged=1; g_state={"i_l":0}; _auth=1; _pinterest_sess=TWc9PSZaY1dZak8xbThQbHVoK2d3REJibDcvS0x6STY4YWQ5eXUwQjRRQnRHU25MdXUveExwMFVhb05kTFV2YWJEd3U3b1YvVGRKUXIyOWhNTXBNTVEwT0RnODB1SzBBY1RBV09MNzNHUXExTUZua2s4aWhvZnY1d0VKME9GTWt3eEozK3dEYWJ1WUF3T2JpVEd2dGFSMGwwVkVuRDA0RDVrZVJBaERrS0czTmI4d01wSWZmQktIWVMzODFhTm5EV0l6QmszNjhtT0Zqdlg3b0ZQUGcxdlNxQXJGeG5qYUNxK2dqdk1sTVo0alQ3L3V3OElpV0NKVklNZit2S1J6VHpIWW1MYkNXazR6QUQyc3QvM0FnMldPSGVIZXdNd2RhVUlhNGJLMC9KSHZ2empzVGQ3eGdadnZ3MS9YWjljc3FnL3VjSms5L3BEd0VvWHZJem1YWWp3M3Z2UjAyODVuMEd4alQ3aEFMQlVTSnpvNDMweXdtVlgrZTlKZVhCQTBSbTl4VVhPMXdhU3p5bkdub25OMXVqRk84RTd4TmZ0UmdKYmpWbXlVdDVqMzFiOWo2Y0xmd3A3SzdnSG1SWnlFa3hRbzRvQnc4QVdJK0ptc08zRmxSbXM1ckVBSFFtK1dHckhvMnRQMW9aUXhRaDYvUFpCUDJhZkcvUU1MZzNWQk05b1loa29zNUY3MnJwZHVybTV5R294NmljY1RwYkdLTDlONkQweDRzM09ZM1NMNmdHMnoxRVhFakJMNEdtVHVDY3Q4M3VaZ2I0Q201U0lZcmZGUGdQcFJ3c3lMZ0I5SS9wclFnWm5SWTdoUDJxVms2Tm5ReEdnMDBXMDI4Qkd0VlVnQ3JXWmtTemVXaHZncndRTFZVUUZrUHVCdmJjV05PUEdpdU5YV2FVVE1sajUzbVZXdzJPeWZ6NVNjak1meVhXNG15K1hUTjFnQm90YkpBTEtFRHZ6YndsQXdXYkFWcG9ZTjczUHlIdXkvbUxOblFsMTlKU2J3MVhnL1lkdDJBVzdBdW0rVjB3cEluZDAwN212UUR1SUdXdTBnPT0mdjV6bHFEZ2hrd2cvTUU0VGJSYllSZ2NNcFZJPQ==; _b="AWE8wCyJ3SZDna0mrM3dRqAlBU5cT0RDGYpjd4nzFS++g5yM3ahjTivMA7KKR5J7XiE="; cm_sub=denied`

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
