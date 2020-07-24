package twitchircgo

import (
	"strconv"
	"strings"
)

// Message represents message received
type Message struct {
	channel  string
	user     User
	contents string
}

// User represents user sending message
type User struct {
	username   string
	moderator  bool
	badges     []Badge
	id         string
	color      string
	subscriber bool
	turbo      bool
	owner      bool
}

// Badge represents badges sent in twitch chat
type Badge struct {
	name  string
	value string
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
				name:  badgeInfo[0],
				value: badgeInfo[1]}

			badges = append(badges, *badge)
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
		username:   tagMap["display-name"],
		moderator:  isMod,
		badges:     badges,
		id:         tagMap["user-id"],
		color:      tagMap["color"],
		owner:      isOwner,
		subscriber: isSubscriber,
		turbo:      isTurbo}
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
		channel:  channel,
		contents: contents,
		user:     *user}
}
