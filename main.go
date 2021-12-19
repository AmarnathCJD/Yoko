package main

import (
	mod "github.com/amarnathcjd/yoko/modules"

	bot "github.com/amarnathcjd/yoko/bot"
	tb "gopkg.in/tucnak/telebot.v3"
)

func main() {
	bot.Bot.Handle("/info", mod.Info)
	bot.Bot.Handle("/imdb", mod.IMDb)
	bot.Bot.Handle("/crypto", mod.Crypto)
	bot.Bot.Handle("/start", mod.Start)
	bot.Bot.Handle("/save", mod.Save, mod.Change_info)
	bot.Bot.Handle("/lock", mod.Lock, mod.Change_info)
	bot.Bot.Handle("/locktypes", mod.Locktypes)
	bot.Bot.Handle("/locks", mod.Check_locks)
	bot.Bot.Handle("/unlock", mod.Unlock, mod.Change_info)
	bot.Bot.Handle("/tr", mod.Translate)
	bot.Bot.Handle("/saved", mod.All_notes)
	bot.Bot.Handle("/notes", mod.All_notes)
	bot.Bot.Handle("/get", mod.Gnote)
	bot.Bot.Handle("/promote", mod.Promote, mod.Add_admins)
	bot.Bot.Handle("/superpromote", mod.Promote, mod.Add_admins)
	bot.Bot.Handle("/demote", mod.Demote, mod.Add_admins)
	bot.Bot.Handle("/adminlist", mod.Adminlist)
	bot.Bot.Handle("/settitle", mod.Set_title)
	bot.Bot.Handle("/title", mod.Set_title)
	bot.Bot.Handle(tb.OnText, mod.Hash_note, mod.Hash_regex)
	bot.Bot.Handle("/ud", mod.Ud)
	bot.Bot.Handle("/newfed", mod.New_fed)
	bot.Bot.Handle("/ban", mod.Ban)
	bot.Bot.Handle("/tban", mod.Ban)
	bot.Bot.Handle("/sban", mod.Ban)
	bot.Bot.Handle("/dban", mod.Ban)
	bot.Bot.Start()
}
