package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	chatToken := refreshAuth(os.Getenv("twitch-bot-refresh-token"), os.Getenv("twitch-bot-client-id"), os.Getenv("twitch-bot-client-secret"))
	bot := newBot(chatToken, []string{"#supcole", "#zeruel132"}, "supcole")
	bot.Connect()

	for {
		message := <-bot.onMessage

		fmt.Printf("   \033[1;45m%s %s:\033[0m   %s\n", message.channel, message.user.username, message.contents)
		if strings.HasPrefix(strings.ToLower(message.contents), "!hello") {
			bot.Message(message.channel, "Hello world")

		}
	}
}
