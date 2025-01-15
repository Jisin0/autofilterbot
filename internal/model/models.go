// Package models contains types shared by various packages.
package model

// FsubChannel is a single force sub channel.
type FsubChannel struct {
	// Telegram id of the channel.
	ID int64 `json:"id" bson:"id"`
	// Name or title of the channel.
	Title string `json:"title" bson:"title"`
	// Invite link for the channel.
	InviteLink string `json:"link" bson:"link"`
}

// Stats are database statistics.
type Stats struct {
	Users  int64
	Groups int64
	Files  interface{} // allows for flexibility, custom types must implement fmt.Stringer
}
