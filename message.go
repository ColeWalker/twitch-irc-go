package twitchircgo

import (
	"strconv"
	"strings"
)

// Message represents message received
type Message struct {
	Channel  string
	User     User
	Contents string
}

// User represents user sending message
type User struct {
	// Username of sender
	Username string
	// Moderator is true if sender is a Moderator
	Moderator bool
	// Badges contains user's badges
	Badges []Badge
	// ID represents the unique user id (static where username can change)
	ID string
	// Color represents the user's chosen chat color
	Color string
	// Subscriber represents whether user is subbed to channel
	Subscriber bool
	// Turbo represents whether user is a Turbo user
	Turbo bool
	// Owner represents whether user is the owner of the channel, owners are NOT moderators
	Owner bool
}

// Badge represents badges sent in twitch chat
type Badge struct {
	Name  string
	Value string
}

func userFromTags(tags []string, channel string) *User {
	tagMap := make(map[string]string)

	for _, v := range tags {
		valuePair := strings.Split(v, "=")
		key := valuePair[0]
		value := valuePair[1]

		tagMap[key] = value
	}

	isMod, err := strconv.ParseBool(tagMap["mod"])

	if err != nil {
		isMod = false
	}

	rawBadges := strings.Split(tagMap["badges"], ",")

	var badges []Badge
	for _, v := range rawBadges {
		badgeInfo := strings.Split(v, "/")
		if len(badgeInfo) > 1 {
			badge := &Badge{
				Name:  badgeInfo[0],
				Value: badgeInfo[1]}

			badges = append(Badges, *badge)
		}

	}

	isSubscriber, err := strconv.ParseBool(tagMap["subscriber"])
	if err != nil {
		isSubscriber = false
	}

	isTurbo, err := strconv.ParseBool(tagMap["turbo"])
	if err != nil {
		isTurbo = false
	}

	isOwner := (strings.ToLower(tagMap["display-name"]) == strings.TrimPrefix(channel, "#"))

	return &User{
		Username:   tagMap["display-name"],
		Moderator:  isMod,
		Badges:     badges,
		ID:         tagMap["user-id"],
		Color:      tagMap["color"],
		Owner:      isOwner,
		Subscriber: isSubscriber,
		Turbo:      isTurbo}
}

func parseMessage(message string, channel string) *Message {
	messageContents := strings.Split(message, ".tmi.twitch.tv PRIVMSG "+channel)

	rawValues := strings.SplitN(messageContents[0], ":", 2)
	contents := messageContents[1][2:len(messageContents[1])]

	tags := strings.Split(rawValues[0], ";")

	// timeRegex := regexp.MustCompile(`tmi-sent-ts=\d*`)
	// timestampPair := timeRegex.FindString(rawValues[0])

	// timestamp, err := strconv.ParseInt(strings.Split(timestampPair, "=")[1], 10, 64)

	// if err != nil {
	// 	timestamp = 0
	// }

	user := userFromTags(tags, channel)

	return &Message{
		Channel:  channel,
		Contents: contents,
		User:     *user}
}
