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
func (f FsubChannel) IsMember(bot *gotgbot.Bot, userId int64) (bool, error) {
	member, err := bot.GetChatMember(f.ID, userId, nil)
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
func (f FsubChannel) HasJoinRequest(db database.Database, userId int64) (bool, error) {
	user, err := db.GetUser(userId)
	if err != nil {
		return true, err
	}

	for _, c := range user.JoinRequests {
		if f.ID == c {
			return true, nil
		}
	}

	return false, nil
}

// FsubChannelArray is a list of all fsub channels.
type FsubChannelArray struct {
	DB   database.Database
	List []FsubChannel
}

// GetNotMemberOrRequest returns all channels the user is not a member of or has pending join request.
//
// - bot: bot client to make api requests.
// - db: database client.
// - userId: id of user to check.
func (f FsubChannelArray) GetNotMemberOrRequest(bot *gotgbot.Bot, userId int64) ([]FsubChannel, error) {
	var (
		l         = f.List
		user      *model.User
		allErrors []error
		notJoined = make([]FsubChannel, len(l))
	)

	for _, c := range l {
		isMember, err := c.IsMember(bot, userId)
		if err != nil {
			allErrors = append(allErrors, err) //TODO: wrap error with channel info
		}

		if isMember {
			continue
		}

		if user == nil {
			user, err = f.DB.GetUser(userId)
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
