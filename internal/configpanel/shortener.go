package configpanel

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/pkg/conversation"
	"github.com/Jisin0/autofilterbot/pkg/panel"
	"github.com/Jisin0/autofilterbot/pkg/shortener"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

func Shortener(app AppPreview) panel.CallbackFunc {
	return func(ctx *panel.Context) (string, [][]gotgbot.InlineKeyboardButton, error) {
		current := app.GetConfig().GetShortener()

		op, _ := ctx.CallbackData.GetArg(0)

		switch op {
		case OperationSet:
			conv := conversation.NewConversatorFromUpdate(ctx.Bot, ctx.Update.Update)

			m, err := conv.Ask(app.GetContext(), "Please Send the URL of the Homepage of your Shortener of Choice for example: gplinks.com (or cancel to abort):", nil)
			if err != nil {
				return "An Unknown Error Occured :/", nil, err
			}

			rawURL := m.Text

			// not optimal but whatever
			if !strings.HasPrefix(m.Text, "https://") {
				rawURL = "https://" + strings.TrimPrefix(rawURL, "http://")
			}

			parsedURL, err := url.Parse(rawURL)
			if err != nil || parsedURL.Host == "" {
				app.GetLog().Debug("configpanel: shortener: failed to parse URL", zap.String("raw_url", rawURL), zap.Error(err))
				return "Failed to parse URL :/\n\nPlease make sure the URL is in a valid format, just copy from your browser's address bar to prevent errors.", nil, nil
			}

			shortenerURL := parsedURL.Host

			m2, err := conv.Ask(app.GetContext(), fmt.Sprintf("Great Work! Now Please Send your API Key Commonly Obtained from %s/member/tools/api:", shortenerURL), nil)
			if err != nil {
				return "An Unknown Error Occured :/", nil, err
			}

			apiKey := m2.Text
			s := shortener.NewShortener(apiKey, shortenerURL)

			_, err = s.ShortenURL("telegram.dog/FractalProjects")
			if err != nil {
				app.GetLog().Warn("configpanel: shortener: testing shortener failed", zap.String("url", shortenerURL), zap.Error(err))
				return "Failed to Shorten URLs with your API key :/\n\nPlease Double Check the Details Provided or Contact Devs.", nil, nil
			}

			err = app.GetDB().UpdateConfig(ctx.Bot.Id, config.FieldNameShortener, s)
			if err != nil {
				app.GetLog().Warn("configpanel: shortener: failed to save shortener details", zap.Error(err))
				return "Saving Shortener Details Failed :/\n\nPlease Check App Logs for More Info.", nil, nil
			}

			go app.RefreshConfig()

			return "✅ URl Shortener Saved Successfully!", nil, nil
		case OperationDelete:
			err := app.GetDB().ResetConfig(ctx.Bot.Id, config.FieldNameShortener)
			go app.RefreshConfig()

			return "✅ URL Shortener Removed Successfully", nil, err // error message is given priority if error is not nil
		default:
			if current != nil {
				text := fmt.Sprintf(`
Current URL Shortener Configuration:

<b>URL</b>: <code>%s</code>
<b>API Key</b>: <tg-spoiler>%s</tg-spoiler>`, current.RootURL, current.ApiKey)

				return text, [][]gotgbot.InlineKeyboardButton{{{Text: "Remove", CallbackData: ctx.CallbackData.RemoveArgs().AddArg(OperationDelete).ToString()}}}, nil
			}

			text := `
ℹ️ URL Shorteners make the User go Through Verification Pages and/or Adwalls in Order to get their Files. This Reduces Bot Abuse and can Also Generate Revenue.

To get started, pick a shortener of choice that follows the same schema for example: gplinks.com

Now click "Add" to start adding your shortener.`

			return text, [][]gotgbot.InlineKeyboardButton{{{Text: "Add", CallbackData: ctx.CallbackData.RemoveArgs().AddArg(OperationSet).ToString()}}}, nil
		}

	}
}
