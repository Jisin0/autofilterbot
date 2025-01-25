package fsub

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// FsubChannel is a single force sub channel.
type FsubChannel struct {
	model.FsubChannel
}

// IsJoined reports wether the user is a member of the channel.
//
// - bot: bot client to make api requests.
// - userId: id of user to check for.
//
// NOTE: If some error other than "user not found" occurs during api call function returns true and error.
func IsMember(bot *gotgbot.Bot, chatId, userId int64) (bool, error) {
	member, err := bot.GetChatMember(chatId, userId, nil)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return false, nil
		}

		// if other error returns true with error
		return true, fmt.Errorf("fsub: failed to get chat member: %w", err)
	}

	switch m := member.(type) {
	case gotgbot.ChatMemberLeft, gotgbot.ChatMemberBanned:
		return false, nil
	case gotgbot.ChatMemberRestricted:
		return m.IsMember, nil
	default:
		return true, nil
	}
}

// HasJoinRequest reports wether the user has a join request pending for the channel saved in the db.
//
// - db: database client.
// - userId: id of user to check.
//
// NOTE: On db level error returns true and error.
func HasJoinRequest(db database.Database, chatId, userId int64) (bool, error) {
	user, err := db.GetUser(userId)
	if err != nil {
		return true, err
	}

	for _, c := range user.JoinRequests {
		if chatId == c {
			return true, nil
		}
	}

	return false, nil
}

// GetNotMemberOrRequest returns all channels the user is not a member of or has pending join request.
//
// - bot: bot client to make api requests.
// - db: database client.
// - f: fsub channels from config.
// - userId: id of user to check.
//
// Returns slice of channels user is not a member. Error returned will be any API call or DB query failure.
func GetNotMemberOrRequest(bot *gotgbot.Bot, db database.Database, f []model.FsubChannel, userId int64) ([]model.FsubChannel, error) {
	var (
		user      *model.User
		allErrors []error
		notJoined = make([]model.FsubChannel, len(f))
	)

	for _, c := range f {
		isMember, err := IsMember(bot, c.ID, userId)
		if err != nil {
			allErrors = append(allErrors, err) //TODO: wrap error with channel info
		}

		if isMember {
			continue
		}

		if user == nil {
			user, err = db.GetUser(userId)
			if err != nil {
				allErrors = append(allErrors, err)
				continue // or break to prevent further db queries?
			}
		}

		var hasJoinRequest bool

		for _, id := range user.JoinRequests {
			if id == c.ID {
				hasJoinRequest = true
				break
			}
		}

		if hasJoinRequest {
			continue
		}

		notJoined = append(notJoined, c)
	}

	return notJoined, errors.Join(allErrors...)
}
