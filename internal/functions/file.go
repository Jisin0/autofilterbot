package functions

import (
	"errors"

	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

var ErrFileNotFound = errors.New("no media was found in the message")

// FileFromMessage extracts data about a file from the message.
func FileFromMessage(m *gotgbot.Message) *model.File {
	if m == nil {
		return nil
	}

	var (
		fileSize                             int64
		fileId, uniqueId, fileName, fileType string
	)

	switch {
	case m.Document != nil:
		fileId = m.Document.FileId
		uniqueId = m.Document.FileUniqueId
		fileName = m.Document.FileName
		fileSize = m.Document.FileSize
		fileType = model.FileTypeDocument
	case m.Video != nil:
		fileId = m.Video.FileId
		uniqueId = m.Video.FileUniqueId
		fileName = m.Video.FileName
		fileSize = m.Document.FileSize
		fileType = model.FileTypeVideo
	case m.Audio != nil:
		fileId = m.Audio.FileId
		uniqueId = m.Audio.FileUniqueId
		fileName = m.Audio.FileName
		fileSize = m.Audio.FileSize
		fileType = model.FileTypeAudio
	default:
		return nil
	}

	fileName = RemoveSymbols(RemoveExtension(fileName))

	return &model.File{
		UniqueId:    uniqueId,
		FileId:      fileId,
		FileName:    fileName,
		FileType:    fileType,
		FileSize:    fileSize,
		Time:        m.Date,
		ChatId:      m.Chat.Id,
		MessageLink: m.GetLink(),
	}
}
