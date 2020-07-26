# Go Twitch IRC Client

This repository contains a Go client for interacting with Twitch IRC channels.

# Example

main.go contains the following example which will receive each message in my channel (supcole) and log information about it and the user who sent it to the console. To be able to run this, you must have Twitch API credentials.

```go

chatToken := twitchircgo.RefreshAuth("refresh token goes here","client id goes here", "client secret goes here")

bot := newBot(chatToken, []string{"#channel", "#anotherchannel"}, "botusername")
bot.Connect()

for {
    message := <-bot.OnMessage
    fmt.Printf("got message %+v\n",message)
}


```

Messages in the above code segment are sent to the channel OnMessage, and saved to a local variable.
Message objects contain user information and the contents of the message itself.
