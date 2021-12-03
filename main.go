package main

import ("fmt"
    "github.com/PaulSonOfLars/gotgbot/v2", "os", "net/http")

func main() {
    fmt.Println("lel")
    bot, err := gotgbot.NewBot(
		os.Getenv("TOKEN"),
		&gotgbot.BotOpts{
			Client:      http.Client{},
			GetTimeout:  gotgbot.DefaultGetTimeout,
			PostTimeout: gotgbot.DefaultPostTimeout,
		},
	)
    fmt.Println("done")
}
