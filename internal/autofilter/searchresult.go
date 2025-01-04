package autofilter

import "github.com/Jisin0/autofilterbot/internal/model"

// SearchResult holds the result of a search query.
type SearchResult struct {
	// Query is the sanitized message that was searched for.
	Query string `json:"query"`
	// FromUser is the id of the user who initiated it.
	FromUser int64 `json:"from_user"`
	// Id of the chat where the query was started.
	ChatID int64 `json:"chat_id,omitempty"`
	// Files are the files fetched from the datbase.
	Files [][]model.File `json:"files,omitempty"`
}
