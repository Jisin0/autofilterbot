package send

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var ErrUnknownMessageType = errors.New("msg is of an unsupported type")

// ResendMessage resends an exact copy of the given *gotgbot.Message to a chat.
// Supports 15 different message types excluding paid media and invoice due to complexity.
func ResendMessage(bot *gotgbot.Bot, msg *gotgbot.Message, chatID int64) (*gotgbot.Message, error) {
	switch {
	case msg.Animation != nil:
		return bot.SendAnimation(chatID, gotgbot.InputFileByID(msg.Animation.FileId), &gotgbot.SendAnimationOpts{
			Duration: msg.Animation.Duration,
			Width:    msg.Animation.Width,
			Height:   msg.Animation.Height,
			// Thumbnails cannot be resent.
			Caption:               msg.Caption,
			CaptionEntities:       msg.CaptionEntities,
			ShowCaptionAboveMedia: msg.ShowCaptionAboveMedia,
			HasSpoiler:            msg.HasMediaSpoiler,
			ProtectContent:        msg.HasProtectedContent,
			ReplyMarkup:           msg.ReplyMarkup,
		})
	case msg.Audio != nil:
		return bot.SendDocument(chatID, gotgbot.InputFileByID(msg.Document.FileId), &gotgbot.SendDocumentOpts{
			Caption:         msg.Caption,
			CaptionEntities: msg.CaptionEntities,
			ProtectContent:  msg.HasProtectedContent,
			ReplyMarkup:     msg.ReplyMarkup,
		})
	case msg.Contact != nil:
		return bot.SendContact(chatID, msg.Contact.PhoneNumber, msg.Contact.FirstName, &gotgbot.SendContactOpts{
			LastName:       msg.Contact.LastName,
			Vcard:          msg.Contact.Vcard,
			ProtectContent: msg.HasProtectedContent,
			ReplyMarkup:    msg.ReplyMarkup,
		})
	case msg.Dice != nil:
		return bot.SendDice(chatID, &gotgbot.SendDiceOpts{
			Emoji:          msg.Dice.Emoji,
			ProtectContent: msg.HasProtectedContent,
			ReplyMarkup:    msg.ReplyMarkup,
		})
	case msg.Document != nil:
		return bot.SendDocument(chatID, gotgbot.InputFileByID(msg.Document.FileId), &gotgbot.SendDocumentOpts{
			Caption:         msg.Caption,
			CaptionEntities: msg.CaptionEntities,
			ProtectContent:  msg.HasProtectedContent,
			ReplyMarkup:     msg.ReplyMarkup,
		})
	case msg.Game != nil:
		return bot.SendGame(chatID, msg.Game.Title, &gotgbot.SendGameOpts{
			ProtectContent: msg.HasProtectedContent,
			ReplyMarkup:    *msg.ReplyMarkup, // possible panic. however, all games should have a markup.
		})
	// skipped invoice, too annoying to implement
	case msg.Location != nil:
		return bot.SendLocation(chatID, msg.Location.Latitude, msg.Location.Longitude, &gotgbot.SendLocationOpts{
			HorizontalAccuracy:   msg.Location.HorizontalAccuracy,
			LivePeriod:           msg.Location.LivePeriod,
			Heading:              msg.Location.Heading,
			ProximityAlertRadius: msg.Location.ProximityAlertRadius,
			ProtectContent:       msg.HasProtectedContent,
			ReplyMarkup:          msg.ReplyMarkup,
		})
	case msg.Text != "":
		return bot.SendMessage(chatID, msg.Text, &gotgbot.SendMessageOpts{
			Entities:           msg.Entities,
			LinkPreviewOptions: msg.LinkPreviewOptions,
			ProtectContent:     msg.HasProtectedContent,
			ReplyMarkup:        msg.ReplyMarkup,
		})
	// skipped paid media, too complicated (but doable)
	case msg.Photo != nil:
		return bot.SendPhoto(chatID, gotgbot.InputFileByID(msg.Photo[0].FileId), &gotgbot.SendPhotoOpts{
			Caption:               msg.Caption,
			CaptionEntities:       msg.CaptionEntities,
			ShowCaptionAboveMedia: msg.ShowCaptionAboveMedia,
			HasSpoiler:            msg.HasMediaSpoiler,
			ProtectContent:        msg.HasProtectedContent,
			ReplyMarkup:           msg.ReplyMarkup,
		})
	case msg.Poll != nil:
		return bot.SendPoll(chatID, msg.Poll.Question, convertPollOptions(msg.Poll.Options), &gotgbot.SendPollOpts{
			QuestionEntities:      msg.Poll.QuestionEntities,
			IsAnonymous:           msg.Poll.IsAnonymous,
			Type:                  msg.Poll.Type,
			AllowsMultipleAnswers: msg.Poll.AllowsMultipleAnswers,
			CorrectOptionId:       msg.Poll.CorrectOptionId,
			Explanation:           msg.Poll.Explanation,
			ExplanationEntities:   msg.Poll.ExplanationEntities,
			OpenPeriod:            msg.Poll.OpenPeriod,
			CloseDate:             msg.Poll.CloseDate,
			IsClosed:              msg.Poll.IsClosed,
			ProtectContent:        msg.HasProtectedContent,
			ReplyMarkup:           msg.ReplyMarkup,
		})
	case msg.Sticker != nil:
		return bot.SendSticker(chatID, gotgbot.InputFileByID(msg.Sticker.FileId), &gotgbot.SendStickerOpts{
			Emoji:          msg.Sticker.Emoji,
			ProtectContent: msg.HasProtectedContent,
			ReplyMarkup:    msg.ReplyMarkup,
		})
	case msg.Venue != nil:
		return bot.SendVenue(chatID, msg.Venue.Location.Latitude, msg.Venue.Location.Longitude, msg.Venue.Title, msg.Venue.Address, &gotgbot.SendVenueOpts{
			FoursquareId:    msg.Venue.FoursquareId,
			GooglePlaceId:   msg.Venue.GooglePlaceId,
			FoursquareType:  msg.Venue.FoursquareType,
			GooglePlaceType: msg.Venue.GooglePlaceType,
			ProtectContent:  msg.HasProtectedContent,
			ReplyMarkup:     msg.ReplyMarkup,
		})
	case msg.Video != nil:
		return bot.SendVideo(chatID, gotgbot.InputFileByID(msg.Video.FileId), &gotgbot.SendVideoOpts{
			Duration:              msg.Video.Duration,
			Width:                 msg.Video.Width,
			Height:                msg.Animation.Height,
			Caption:               msg.Caption,
			CaptionEntities:       msg.CaptionEntities,
			ShowCaptionAboveMedia: msg.ShowCaptionAboveMedia,
			HasSpoiler:            msg.HasMediaSpoiler,
			ProtectContent:        msg.HasProtectedContent,
			ReplyMarkup:           msg.ReplyMarkup,
			//	SupportsStreaming: , not sure:/
		})
	case msg.VideoNote != nil:
		return bot.SendVideoNote(chatID, gotgbot.InputFileByID(msg.VideoNote.FileId), &gotgbot.SendVideoNoteOpts{
			Duration:       msg.Video.Duration,
			Length:         msg.VideoNote.Length,
			ProtectContent: msg.HasProtectedContent,
			ReplyMarkup:    msg.ReplyMarkup,
			//	SupportsStreaming: , not sure:/
		})
	case msg.Voice != nil:
		return bot.SendVoice(chatID, gotgbot.InputFileByID(msg.Voice.FileId), &gotgbot.SendVoiceOpts{
			Duration:        msg.Video.Duration,
			Caption:         msg.Caption,
			CaptionEntities: msg.CaptionEntities,
			ProtectContent:  msg.HasProtectedContent,
			ReplyMarkup:     msg.ReplyMarkup,
		})
	default:
		return nil, ErrUnknownMessageType
	}
}

func convertPollOptions(i []gotgbot.PollOption) []gotgbot.InputPollOption {
	var l []gotgbot.InputPollOption
	for _, o := range i {
		l = append(l, gotgbot.InputPollOption{
			Text:         o.Text,
			TextEntities: o.TextEntities,
		})
	}
	return l
}
