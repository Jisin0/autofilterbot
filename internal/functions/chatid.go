package functions

const MaxChannelID int64 = -1000000000000

// ChatIdToMtproto converts a bot api channel id starting with -100 to an mtproto id.
func ChatIdToMtproto(id int64) int64 {
	if id < 0 {
		id = MaxChannelID - id
	}
	return id
}

// MtprotoToChatId converts a mtproto chat id to a bot api one.
// Should only be used for channel or supergroup ids.
func MtprotoToChatId(id int64) int64 {
	return MaxChannelID - id
}
