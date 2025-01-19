package button

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
)

const (
	maxButtonsPerRow = 10
)

var buttonRegex = regexp.MustCompile(`\[([^\[]+?)\]\((url|cmd|inline|copy):(?:/{0,2})?(.*?)\)`)

// ParseFromText parses text into keyboard buttons to be saved to the database or used in a message.
//
// Text must follow the syntax: [<label>](<type>:<value>)
//
// Possible types:
//   - url: value is the target url
//   - cmd: value is the command to be redirected to for ex: start, about, help, privacy etc.
//   - inline: value is the search query and can be empty
//   - copy: value is text to be copied when clicked.
//
// label is always required; value is required unless explicitly stated above.
// Buttons are split into different rows using new lines or the enter key.
// Returns the text with parsed buttons removed and the buttons parsed or an error.
func ParseFromText(text string) (returnText string, buttons [][]InlineKeyboardButton, err error) {
	returnText = text

	var allErrors []error

	for _, textRows := range strings.Split(text, "\n") {
		find := buttonRegex.FindAllStringSubmatch(textRows, maxButtonsPerRow)
		row := make([]InlineKeyboardButton, 0, len(find))

		for _, m := range find {
			if len(m) < 4 {
				// what each index in m is:
				//	- 0: raw text
				// 	- 1: label
				//	- 2: type
				//	- 3: value
				continue
			}

			var (
				newButton = InlineKeyboardButton{Text: m[1]}
				btnType   = m[2]
				val       = m[3]
			)

			if val == "" && btnType != "inline" { // only inline can have empty value
				allErrors = append(allErrors, fmt.Errorf("value is required for %s button", btnType))
			}

			switch btnType {
			case "url":
				newButton.Url = val
			case "cmd":
				newButton.CallbackData = "cmd" + string(callbackdata.PathDelimiter) + val
			case "inline":
				newButton.IsInline = true
				newButton.SwitchInlineQueryCurrentChat = val
			case "copy":
				newButton.CopyText = val
			default: // ideally shouldnt happen as regex pattern wouldnt match
				allErrors = append(allErrors, fmt.Errorf("invalid button type %s", btnType))
				continue
			}

			row = append(row, newButton)
			returnText = strings.Replace(returnText, m[0], "", 1)
		}

		if len(row) != 0 {
			buttons = append(buttons, row)
		}
	}

	returnText = strings.TrimSpace(returnText)

	return returnText, buttons, errors.Join(allErrors...)
}
