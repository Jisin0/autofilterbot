package functions

const MaxChannelID int64 = -1000000000000

// ChatIdToMtproto converts a bot api channel id starting with -100 to an mtproto id.
func ChatIdToMtproto(id int64) int64 {
	if id < 0 {
		id = MaxChannelID - id
	}
	return id
}
