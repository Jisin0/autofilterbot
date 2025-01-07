package autofilter

import (
	"regexp"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// IsBadQuery report whether the query is bad based on the message entities and length of the text.
func IsBadQuery(text string, entities []gotgbot.MessageEntity) bool {
	if len(text) > 30 || len(text) < 2 { //TODO: might wanna change these hardocded values
		return true
	}

	if text[0] == '!' { // starting a message with ! is a fail safe to not process it
		return true
	}

	for _, e := range entities {
		switch e.Type {
		case "url", "bot_command", "mention", "text_link", "text_mention", "email", "phone_number":
			return true
		}
	}

	return false
}

var nonAlphaNumeric = regexp.MustCompile(`[^\w\s]`)

// Sanitize removes all unwanted character from the query.
//
// NOTE: Also changes string to lower case which may have side effects.
func Sanitize(s string) string {
	// removes all non-alphanumeric characters
	return strings.Trim(nonAlphaNumeric.ReplaceAllString(strings.ToLower(s), ""), " ")
}
