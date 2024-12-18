package handlers

import (
	"regexp"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var stringFormatRegex = regexp.MustCompile(`\{(\w+)\}`)

// FormatString formats a Pythonic syntax string using a map of key-value pairs.
func FormatString(template string, values map[string]string) string {
	result := stringFormatRegex.ReplaceAllStringFunc(template, func(match string) string {
		key := stringFormatRegex.FindStringSubmatch(match)[1] // Extract key from {key}
		if value, ok := values[key]; ok {
			return value
		}
		return match // Leave unchanged as no matching key is found
	})

	return result
}

// hasMedia reports whether message has media.
func hasMedia(m *gotgbot.Message) bool {
	return m.Photo != nil || m.Audio != nil || m.Document != nil || m.Video != nil || m.Animation != nil || m.Audio != nil
}
