package modules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/amarnathcjd/yoko/bot"
	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

type HANDLE struct {
	FUNC       func(tb.Context) error
	MIDDLEWARE func(tb.HandlerFunc) tb.HandlerFunc
}

var HANDLERS = GatherHandlers()

func GatherHandlers() map[string]HANDLE {
	var HANDLERS = make(map[string]HANDLE)
	// misc.go
	HANDLERS["info"] = HANDLE{FUNC: UserInfo}
	HANDLERS["imdb"] = HANDLE{FUNC: Imdb}
	HANDLERS["translate"] = HANDLE{FUNC: Translate}
	HANDLERS["ud"] = HANDLE{FUNC: UDict}
	HANDLERS["bin"] = HANDLE{FUNC: BinCheck}
	HANDLERS["telegraph"] = HANDLE{FUNC: telegraph}
	HANDLERS["math"] = HANDLE{FUNC: Math}
	HANDLERS["id"] = HANDLE{FUNC: GetID}
	HANDLERS["paste"] = HANDLE{FUNC: Paste}
	HANDLERS["stat"] = HANDLE{FUNC: GroupStat}
	HANDLERS["tr"] = HANDLE{FUNC: Translate}
	HANDLERS["chatinfo"] = HANDLE{FUNC: ChatInfo}
	HANDLERS["music"] = HANDLE{FUNC: Music}
	HANDLERS["rs"] = HANDLE{FUNC: RsStripe}
	HANDLERS["wiki"] = HANDLE{FUNC: WikiPedia}
	HANDLERS["fake"] = HANDLE{FUNC: FakeGen}
	HANDLERS["insta"] = HANDLE{FUNC: InstaCSearch}
	HANDLERS["roll"] = HANDLE{FUNC: Roll}
	HANDLERS["json"] = HANDLE{FUNC: Json}
	HANDLERS["doge"] = HANDLE{FUNC: DogeSticker}
	// start.go
	HANDLERS["start"] = HANDLE{FUNC: Start}
	HANDLERS["help"] = HANDLE{FUNC: Help_Menu}
	// notes.go
	HANDLERS["save"] = HANDLE{FUNC: Save, MIDDLEWARE: Change_info}
	HANDLERS["saved"] = HANDLE{FUNC: All_notes}
	HANDLERS["notes"] = HANDLE{FUNC: All_notes}
	HANDLERS["get"] = HANDLE{FUNC: Gnote}
	HANDLERS["clear"] = HANDLE{FUNC: clear_note, MIDDLEWARE: Change_info}
	HANDLERS["clearall"] = HANDLE{FUNC: clear_all, MIDDLEWARE: Change_info}
	HANDLERS["privatenotes"] = HANDLE{FUNC: private_notes, MIDDLEWARE: Change_info}
	// locks.go
	HANDLERS["lock"] = HANDLE{FUNC: Lock, MIDDLEWARE: Change_info}
	HANDLERS["locks"] = HANDLE{FUNC: Check_locks}
	HANDLERS["locktypes"] = HANDLE{FUNC: Locktypes}
	HANDLERS["unlock"] = HANDLE{FUNC: Unlock, MIDDLEWARE: Change_info}
	// admin.go
	HANDLERS["promote"] = HANDLE{FUNC: Promote, MIDDLEWARE: Add_admins}
	HANDLERS["demote"] = HANDLE{FUNC: Demote, MIDDLEWARE: Add_admins}
	HANDLERS["superpromote"] = HANDLE{FUNC: Promote, MIDDLEWARE: Add_admins}
	HANDLERS["settitle"] = HANDLE{FUNC: Set_title, MIDDLEWARE: Add_admins}
	HANDLERS["title"] = HANDLE{FUNC: Set_title, MIDDLEWARE: Add_admins}
	HANDLERS["adminlist"] = HANDLE{FUNC: Adminlist}
	// chatbot.go
	HANDLERS["chatbot"] = HANDLE{FUNC: Chatbot_mode, MIDDLEWARE: Change_info}
	// connect.go
	HANDLERS["connect"] = HANDLE{FUNC: ConnectChat}
	// greetings.go
	HANDLERS["welcome"] = HANDLE{Welcome_set, Change_info}
	HANDLERS["setwelcome"] = HANDLE{Set_welcome, Change_info}
	HANDLERS["resetwelcome"] = HANDLE{ResetWelcome, Change_info}
	// warnings.go
	HANDLERS["warn"] = HANDLE{WARN, Ban_users}
	HANDLERS["setwarnmode"] = HANDLE{Set_warn_mode_hn, Ban_users}
	HANDLERS["warnmode"] = HANDLE{Set_warn_mode_hn, Ban_users}
	HANDLERS["setwarnlimit"] = HANDLE{Set_warn_limit, Change_info}
	HANDLERS["warnlimit"] = HANDLE{Set_warn_limit, Change_info}
	HANDLERS["warnings"] = HANDLE{FUNC: Warnings_info}
	// eval.go
	HANDLERS["sh"] = HANDLE{FUNC: Exec}
	HANDLERS["media"] = HANDLE{FUNC: MediaInfo}
	// stickers.go
	HANDLERS["kang"] = HANDLE{FUNC: AddSticker}
	HANDLERS["packs"] = HANDLE{FUNC: MyPacks}
	HANDLERS["mypacks"] = HANDLE{FUNC: MyPacks}
	HANDLERS["stickers"] = HANDLE{FUNC: CombotSticker}
	// pin.go
	HANDLERS["pin"] = HANDLE{pin_message, Pin_messages}
	HANDLERS["unpin"] = HANDLE{unpin_msg, Pin_messages}
	HANDLERS["pinned"] = HANDLE{FUNC: pinned_msg}
	HANDLERS["permapin"] = HANDLE{PermaPin, Pin_messages}
	// ban.go
	HANDLERS["ban"] = HANDLE{Ban, Ban_users}
	HANDLERS["dban"] = HANDLE{Ban, Ban_users}
	HANDLERS["sban"] = HANDLE{Ban, Ban_users}
	HANDLERS["tban"] = HANDLE{Ban, Ban_users}
	HANDLERS["unban"] = HANDLE{Ban, Ban_users}
	HANDLERS["mute"] = HANDLE{Mute, Ban_users}
	HANDLERS["tmute"] = HANDLE{Mute, Ban_users}
	HANDLERS["dmute"] = HANDLE{Mute, Ban_users}
	HANDLERS["smute"] = HANDLE{Mute, Ban_users}
	HANDLERS["unmute"] = HANDLE{Mute, Ban_users}
	HANDLERS["kick"] = HANDLE{Kick, Ban_users}
	HANDLERS["skick"] = HANDLE{Kick, Ban_users}
	HANDLERS["dkick"] = HANDLE{Kick, Ban_users}
	HANDLERS["kickme"] = HANDLE{FUNC: KickMe}
	// feds.go
	HANDLERS["newfed"] = HANDLE{FUNC: New_fed}
	HANDLERS["delfed"] = HANDLE{FUNC: Delete_fed}
	HANDLERS["renamefed"] = HANDLE{FUNC: Rename_fed}
	HANDLERS["joinfed"] = HANDLE{FUNC: Join_fed}
	HANDLERS["leavefed"] = HANDLE{FUNC: Leave_fed}
	HANDLERS["fpromote"] = HANDLE{FUNC: Fpromote}
	HANDLERS["fedpromote"] = HANDLE{FUNC: Fpromote}
	HANDLERS["fdemote"] = HANDLE{FUNC: Fdemote}
	HANDLERS["feddemote"] = HANDLE{FUNC: Fdemote}
	HANDLERS["ftransfer"] = HANDLE{FUNC: Transfer_fed_user}
	HANDLERS["fedtransfer"] = HANDLE{FUNC: Transfer_fed_user}
	HANDLERS["fedexport"] = HANDLE{FUNC: FedExport}
	HANDLERS["exportfed"] = HANDLE{FUNC: FedExport}
	HANDLERS["fexport"] = HANDLE{FUNC: FedExport}
	HANDLERS["fban"] = HANDLE{FUNC: Fban}
	HANDLERS["fedban"] = HANDLE{FUNC: Fban}
	HANDLERS["unfban"] = HANDLE{FUNC: Unfban}
	HANDLERS["unfedban"] = HANDLE{FUNC: Unfban}
	HANDLERS["fedunban"] = HANDLE{FUNC: Unfban}
	HANDLERS["fednotif"] = HANDLE{FUNC: Fednotif}
	HANDLERS["subfed"] = HANDLE{FUNC: sub_fed}
	HANDLERS["unsubfed"] = HANDLE{FUNC: unsub_fed}
	HANDLERS["fedinfo"] = HANDLE{FUNC: fed_info}
	HANDLERS["chatfed"] = HANDLE{FUNC: Chat_fed}
	// filters.go
	HANDLERS["filter"] = HANDLE{FUNC: SaveFilter, MIDDLEWARE: Change_info}
	HANDLERS["filters"] = HANDLE{FUNC: AllFilters}
	HANDLERS["stop"] = HANDLE{StopFilter, Change_info}
	HANDLERS["stopall"] = HANDLE{StopAllFIlters, Change_info}
	// antiflood.go
	HANDLERS["flood"] = HANDLE{FUNC: Flood}
	HANDLERS["setflood"] = HANDLE{SetFlood, Ban_users}
	HANDLERS["setfloodmode"] = HANDLE{SetFloodMode, Ban_users}
	// gbans.go
	HANDLERS["gban"] = HANDLE{FUNC: Gban}
	// devs.go
	HANDLERS["addsudo"] = HANDLE{FUNC: AddSudo}
	HANDLERS["adddev"] = HANDLE{FUNC: AddDev}
	HANDLERS["devs"] = HANDLE{FUNC: ListDev}
	HANDLERS["sudolist"] = HANDLE{FUNC: ListSudo}
	HANDLERS["remdev"] = HANDLE{FUNC: RemoveDev}
	HANDLERS["remsudo"] = HANDLE{FUNC: RemoveSudo}
	HANDLERS["logs"] = HANDLE{FUNC: Logs}
	HANDLERS["ping"] = HANDLE{FUNC: Ping}
	HANDLERS["stats"] = HANDLE{FUNC: Stats}
	HANDLERS["pong"] = HANDLE{FUNC: Te}
	HANDLERS["sendmessage"] = HANDLE{FUNC: SendMessage}
	// rules.go
	HANDLERS["rules"] = HANDLE{FUNC: Rules}
	HANDLERS["setrules"] = HANDLE{FUNC: SetRules, MIDDLEWARE: Change_info}
	HANDLERS["resetrules"] = HANDLE{FUNC: ResetRules, MIDDLEWARE: Change_info}
	HANDLERS["setrulesbutton"] = HANDLE{FUNC: SetRulesButton, MIDDLEWARE: Change_info}
	HANDLERS["resetsetrulesbutton"] = HANDLE{FUNC: ResetRulesButton, MIDDLEWARE: Change_info}
	HANDLERS["privaterules"] = HANDLE{FUNC: PrivateRules, MIDDLEWARE: Change_info}
	// purge.go
	HANDLERS["purge"] = HANDLE{FUNC: Purge, MIDDLEWARE: DeleteMessages}
	HANDLERS["del"] = HANDLE{FUNC: Delete, MIDDLEWARE: DeleteMessages}
	HANDLERS["purgefrom"] = HANDLE{FUNC: PurgeFrom, MIDDLEWARE: DeleteMessages}
	HANDLERS["purgeto"] = HANDLE{FUNC: PurgeTo, MIDDLEWARE: DeleteMessages}
	HANDLERS["pinterest"] = HANDLE{FUNC: PinterestSearch}
	return HANDLERS

}

func RegisterHandlers() {
	for endpoint, function := range HANDLERS {
		if function.MIDDLEWARE != nil {
			bot.Bot.Handle("/"+endpoint, function.FUNC, function.MIDDLEWARE, ConnectFunc)
		} else {
			bot.Bot.Handle(fmt.Sprintf("/%s", endpoint), function.FUNC, ConnectFunc)
		}
	}
	CallBackHandlers()
}

func CallBackHandlers() {
	bot.Bot.Handle(&help_button, HelpCB)
	bot.Bot.Handle(&back_button, back_cb)
	// notes.go
	bot.Bot.Handle(&dall_btn, del_all_notes_cb)
	bot.Bot.Handle(&cdall_btn, cancel_delall_cb)
	// feds.go
	bot.Bot.Handle(&accept_fpromote, Fpromote_cb)
	bot.Bot.Handle(&deny_fpromote, Fpromote_deny_cb)
	bot.Bot.Handle(&accept_ftransfer, Accept_Transfer_fed_cb)
	bot.Bot.Handle(&deny_ftransfer, Decline_Transfer_fed_cb)
	bot.Bot.Handle(&confirm_ftransfer, Confirm_Transfer_Fed_cb)
	bot.Bot.Handle(&reject_ftransfer, Deny_Transfer_Fed_cb)
	bot.Bot.Handle(&check_fed_admins, Check_f_admins_cb)
	// common handlers
	bot.Bot.Handle(tb.OnQuery, InlineQueryHandler)
	bot.Bot.Handle(tb.OnText, OnTextHandler)
	bot.Bot.Handle(tb.OnChatMember, OnChatMemberHandler)
	// warns.go
	bot.Bot.Handle(&unwarn_btn, UnWarnCb)
	// filters.go
	bot.Bot.Handle(&del_all_filters, DelAllFCB)
	bot.Bot.Handle(&cancel_del_all_filters, CancelDALL)
	bot.Bot.Handle(&imdb_btn, ImdbCB)
	bot.Bot.Handle(&ChatBTN, ChatbotCB)
	bot.Bot.Handle(&anon_button, AnonCB)
	bot.Bot.Handle(tb.OnAddedToGroup, AddedToGroupHandler)
}

func OnTextHandler(c tb.Context) error {
	if c.Sender().Username != string("") && c.Sender().Username == "nglnah" {
		c.Bot().Send(&tb.User{ID: OWNER_ID}, "She's online.")
	}
	if strings.HasPrefix(c.Message().Text, "!") || strings.HasPrefix(c.Message().Text, "?") {
		cmd := strings.Split(c.Message().Text, " ")[0][1:]
		for endpoint, function := range HANDLERS {
			if endpoint == cmd {
				c = AddPayload(c)
				if c.Message().Private() {
					c = ChatContext(c)
				}
				if function.MIDDLEWARE == nil {
					return function.FUNC(c)
				} else {
					if proc := function.MIDDLEWARE(function.FUNC); proc != nil {
						return proc(c)
					}
				}
			}
		}
	}
	FLOOD_EV(c)
	Chat_bot(c)
	if ok, err := FilterEvent(c); ok {
		return err
	}
	if match, _ := regexp.MatchString("\\#(\\S+)", c.Message().Text); match {
		Hash_note(c)
		return nil
	}
	if afk := AFK(c); afk {
		return nil
	}
	return nil
}

func AddedToGroupHandler(c tb.Context) error {
	if !db.IsChat(c.Chat().ID) {
		db.AddChat(db.Chat{Id: c.Chat().ID, Title: c.Chat().Title})
	}
	return nil
}
