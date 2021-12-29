package modules

import (
	bot "github.com/amarnathcjd/yoko/bot"
	tb "gopkg.in/tucnak/telebot.v3"
)

func RegHandlers() {
	// misc.go
	bot.Bot.Handle("/info", Info)
	bot.Bot.Handle("/imdb", IMDb)
	bot.Bot.Handle("/crypto", Crypto)
	bot.Bot.Handle("/tr", Translate)
	bot.Bot.Handle("/ud", Ud)
	bot.Bot.Handle("/bin", Bin_check)
	bot.Bot.Handle("/telegraph", telegraph)
	bot.Bot.Handle("/math", Math)
	// start.go
	bot.Bot.Handle("/start", Start)
	// notes.go
	bot.Bot.Handle("/save", Save, Change_info)
	bot.Bot.Handle("/saved", All_notes)
	bot.Bot.Handle("/notes", All_notes)
	bot.Bot.Handle("/get", Gnote)
	bot.Bot.Handle(tb.OnText, Hash_note, Hash_regex)
	bot.Bot.Handle("/clear", clear_note, Change_info)
	bot.Bot.Handle("/clearall", clear_all)
	bot.Bot.Handle("/privatenotes", private_notes, Change_info)
	bot.Bot.Handle(&dall_btn, del_all_notes_cb)
	bot.Bot.Handle(&cdall_btn, cancel_delall_cb)
	// locks.go
	bot.Bot.Handle("/lock", Lock, Change_info)
	bot.Bot.Handle("/locktypes", Locktypes)
	bot.Bot.Handle("/locks", Check_locks)
	bot.Bot.Handle("/unlock", Unlock, Change_info)
	// admin.go
	bot.Bot.Handle("/promote", Promote, Add_admins)
	bot.Bot.Handle("/superpromote", Promote, Add_admins)
	bot.Bot.Handle("/demote", Demote, Add_admins)
	bot.Bot.Handle("/adminlist", Adminlist)
	bot.Bot.Handle("/settitle", Set_title)
	bot.Bot.Handle("/title", Set_title)
	// feds.go
	bot.Bot.Handle("/newfed", New_fed)
	bot.Bot.Handle("/delfed", Delete_fed)
	bot.Bot.Handle("/renamefed", Rename_fed)
	// bans.go
	bot.Bot.Handle("/ban", Ban)
	bot.Bot.Handle("/tban", Ban)
	bot.Bot.Handle("/sban", Ban)
	bot.Bot.Handle("/dban", Ban)
	// inline.go
	bot.Bot.Handle(tb.OnQuery, InlineQueryHandler)
	// pin.go
	bot.Bot.Handle("/pin", pin_message, Pin_messages)
	// stickers.go
	bot.Bot.Handle("/kang", AddSticker)
	bot.Bot.Handle("/d", TGStest)
	// eval.go
	bot.Bot.Handle("/sh", Exec)
	bot.Bot.Handle("/eval", Eval)
	// warns.go
	bot.Bot.Handle("/warn", WARN)
}
