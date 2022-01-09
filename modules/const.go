package modules

import (
	"go.mongodb.org/mongo-driver/bson"
)

var notes_help = `✨ Here is the help for <b>Notes:<b>
<b>Command for Members</b>
<b>-></b> <code>/get notename</code>: get the note with this notename
<b>-</b> <code>#notename</code>: same as /get
<b>-></b> <code>/notes<code>: list all saved notes in this chat
<b>Command for Admins</b>
<b>-></b> <code>/save notename notedata</code>: saves notedata as a note with name notename, reply to a message or document to save it
<b>-></b> <code>/clear notename</code>: clear note with this name
<b>-></b> <code>/privatenote on/yes/off/no</code>: whether or not to send the note in PM. Write del besides on/off to delete hashtag message on group.
<b>Note</b>
<b>-</b> Only admins can use This module
<b>-</b> To save a document (like photo, audio, etc.), reply to a document or media then type /save
<b>-</b> Need help for parsing text? Check /markdownhelp
Save data for future users with notes!
Notes are great to save random tidbits of information; a phone number, a nice gif, a funny picture - anything!
Also you can save a text/document with buttons, you can even save it in here.`

var help = bson.M{"notes": notes_help}

var PLUGIN_LIST = []string{"admin", "bans", "chatbot", "feds", "greetings", "inline", "lock", "misc", "notes", "pin", "stickers", "warns"}

var help_caption = `
Hey!, My name is Yoko.
I am a group management bot, here to help you get around and keep the order in your groups!
I have lots of handy features.
So what are you waiting for?
Add me in your groups and give me full rights to make me function well.`

var COUNTRY_CODES = bson.M{"Argentina": "ar", "Australia": "au", "Bangladesh": "bd", "Belgium": "be", "Brazil": "br", "Canada": "ca", "China": "cn", "Czech Republic": "cz", "France": "fr", "Germany": "de", "Greece": "gr", "Hungary": "hu", "India": "in", "Indonesia": "id", "Iran": "ir", "Italy": "it", "Japan": "jp", "Malaysia": "my", "Mexico": "mx", "Netherlands": "nl", "Nigeria": "ng", "Peru": "pe", "Philippines": "ph", "Poland": "pl", "Portugal": "pt", "Romania": "ro", "Russia": "ru", "Saudi Arabia": "sa", "Singapore": "sg", "South Africa": "za", "South Korea": "kr", "Spain": "es", "Sweden": "se", "Thailand": "th", "Turkey": "tr", "Uganda": "ug", "Ukraine": "ua", "United Kingdom": "uk", "United States": "us", "Vietnam": "vn"}
var CODE_C = []string{"ar", "au", "bd", "be", "br", "ca", "cn", "cz", "fr", "de", "gr", "hu", "in", "id", "ir", "it", "jp", "my", "mx", "nl", "ng", "pe", "ph", "pl", "pt", "ro", "ru", "sa", "sg", "za", "kr", "es", "se", "th", "tr", "ug", "ua", "uk", "us", "vn"}
