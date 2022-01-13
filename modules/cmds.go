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
	bot.Bot.Handle("/id", ID_info)
	bot.Bot.Handle("/paste", Paste)
	bot.Bot.Handle("/fake", Fake_gen)
	// start.go
	bot.Bot.Handle("/start", Start)
	bot.Bot.Handle("/help", Help_Menu)
	bot.Bot.Handle(&help_button, HelpCB)
	bot.Bot.Handle(&back_button, back_cb)
	// notes.go
	bot.Bot.Handle("/save", Save, Change_info)
	bot.Bot.Handle("/saved", All_notes)
	bot.Bot.Handle("/notes", All_notes)
	bot.Bot.Handle("/get", Gnote)
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
	bot.Bot.Handle("/joinfed", Join_fed)
	bot.Bot.Handle("/leavefed", Leave_fed)
	bot.Bot.Handle("/chatfed", Chat_fed)
	bot.Bot.Handle("/fpromote", Fpromote)
	bot.Bot.Handle(&accept_fpromote, Fpromote_cb)
	bot.Bot.Handle(&deny_fpromote, Fpromote_deny_cb)
	bot.Bot.Handle("/fdemote", Fdemote)
	bot.Bot.Handle("/ftransfer", Transfer_fed_user)
	bot.Bot.Handle(&accept_ftransfer, Accept_Transfer_fed_cb)
	bot.Bot.Handle(&deny_ftransfer, Decline_Transfer_fed_cb)
	bot.Bot.Handle(&confirm_ftransfer, Confirm_Transfer_Fed_cb)
	bot.Bot.Handle(&reject_ftransfer, Deny_Transfer_Fed_cb)
	bot.Bot.Handle("/fban", Fban)
	bot.Bot.Handle("/unfban", Unfban)
	bot.Bot.Handle("/subfed", sub_fed)
	bot.Bot.Handle("/unsubfed", unsub_fed)
	bot.Bot.Handle("/fedinfo", fed_info)
	bot.Bot.Handle(&check_fed_admins, Check_f_admins_cb)
	bot.Bot.Handle("/fednotif", Fednotif)
	// bans.go
	bot.Bot.Handle("/ban", Ban, Ban_users)
	bot.Bot.Handle("/tban", Ban, Ban_users)
	bot.Bot.Handle("/sban", Ban, Ban_users)
	bot.Bot.Handle("/dban", Ban, Ban_users)
	bot.Bot.Handle("/unban", Ban, Ban_users)
	bot.Bot.Handle("/mute", Mute, Ban_users)
	bot.Bot.Handle("/dmute", Mute, Ban_users)
	bot.Bot.Handle("/smute", Mute, Ban_users)
	bot.Bot.Handle("/tmute", Mute, Ban_users)
	bot.Bot.Handle("/unmute", Mute, Ban_users)
	bot.Bot.Handle("/kick", Kick, Ban_users)
	// common handlers
	bot.Bot.Handle(tb.OnQuery, InlineQueryHandler)
	bot.Bot.Handle(tb.OnText, OnTextHandler)
	// pin.go
	bot.Bot.Handle("/pin", pin_message, Pin_messages)
	bot.Bot.Handle("/pinned", pinned_msg)
	// stickers.go
	bot.Bot.Handle("/kang", AddSticker)
	bot.Bot.Handle("/d", TGStest)
	// eval.go
	bot.Bot.Handle("/sh", Exec)
	bot.Bot.Handle("/eval", Eval)
	// warns.go
	bot.Bot.Handle("/warn", WARN, Ban_users)
	bot.Bot.Handle("/setwarnmode", Set_warn_mode_hn, Ban_users)
	bot.Bot.Handle("/warnmode", Set_warn_mode_hn, Ban_users)
	bot.Bot.Handle("/warnings", Warnings_info)
	// greetings.go
	bot.Bot.Handle("/welcome", Welcome_set)
	// chatbot.go
	bot.Bot.Handle("/chatbot", Chatbot_mode, Change_info)
}
