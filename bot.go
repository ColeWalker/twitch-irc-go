package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
	"time"
)

//The Bot object which stores irc info
type Bot struct {
	// server to connect to
	server string
	// port to connect to
	port string
	// nickname of bot
	nick string
	// channel to join - #username format
	channel string
	// mods of channel -> currently unused
	mods map[string]bool
	// connection to irc
	conn net.Conn
	// oauth token containing chat:read and chat:write scopes
	AuthToken string
	// Event listeners
	onMessage map[string]chan Message
}

//constructor
func newBot(token string, channel string, nick string) *Bot {
	defaultChannel := make(chan Message)
	onMessage := make(map[string]chan Message)

	onMessage["default"] = defaultChannel
	return &Bot{
		server:    "irc.twitch.tv",
		port:      "6667",
		nick:      nick,
		channel:   channel,
		mods:      make(map[string]bool),
		conn:      nil,
		AuthToken: token,
		onMessage: onMessage}
}

//Connect and connect bot to IRC server
func (bot *Bot) Connect() {
	//close connection if it exists
	if bot.conn != nil {
		err := bot.conn.Close()
		if err != nil {
			fmt.Println("Error closing existing connection.")
		}
	}

	var err error
	fmt.Println("Attempting to connect to Twitch IRC server...")
	bot.conn, err = net.Dial("tcp", bot.server+":"+bot.port)

	if err != nil {
		fmt.Printf("Unable to connect to Twitch IRC server! Reconnecting in 10 seconds...\n")
		time.Sleep(10 * time.Second)
		bot.Connect()
	}

	fmt.Printf("Connected to IRC server %s\n", bot.server)

	// Initial pings need to be set
	fmt.Fprintf(bot.conn, "USER %s 8 * :%s\r\n", bot.nick, bot.nick)
	fmt.Fprintf(bot.conn, "PASS oauth:%s\r\n", bot.AuthToken)
	fmt.Fprintf(bot.conn, "NICK %s\r\n", bot.nick)
	fmt.Fprintf(bot.conn, "JOIN %s\r\n", bot.channel)
	fmt.Fprintf(bot.conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(bot.conn, "CAP REQ :twitch.tv/tags\r\n")

	go bot.ReadLoop()
}

//Message to IRC channel
func (bot *Bot) Message(message string) {
	if message != "" {
		fmt.Fprintf(bot.conn, "PRIVMSG "+bot.channel+" :"+message+"\r\n")
	}
	fmt.Println("PRIVMSG " + bot.channel + " :" + message + "\r\n")
}

//ReadLoop for Bot
func (bot *Bot) ReadLoop() {
	reader := bufio.NewReader(bot.conn)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLine()

		if err != nil {
			fmt.Println("Bot read loop exited due to error")
			bot.Connect()
			break
		} else if strings.Contains(line, "PING") {
			fmt.Fprintf(bot.conn, "PONG :tmi.twitch.tv")
		} else if strings.Contains(line, ".tmi.twitch.tv PRIVMSG "+bot.channel) {
			message := parseMessage(line, bot.channel)
			bot.onMessage["default"] <- *message
		}
	}

}
