package modules

import (
	"go.mongodb.org/mongo-driver/bson"
)

const BOT_ID = int64(5181620060)
const OWNER_ID = int64(1833850637)
const YOUTUBE_API_KEY = "AIzaSyAEz0eRkbsEE7TrHGKmd_iXh4AmYJlMKDs"

var notes_help = "✨ Here is the help for **Notes:**\n**Command for Members**\n**->** `/get notename`: get the note with this notename\n**-** #notename: same as /get\n**->** `/notes`: list all saved notes in this chat\n**Command for Admins**\n**->** `/save notename notedata`: saves notedata as a note with name notename, reply to a message or document to save it\n**->** `/clear notename`: clear note with this name\n**->** `/privatenote on/yes/off/no`: whether or not to send the note in PM. Write del besides on/off to delete hashtag message on group.\n**Note**\n **-** Only admins can use This module\n **-** To save a document (like photo, audio, etc.), reply to a document or media then type /save\n **-** Need help for parsing text? Check /markdownhelp\nSave data for future users with notes!\nNotes are great to save random tidbits of information; a phone number, a nice gif, a funny picture - anything!\nAlso you can save a text/document with buttons, you can even save it in here."

var help = bson.M{"notes": notes_help}

var PLUGIN_LIST = []string{"admin", "bans", "chatbot", "feds", "greetings", "inline", "lock", "misc", "notes", "pin", "stickers", "warns"}

var help_caption = `
Hey!, My name is Mika.
I am a group management bot, here to help you get around and keep the order in your groups!
I have lots of handy features.
So what are you waiting for?
Add me in your groups and give me full rights to make me function well.`

var COUNTRY_CODES = bson.M{"Argentina": "ar", "Australia": "au", "Bangladesh": "bd", "Belgium": "be", "Brazil": "br", "Canada": "ca", "China": "cn", "Czech Republic": "cz", "France": "fr", "Germany": "de", "Greece": "gr", "Hungary": "hu", "India": "in", "Indonesia": "id", "Iran": "ir", "Italy": "it", "Japan": "jp", "Malaysia": "my", "Mexico": "mx", "Netherlands": "nl", "Nigeria": "ng", "Peru": "pe", "Philippines": "ph", "Poland": "pl", "Portugal": "pt", "Romania": "ro", "Russia": "ru", "Saudi Arabia": "sa", "Singapore": "sg", "South Africa": "za", "South Korea": "kr", "Spain": "es", "Sweden": "se", "Thailand": "th", "Turkey": "tr", "Uganda": "ug", "Ukraine": "ua", "United Kingdom": "uk", "United States": "us", "Vietnam": "vn"}
var CODE_C = []string{"ar", "au", "bd", "be", "br", "ca", "cn", "cz", "fr", "de", "gr", "hu", "in", "id", "ir", "it", "jp", "my", "mx", "nl", "ng", "pe", "ph", "pl", "pt", "ro", "ru", "sa", "sg", "za", "kr", "es", "se", "th", "tr", "ug", "ua", "uk", "us", "vn"}

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

var dead_cc = `
<b>❌ Stripe/15$ - DEAD</b>
<b>CC:</b> <code>%s|%s|%s|%s</code>
<u>Result:</u> <b>%s</b>
%s
<U>Time:</u> <b>%d's</b>
<b>Checked By <a href='tg://user?id=%d'>%s</a> [%s]</b>`

var ccn_cc = `
<b>✅ Stripe/15$ - CCN Live</b>
<b>CC:</b> <code>%s|%s|%s|%s</code>
<u>Result:</u> <b>%s</b>
%s
<U>Time:</u> <b>%d's</b>
<b>Checked By <a href='tg://user?id=%d'>%s</a> [%s]</b>`

var insuf_funds = `
<b>✅ Stripe/15$ - Insuffient Funds</b>
<b>CC:</b> <code>%s|%s|%s|%s</code>
<u>Result:</u> <b>%s</b>
%s
<U>Time:</u> <b>%d's</b>
<b>Checked By <a href='tg://user?id=%d'>%s</a> [%s]</b>`

var stripe_rs = `
>. <b>𝐆𝐚𝐭𝐞𝐬/𝐂𝐡𝐚𝐫𝐠𝐞/𝐒𝐭𝐫𝐢𝐩𝐞-$1</b>

| —  <b>𝐑𝐄𝐒𝐔𝐋𝐓</b>
|- <b>CARD:</b> <code>%s|%s|%s|%s</code>
|- <b>STATUS:</b> %s %s
|- <b>RESPONSE:</b> %s
|- <b>CODE:</b> %s
| —  <b>𝐁𝐈𝐍-𝐈𝐍𝐅𝐎</b>
|- <b>BANK/TYPE:</b> .
|- <b>COUNTRY:</b> . .
| —  <b>𝐈𝐍𝐅𝐎𝐒</b>
|- <b>Checked By:</b> @%s[%s]
|- <b>Time Taken:</b>  "."

| — 𝙱𝚘𝚝 𝙱𝚢 <b>RoseLoverX</b>
`

var stripe_1 = `
<b>CC ✑</b> <code>%s|%s|%s|%s</code>
<b>⌥ Status ✑ %s</b> %s %s
<b>⌥ Response ✑</b> <i>%s</i>
<b>⌥ Gateway ✑ Stripe 1$</b>

<b>⎋ Card Details: %s</b>

<b>✁Checked by %s [%s]</b>
<b>✁Bot by 𝙉ᴇρтυиᴇ 🇬🇷[𝙸𝚅𝙰𝚁]</b>
`
