package autofilter

// SearchResult holds the result of a search query.
type SearchResult struct {
	// Unique id used to identify the query.
	UniqueId string
	// Query is the sanitized message that was searched for.
	Query string `json:"query"`
	// FromUser is the id of the user who initiated it.
	FromUser int64 `json:"from_user"`
	// Id of the chat where the query was started.
	ChatID int64 `json:"chat_id,omitempty"`
	// Files are the files fetched from the datbase.
	Files []Files `json:"files,omitempty"`
}
