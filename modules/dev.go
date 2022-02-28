package modules

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/amarnathcjd/yoko/modules/db"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	tb "gopkg.in/telebot.v3"
)

func AddSudo(c tb.Context) error {
	if !IsDev(c.Sender().ID) && c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := GetUser(c)
	if u.ID == 0 {
		return nil
	}
	if IsSudo(u.ID) {
		return c.Reply("This user is already a sudo user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't add yourself to the bot admin list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna add myself to the bot admin list!")
	}
	db.AddSudo(u.ID, u.First)
	return c.Reply("Added user to sudo list !")
}

func AddDev(c tb.Context) error {
	if c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := GetUser(c)
	if u.ID == 0 {
		return nil
	}
	if IsDev(u.ID) {
		return c.Reply("This user is already a dev user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't add yourself to the dev list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna add myself to the dev list!")
	}
	db.AddDev(u.ID, u.First)
	return c.Reply("Added user to dev list !")
}

func ListSudo(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	return c.Reply(db.ListSudo())
}

func ListDev(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	return c.Reply(db.ListDev())
}

func RemoveSudo(c tb.Context) error {
	if !IsDev(c.Sender().ID) && c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := GetUser(c)
	if u.ID == 0 {
		return nil
	}
	if !IsSudo(u.ID) {
		return c.Reply("This user is not a sudo user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't remove yourself from the sudo list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna remove myself from the sudo list!")
	}
	db.RemSudo(u.ID)
	return c.Reply("Removed user from sudo list !")
}

func RemoveDev(c tb.Context) error {
	if c.Sender().ID != OWNER_ID {
		return c.Reply("You dont have access to use this !")
	}
	u, _ := GetUser(c)
	if u.ID == 0 {
		return nil
	}
	if !IsDev(u.ID) {
		return c.Reply("This user is not a dev user !")
	} else if u.ID == c.Message().Sender.ID {
		return c.Reply("You can't remove yourself from the dev list !")
	} else if u.ID == BOT_ID {
		return c.Reply("You are a funny one aren't you?, I not gonna remove myself from the dev list!")
	}
	db.RemDev(u.ID)
	return c.Reply("Removed user from dev list !")
}

func Logs(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	} else {
		return c.Reply(&tb.Document{
			File:     tb.File{FileLocal: "log.txt"},
			Caption:  time.Now().String(),
			FileName: "log.txt",
		})
	}
}

func Ping(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	a := time.Now()
	alive := a.Sub(StartTime)
	msg, _ := c.Bot().Send(c.Chat(), "<code>Pinging!</code>")
	b := time.Now()
	_, err := c.Bot().Edit(msg, fmt.Sprintf("<b>► Ping</b>: <code>%s</code>\n<b>► Uptime:</b> %s", b.Sub(a).String(), alive.Truncate(time.Second).String()))
	return err
}

func Stats(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	var Stats = db.GatherStats()
	Stats += "\n\n"
	Sys := GetSystemDetails()
	Stats += "<b>CPU:</b> " + Sys.CPU + "\n"
	Stats += "<b>RAM:</b> " + Sys.Memory + "\n"
	Stats += "<b>Disk:</b> " + Sys.Disk + "\n"
	Stats += "<b>Hostname:</b> " + Sys.Hostname + "\n"
	Stats += "<b>OS:</b> " + Sys.Platform + "\n"
	Stats += "<b>Uptime:</b> " + fmt.Sprint() + "\n"
	return c.Reply(Stats)
}

func Json(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	if c.Message().IsReply() {
		b, _ := json.MarshalIndent(c.Message().ReplyTo, "", "    ")
		return c.Reply(string(b))
	} else {
		b, _ := json.MarshalIndent(c.Message(), "", "    ")
		return c.Reply(string(b))
	}
}

func SendMessage(c tb.Context) error {
	if !IsBotAdmin(c.Sender().ID) {
		return nil
	}
	Args := GetArgs(c)
	var Chat *tb.Chat
	if len(Args) < 2 {
		Chat = c.Chat()
	} else {
		Arg := strings.Split(Args, " ")
		if isInt(Arg[0]) {
			ID, _ := strconv.Atoi(Arg[0])
			Chat = &tb.Chat{ID: int64(ID)}
		} else {
			Chat = c.Chat()
		}
	}
	_, err := c.Bot().Send(Chat, Args)
	return err
}

type SysInfo struct {
	Hostname string
	Platform string
	Uptime   string
	CPU      string
	Memory   string
	Disk     string
}

func GetSystemDetails() SysInfo {
	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")
	return SysInfo{
		Hostname: hostStat.HostID,
		Platform: hostStat.Platform,
		Uptime:   time.Now().Sub(time.Unix(int64(hostStat.BootTime), 0)).String(),
		CPU:      cpuStat[0].ModelName + "(" + fmt.Sprint(cpuStat[0].Cores) + ")",
		Memory:   fmt.Sprintf("%v/%v", ByteCount(int64(vmStat.Used)), ByteCount(int64(vmStat.Total))),
		Disk:     fmt.Sprintf("%v/%v", ByteCount(int64(diskStat.Used)), ByteCount(int64(diskStat.Total))),
	}
}
