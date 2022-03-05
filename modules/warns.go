package modules

import (
	"fmt"
	"strconv"
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

var (
	unwarn_btn        = sel.Data("Remove warn (admin only)", "remove_user_warning")
	resetall_btn      = sel.Data("Reset all warns", "reset_all_warns")
	cancel_resetwarns = sel.Data("Cancel", "cancel_reset_warns")
)

func WarnUser(c tb.Context) error {
	cmd := strings.SplitN(c.Message().Text, " ", 1)[0][1:]
	if cmd == "dwarn" && !c.Message().IsReply() {
		c.Reply("You have to reply to a message to delete it and warn the user.")
		return nil
	}
	user, extra := GetUser(c)
	if user.ID == 0 {
		return nil
	}
	if user.ID == BOT_ID {
		return c.Reply("Do you really think I can do that to myself <b>:p</b>")
	}
	p, err := c.Bot().ChatMemberOf(c.Chat(), user.User())
	if err != nil {
		return c.Reply(err.Error())
	}
	if p.Role == tb.Administrator || p.Role == tb.Creator {
		return c.Reply("I'm not going to warn an admin!")
	}
	exceeded, limit, count := db.WarnUser(c.Chat().ID, user.ID, extra)
	if extra == string("") {
		extra = "No reason given."
	}
	if cmd == "dwarn" {
		c.Bot().Delete(c.Message().ReplyTo)
	}
	var ReplyTo *tb.Message
	if c.Message().ReplyTo != nil {
		ReplyTo = c.Message().ReplyTo
	} else {
		ReplyTo = c.Message()
	}
	if !exceeded && cmd != "swarn" {
		unwarn_btn.Data = strconv.Itoa(int(user.ID))
		sel.Inline(sel.Row(unwarn_btn))
		return c.Send(fmt.Sprintf("User %s has %d/%d warnings; be careful!\n<b>Reason</b>: %s", user.Mention, count, limit, extra), &tb.SendOptions{ReplyTo: ReplyTo, ReplyMarkup: sel})
	} else {
		Warns := db.GetWarns(c.Chat().ID, user.ID)
		var reasons = ""
		if len(Warns.Reasons) > 0 {
			reasons += "\n<b>Reasons:</b>"
			for i, reason := range Warns.Reasons {
				reasons += fmt.Sprintf("\n%d. ", i+1) + reason
			}
		}
		c.Send(fmt.Sprintf("That's %d/%d warnings; %s is %s%s", count, limit, user.Mention, extra, reasons), &tb.SendOptions{ReplyTo: ReplyTo})
		Settings := db.GetSettings(c.Chat().ID)
		return TakeAction(Settings.Mode, Settings.Time, *user.User(), *c.Chat())
	}
}

func SetWarnMode(c tb.Context) error {
	if c.Message().Private() {
		return c.Reply("This command is made to be used in group chats!")
	}
	arg, ctime := c.Message().Payload, int64(0)
	args := strings.SplitN(arg, " ", 2)
	if arg == string("") && strings.SplitN(c.Message().Text, " ", 2)[0][1:] == "warnmode" {
		Settings := db.GetSettings(c.Chat().ID)
		return c.Reply(fmt.Sprintf(`Users who go over the warning limit are currently: <code>%s</code>

To change the warn mode, use this command again, with one of ban/kick/mute/tban/tmute.
eg: <code>/warnmode ban</code>`, Convert_action(Settings.Mode, Settings.Time)))
	} else if arg == string("") {
		return c.Reply("You need to specify an action to take upon too many warns.")
	} else if stringInSlice(args[0], []string{"ban", "mute", "kick", "tban", "tmute"}) {
		if strings.HasPrefix(args[0], "t") {
			if len(args) < 2 {
				return c.Reply("It looks like you tried to set time value for warns but you didn't specified time; Try, <code>/setwarnmode [tban/tmute] <timevalue></code>.\n<b>Examples of time value:</b> <code>4m = 4 minutes</code>, <code>3h = 3 hours</code>, <code>6d = 6 days</code>, <code>5w = 5 weeks</code>.")
			}
			ctime = ExtractTime(c, args[1])
			if ctime == 0 {
				return nil
			}
		}
	} else {
		return c.Reply(fmt.Sprintf("Unknown type '%s'. Please use one of: ban/kick/mute/tban/tmute", args[0]))
	}
	c.Reply(fmt.Sprintf("Updated warning mode to: %s", Convert_action(args[0], int32(ctime))))
	db.SetWarnMode(c.Chat().ID, args[0], int(ctime))
	return nil
}

func SetWarnLimit(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	arg := c.Message().Payload
	if arg == string("") && strings.SplitN(c.Message().Text, " ", 2)[0][1:] == "warnlimit" {
		Settings := db.GetWarnSettings(c.Chat().ID)
		return c.Reply(fmt.Sprintf(`The current warn limit is: <code>%d</code>
After this is exceeded, users will be %s.
		
To change the warn limit, use this command again, specifying the amount of warns.
eg: <code>/warnlimit 6</code>`, Settings.Limit, Convert_action(Settings.Mode, Settings.Time)))
	} else if arg == string("") {
		return c.Reply("Please specify how many warns a user should be allowed to receive before being acted upon.")
	}
	arg = strings.SplitN(arg, "", -1)[0]
	if !isInt(arg) {
		return c.Reply(fmt.Sprintf("<b>%s</b> is not a valid integer.", arg))

	} else if arg, b := strconv.Atoi(arg); b == nil && arg < 1 {
		return c.Reply("The warning limit has to be set to a number bigger than 0.")
	} else if arg > 50 {
		return c.Reply("The maximum warning limit is 50.")
	} else {
		c.Reply(fmt.Sprintf("Warn limit settings have been updated to <b>%d</b>.", arg))
		db.SetWarnLimit(c.Chat().ID, arg)
		return nil
	}
}

func WarningsInfo(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	Settings := db.GetWarnSettings(c.Chat().ID)
	c.Reply(fmt.Sprintf("There is a %d warning limit in %s. When that limit has been exceeded, the user will be %s.", Settings.Limit, c.Chat().Title, Convert_action(Settings.Mode, Settings.Time)))
	return nil
}

func UnWarnCb(c tb.Context) error {
	p, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if err != nil {
		return nil
	}
	if p.Role == tb.Member || p.Role == tb.Left {
		return c.Respond(&tb.CallbackResponse{Text: "You need to be an admin to do this!", ShowAlert: true})
	} else if p.Role == tb.Administrator && !p.Rights.CanRestrictMembers {
		return c.Edit("You are missing the following rights to use this command: CanRestrictMembers")
	}
	user, _ := strconv.Atoi(c.Callback().Data)
	user_id := int64(user)
	r := db.RemoveWarn(c.Chat().ID, user_id)
	if !r {
		return c.Respond(&tb.CallbackResponse{Text: "This warning has already been removed.", ShowAlert: true})
	} else {
		name := "User"
		usr, err := c.Bot().ChatByID(user_id)
		if err == nil {
			name = usr.FirstName
		}
		return c.Edit(fmt.Sprintf("Admin %s has removed %s's warning.", c.Sender().FirstName, name))
	}
}

func ResetWarns(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	User, Arg := GetUser(c)
	if Arg != string("") {
		Arg = "\n<b>Reason:</b> " + Arg
	}
	if Success := db.ResetWarns(c.Chat().ID, User.ID); Success {
		return c.Reply(fmt.Sprintf("Admin %s has reset %s's warnings.%s", c.Sender().FirstName, User.Mention, Arg))
	} else {
		return c.Reply("This user has no warnings to reset.")
	}
}

func ResetChatWarns(c tb.Context) error {
	if c.Message().Private() {
		return c.Reply("This command is made to be used in group chats!")
	}
	p, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if err != nil {
		return nil
	}
	if p.Role == tb.Member || p.Role == tb.Left {
		return c.Reply("You need to be an admin to do this!")
	} else if p.Role == tb.Administrator && !p.Rights.CanRestrictMembers {
		return c.Reply("You need to be the chat creator to use this command.")
	} else if p.Role == tb.Creator {
		sel.Inline(sel.Row(resetall_btn), sel.Row(cancel_resetwarns))
		return c.Reply(fmt.Sprintf("Are you sure you would like to reset **ALL** warnings in %s? This action cannot be undone.", c.Chat().Title), sel)
	} else {
		return c.Reply("You need to be an admin to do this!")
	}
}

func RemoveWarn(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	user, extra := GetUser(c)
	Rm := db.RemoveWarn(c.Chat().ID, user.ID)
	if Rm {
		var Reason string
		if extra != string("") {
			Reason = "\n<b>Reason:</b> " + extra
		}
		return c.Reply(fmt.Sprintf("Admin %s has removed %s's warning.%s", c.Sender().FirstName, user.Mention, Reason))
	} else {
		return c.Reply("This user has no warnings to remove.")
	}
}

func GetWarns(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	User, _ := GetUser(c)
	Warns := db.GetWarns(c.Chat().ID, User.ID)
	if Warns.Count == 0 {
		return c.Reply("User " + User.Mention + " has no warnings!")
	} else {
		var Reasons string
		for i, warn := range Warns.Reasons {
			Reasons += "\n" + fmt.Sprint(i) + ". " + warn
		}
		Settings := db.GetSettings(c.Chat().ID)
		return c.Reply(fmt.Sprintf("User %s has %d/%d warnings.\n%s", User.Mention, Warns.Count, Settings.Limit, Reasons))
	}
}
