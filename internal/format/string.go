package format

import (
	"regexp"
)

var keyValueRegex = regexp.MustCompile(`\{(\w+)\}`)

// KeyValueFormat formats a template string using a pythonic syntax from values in the input map.
func KeyValueFormat(template string, values map[string]string) string {
	return keyValueRegex.ReplaceAllStringFunc(template, func(match string) string {
		key := match[1 : len(match)-1] // extract key from {key}
		if value, ok := values[key]; ok {
			return value
		}

		return match // leave as is if no matching key is found
	})
}
