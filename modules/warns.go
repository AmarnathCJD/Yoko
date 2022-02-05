package modules

import (
	"fmt"
	"strconv"
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

var unwarn_btn = sel.Data("Remove warn (admin only)", "remove_user_warning")

func WARN(c tb.Context) error {
	cmd := strings.SplitN(c.Message().Text, " ", 2)[0][1:]
	if cmd == "dwarn" && !c.Message().IsReply() {
		c.Reply("You have to reply to a message to delete it and warn the user.")
		return nil
	}
	user, extra := get_user(c.Message())
	if user.ID == int64(6) {
		c.Reply("Do you really think I can do that to myself <b>:p</b>")
		return nil
	}
	p, err := c.Bot().ChatMemberOf(c.Chat(), user)
	if err != nil {
		c.Reply(err.Error())
		return nil
	}
	if stringInSlice(string(p.Role), []string{"administrator", "creator"}) {
		c.Reply("✨ I'm not going to warn an admin!")
		return nil
	}
	exceeded, limit, count := db.Warn_user(c.Chat().ID, user.ID, extra)
	if extra == string("") {
		extra = "No reason given."
	}
	if !exceeded && cmd != "swarn" {
		unwarn_btn.Data = strconv.Itoa(int(user.ID))
		sel.Inline(sel.Row(unwarn_btn))
		c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> has %d/%d warnings; be careful!\n<b>Reason</b>: %s", user.ID, user.FirstName, count, limit, extra), sel)
		return nil
	}
	return nil
}

func Set_warn_mode_hn(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	arg, ctime := c.Message().Payload, int64(0)
	args := strings.SplitN(arg, " ", 2)
	if arg == string("") && strings.SplitN(c.Message().Text, " ", 2)[0][1:] == "warnmode" {
		_, mode, time := db.Get_warn_settings(c.Chat().ID)
		c.Reply(fmt.Sprintf(`Users who go over the warning limit are currently: <code>%s</code>

To change the warn mode, use this command again, with one of ban/kick/mute/tban/tmute.
eg: <code>/warnmode ban</code>`, Convert_action(mode, time)))
		return nil
	} else if arg == string("") {
		c.Reply("You need to specify an action to take upon too many warns.")
		return nil
	} else if stringInSlice(args[0], []string{"ban", "mute", "kick", "tban", "tmute"}) {
		if strings.HasPrefix(args[0], "t") {
			if len(args) < 2 {
				c.Reply("It looks like you tried to set time value for warns but you didn't specified time; Try, <code>/setwarnmode [tban/tmute] <timevalue></code>.\n<b>Examples of time value:</b> <code>4m = 4 minutes</code>, <code>3h = 3 hours</code>, <code>6d = 6 days</code>, <code>5w = 5 weeks</code>.")
				return nil
			}
			ctime = Extract_time(c, args[1])
			if ctime == 0 {
				return nil
			}
		}
	} else {
		c.Reply(fmt.Sprintf("Unknown type '%s'. Please use one of: ban/kick/mute/tban/tmute", args[0]))
		return nil
	}
	c.Reply(fmt.Sprintf("✨ Updated warning mode to: %s", Convert_action(args[0], int32(ctime))))
	db.Set_warn_mode(c.Chat().ID, args[0], int(ctime))
	return nil
}

func Set_warn_limit(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	arg := c.Message().Payload
	if arg == string("") && strings.SplitN(c.Message().Text, " ", 2)[0][1:] == "warnlimit" {
		limit, mode, time := db.Get_warn_settings(c.Chat().ID)
		c.Reply(fmt.Sprintf(`The current warn limit is: <code>%d</code>
After this is exceeded, users will be %s.
		
To change the warn limit, use this command again, specifying the amount of warns.
eg: <code>/warnlimit 6</code>`, limit, Convert_action(mode, time)))
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
		fmt.Println(arg)
		db.Set_warn_limit(c.Chat().ID, arg)
		return nil
	}
}

func Warnings_info(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats!")
		return nil
	}
	limit, mode, time := db.Get_warn_settings(c.Chat().ID)
	c.Reply(fmt.Sprintf("There is a %d warning limit in %s. When that limit has been exceeded, the user will be %s.", limit, c.Chat().Title, Convert_action(mode, time)))
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
	data := c.Callback().Data
	user, _ := strconv.Atoi(data)
	user_id := int64(user)
	r := db.Remove_warn(c.Chat().ID, user_id)
	if !r {
		return c.Reply("This warning has already been removed.")
	} else {
		name := "User"
		us, err := c.Bot().ChatByID(user_id)
		if err == nil {
			name = us.FirstName
		}
		return c.Reply(fmt.Sprintf("Admin %s has removed %s's warning.", c.Sender().FirstName, name))
	}
}
