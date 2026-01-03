package configpanel

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/conversation"
	"github.com/Jisin0/autofilterbot/pkg/panel"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ChannelFieldOpts provides optional parameters to ChannelField.
type ChannelFieldOpts struct {
	// Description for the field.
	Description string
	// Maximum number of channels allowed.
	MaxAmount int
	// Indicates whether the user should be asked if a request invite link should be generated.
	AllowRequestInvite bool
}

func ChannelField(app AppPreview, fieldName string, opts ChannelFieldOpts) panel.CallbackFunc {
	return func(ctx *panel.Context) (string, [][]gotgbot.InlineKeyboardButton, error) {
		var (
			op   string
			data = ctx.CallbackData
		)

		if len(data.Args) != 0 {
			op = data.Args[0]
		}

		currentChannels := app.GetConfig().GetFsubChannels()

		switch op {
		case OperationDelete:
			if len(data.Args) < 2 {
				return "", nil, errors.New("configpanel: channel: insufficient data for delete operation")
			}

			channelID, err := strconv.ParseInt(data.Args[1], 10, 64)
			if err != nil {
				return "", nil, err
			}

			for i, c := range currentChannels {
				if c.ID == channelID {
					currentChannels = slices.Delete(currentChannels, i, i+1)

					app.GetDB().UpdateConfig(ctx.Bot.Id, config.FieldNameFsub, currentChannels)
					go app.RefreshConfig()

					return fieldName + " Channel was Deleted Successfully ✅", nil, nil
				}
			}

			return fieldName + " Channel to Delete Was not Found 🫤", nil, nil
		case OperationReset:
			conv := conversation.NewConversatorFromUpdate(ctx.Bot, ctx.Update.Update)

			m, err := conv.Ask(app.GetContext(), "Are you sure you want to delete all Channels? (y/N)", nil)
			if err != nil {
				return "", nil, errors.Wrap(err, "configpanel: channel: send reset confirmation message failed")
			}

			if strings.ToLower(m.Text) != "y" {
				return "Reset Operation Cancelled!", nil, nil
			}

			app.GetDB().ResetConfig(ctx.Bot.Id, config.FieldNameFsub)
			go app.RefreshConfig()

			return fieldName + " Channels Have Been Reset Succesfully ✅", nil, nil
		case OperationSet:
			if opts.MaxAmount != 0 && len(currentChannels) >= opts.MaxAmount {
				ctx.CallbackQuery.Answer(ctx.Bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Channel Limit Reached.\n\nPlease delete a value to try again.", ShowAlert: true})
				return "", nil, nil
			}

			chatIDStr, _ := data.GetArg(1)
			chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

			if chatID == 0 {
				conv := conversation.NewConversatorFromUpdate(ctx.Bot, ctx.Update.Update)

				m, err := conv.Ask(app.GetContext(), "Please Forward a Post from the Channel (with quotes) or Send the Chat id in the Format -100xxxxxxx: ", nil)
				if err != nil {
					return "", nil, errors.Wrap(err, "configpanel: channel: send channel request message failed")
				}

				if m.ForwardOrigin != nil {
					if f, ok := m.ForwardOrigin.(gotgbot.MessageOriginChannel); ok {
						chatID = f.Chat.Id
					}
				} else {
					chatID, _ = strconv.ParseInt(strings.TrimSpace(m.Text), 10, 64)
				}
			}

			if chatID == 0 {
				return "Message was not forwarded from a channel nor contains a channel ID!", nil, nil
			}

			chat, err := ctx.Bot.GetChat(chatID, nil)
			if err != nil {
				return "", nil, err
			}

			for _, c := range currentChannels {
				if c.ID == chat.Id {
					return "New channel is already add!", nil, nil
				}
			}

			var isRequest bool

			if s, ok := data.GetArg(2); ok {
				isRequest, err = strconv.ParseBool(s)
				if err != nil {
					return "", nil, errors.Wrap(err, "configpanel: channel: failed to parse 'is request' arg")
				}
			} else if opts.AllowRequestInvite {
				return "<b>Would you like to add the channel as a request invite channel?</b> \n\nThe user will send a join request, requiring admin approval, instead of joining directly.",
					[][]gotgbot.InlineKeyboardButton{{
						{Text: "Yes", CallbackData: ctx.CallbackData.AddArg(fmt.Sprint(chatID)).AddArg("1").ToString()},
						{Text: "No", CallbackData: ctx.CallbackData.AddArg(fmt.Sprint(chatID)).AddArg("0").ToString()},
					}},
					nil
			}

			link, err := ctx.Bot.CreateChatInviteLink(chat.Id, &gotgbot.CreateChatInviteLinkOpts{Name: fieldName, CreatesJoinRequest: isRequest})
			if err != nil {
				app.GetLog().Debug("configpanel: channel: failed to generate invite link", zap.Int64("id", chat.Id), zap.Error(err))
				return "Failed to Create Invite Link. Please Make Sure the bot has Permissions to Add Users", nil, nil
			}

			currentChannels = append(currentChannels, model.Channel{
				ID:                 chat.Id,
				Title:              chat.Title,
				InviteLink:         link.InviteLink,
				CreatesJoinRequest: isRequest,
			})

			app.GetDB().UpdateConfig(ctx.Bot.Id, config.FieldNameFsub, currentChannels)
			go app.RefreshConfig()

			return fmt.Sprintf("%s has been Saved as a %s Channel Successfully ✅", chat.Title, fieldName), nil, nil
		case OperationRefresh:
			if len(data.Args) < 2 {
				return "", nil, errors.New("configpanel: channel: insufficient data for refresh operation")
			}

			channelID, err := strconv.ParseInt(data.Args[1], 10, 64)
			if err != nil {
				return "", nil, err
			}

			chat, err := ctx.Bot.GetChat(channelID, nil)
			if err != nil {
				return "", nil, err
			}

			var (
				channel      model.Channel
				channelIndex *int
			)

			for i, c := range currentChannels {
				if c.ID != channelID {
					continue
				}

				channel = c
				channelIndex = &i
			}

			if channelIndex == nil {
				return "Channel was not found in saved channels!", nil, nil
			}

			link, err := ctx.Bot.CreateChatInviteLink(chat.Id, &gotgbot.CreateChatInviteLinkOpts{Name: fieldName, CreatesJoinRequest: channel.CreatesJoinRequest})
			if err != nil {
				app.GetLog().Debug("configpanel: channel: failed to generate invite link", zap.Int64("id", chat.Id), zap.Error(err))
				return "Failed to Create Invite Link. Please Make Sure the bot has Permissions to Add Users", nil, nil
			}

			currentChannels[*channelIndex].Title = chat.Title
			currentChannels[*channelIndex].InviteLink = link.InviteLink

			app.GetDB().UpdateConfig(ctx.Bot.Id, config.FieldNameFsub, currentChannels)
			go app.RefreshConfig()

			return "Channel Information has been Updated Successfully ✅", nil, nil
		default:
			var s strings.Builder

			if opts.Description != "" {
				s.WriteString("ℹ️ <i>" + opts.Description + "</i>\n\n")
			}

			s.WriteString(`<b><u>Options</u></b>
<b>Refresh</b> - Refresh channel information (title and invite link)
<b>Add</b> - Add a new channel
<b>Delete</b> - Delete a single channel.
<b>Reset</b> - Reset to default`)

			if opts.MaxAmount != 0 {
				s.WriteString(fmt.Sprintf("\n\n<b>🗒️ Upto %d channel(s) can be added.</b>", opts.MaxAmount))
			}

			var keybaord [][]gotgbot.InlineKeyboardButton

			for _, c := range currentChannels {
				keybaord = append(keybaord, []gotgbot.InlineKeyboardButton{{Text: c.Title, Url: c.InviteLink}})
				keybaord = append(keybaord, []gotgbot.InlineKeyboardButton{
					{Text: "🗑️ Delete", CallbackData: ctx.CallbackData.AddArgs(OperationDelete, fmt.Sprint(c.ID)).ToString()},
					{Text: "🔄 Refresh", CallbackData: ctx.CallbackData.AddArgs(OperationRefresh, fmt.Sprint(c.ID)).ToString()},
				})
			}

			keybaord = append(keybaord, []gotgbot.InlineKeyboardButton{
				{Text: "⏪ Reset", CallbackData: ctx.CallbackData.AddArg(OperationReset).ToString()},
				{Text: "➕ Add", CallbackData: ctx.CallbackData.AddArgs(OperationSet).ToString()},
			})

			return s.String(), keybaord, nil
		}
	}
}
