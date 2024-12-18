package configpanel

import (
	"github.com/Jisin0/autofilterbot/pkg/configpanel/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var closeButton = gotgbot.InlineKeyboardButton{
	Text:         "X ᴄʟᴏsᴇ X",
	CallbackData: "close",
}

// buttonsFromPages creates keyboard button with subpages from given CallbackData.
func buttonsFromPages(callbackData callbackdata.CallbackData, pages []*Page) [][]gotgbot.InlineKeyboardButton {
	var backRow []gotgbot.InlineKeyboardButton

	if len(callbackData.Path) <= 1 {
		// root page so add close button
		backRow = []gotgbot.InlineKeyboardButton{closeButton}
	} else {
		// nested page so add back button
		backRow = []gotgbot.InlineKeyboardButton{backButton(callbackData.RemoveLastPath().ToString())}
	}

	if len(pages) == 0 {
		return [][]gotgbot.InlineKeyboardButton{backRow}
	}

	// determine the number of buttons per row (2 by default, 3 if divisible by 3)
	buttonsPerRow := 2
	if len(pages)%3 == 0 {
		buttonsPerRow = 3
	}

	totalButtons := make([][]gotgbot.InlineKeyboardButton, 0, (len(pages)+buttonsPerRow-1)/buttonsPerRow)

	for i := 0; i < len(pages); i += buttonsPerRow {
		end := i + buttonsPerRow
		if end > len(pages) {
			end = len(pages)
		}

		row := make([]gotgbot.InlineKeyboardButton, 0, end-i)
		for _, page := range pages[i:end] {
			row = append(row, gotgbot.InlineKeyboardButton{
				Text:         page.DisplayName,
				CallbackData: callbackData.AddPath(page.Name).ToString(),
			})
		}
		totalButtons = append(totalButtons, row)
	}

	totalButtons = append(totalButtons, backRow)

	return totalButtons
}

// backButton creates a button with back text with given callback data.
func backButton(callbackData string) gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         "<- Back",
		CallbackData: callbackData,
	}
}
