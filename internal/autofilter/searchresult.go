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

// SelectFile sets the IsSelected field of the file on given page.
func (r *SearchResult) SelectFile(pageIndex int, fileUniqueId string) bool {
	if pageIndex >= len(r.Files) {
		return false
	}

	for i, f := range r.Files[pageIndex] {
		if f.UniqueId == fileUniqueId {
			r.Files[pageIndex][i].IsSelected = !f.IsSelected

			return true
		}
	}

	return false
}
