package main

import (
	"fmt"
	"os"
)

func main() {
	chatToken := refreshAuth(os.Getenv("twitch-bot-refresh-token"), os.Getenv("twitch-bot-client-id"), os.Getenv("twitch-bot-client-secret"))
	bot := newBot(chatToken, "#supcole", "supcole")
	bot.Connect()

	for {
		message := <-bot.onMessage["default"]
		fmt.Printf("got message %+v\n",message)
	}
}
