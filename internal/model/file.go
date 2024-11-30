package model

// File is a single file stored in the database.
type File struct {
	// Unique id of the file used to save, fetch and send files.
	FileId string `json:"file_id"`
	// Name of the file including extension
	FileName string `json:"file_name"`
	// Type of file either document/video/audio.
	FileType string `json:"file_type"`
	// Size of the file in bytes.
	FileSize int64 `json:"file_size"`
	// Unix timestamp of time when file was saved.
	Time int64 `json:"time,omitempty"`
	// Id of the chat/channel where the file is posted.
	ChatId int64 `json:"chat_id,omitempty"`
	// Link to the original message containing the file.
	MessageLink string `json:"file_link,omitempty"`
}
