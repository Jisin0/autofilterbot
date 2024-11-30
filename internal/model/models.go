// Package models contains types shared by various packages.
package model

// FsubChannel is a single force sub channel.
type FsubChannel struct {
	// Telegram id of the channel.
	ID int64 `json:"id"`
	// Name or title of the channel.
	Title string `json:"title"`
	// Invite link for the channel.
	InviteLink string `json:"link"`
}
