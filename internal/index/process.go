package index

import (
	"context"
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/fileid"
	"github.com/amarnathcjd/gogram/telegram"
	"go.uber.org/zap"
)

func (o *Operation) MessageProcessor(ctx context.Context, c chan []telegram.Message) {
	for {
		select {
		case msgs := <-c:
			o.log.Debug("index: msgs to save received", zap.String("pid", o.ID), zap.Int("length", len(msgs)))

			for _, m := range msgs {
				msg, ok := m.(*telegram.MessageObj)
				if !ok {
					o.log.Debug("index: unspported msg type", zap.String("pid", o.ID), zap.String("type", fmt.Sprintf("%T", m)))
					o.Failed++
					continue
				}

				if msg.Media == nil {
					o.log.Debug("index: msg has no media", zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID))
					o.Failed++
					continue
				}

				media, ok := msg.Media.(*telegram.MessageMediaDocument)
				if !ok {
					o.log.Debug("index: unsupported media type", zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID), zap.String("type", fmt.Sprintf("%T", msg.Media)))
					o.Failed++
					continue
				}

				doc, ok := media.Document.(*telegram.DocumentObj)
				if !ok {
					o.log.Debug("index: document is empty", zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID), zap.String("type", fmt.Sprintf("%T", media.Document)))
					o.Failed++
					continue
				}

				var (
					fileType            = model.FileTypeDocument
					fileIDType          = fileid.Document
					fileName            string
					unsupportedDocument bool
				)

				for _, attr := range doc.Attributes {
					switch a := attr.(type) {
					case *telegram.DocumentAttributeAnimated, *telegram.DocumentAttributeHasStickers, *telegram.DocumentAttributeImageSize, *telegram.DocumentAttributeSticker:
						o.log.Debug("unsupported document type", zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID), zap.Any("attr", a))
						unsupportedDocument = true
					case *telegram.DocumentAttributeAudio:
						if a.Voice {
							fileType = model.FileTypeVoice
							fileIDType = fileid.Voice
						} else {
							fileType = model.FileTypeAudio
							fileIDType = fileid.Audio
						}
					case *telegram.DocumentAttributeVideo:
						fileType = model.FileTypeVideo
						fileIDType = fileid.Video
					case *telegram.DocumentAttributeFilename:
						fileName = a.FileName
					}
				}

				if unsupportedDocument {
					o.Failed++
					continue
				}

				if fileName == "" {
					o.log.Debug("filename attribute not found", zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID))
					o.Failed++
					continue
				}

				f := fileid.FileID{
					Type:          fileIDType,
					DC:            int(doc.DcID),
					ID:            doc.ID,
					AccessHash:    doc.AccessHash,
					FileReference: doc.FileReference,
				}

				fileID, err := fileid.EncodeFileID(f)
				if err != nil {
					o.log.Warn("encode file id failed", zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID), zap.Any("file", f))
					o.Failed++
					continue
				}

				file := model.File{
					UniqueId: functions.RandString(15),
					FileId:   fileID,
					FileName: fileName,
					FileType: fileType,
					FileSize: int64(doc.Size),
					Time:     int64(msg.Date),
				}

				err = o.db.SaveFile(&file)
				if err != nil {
					if _, ok := err.(database.FileAlreadyExistsError); ok {
						o.log.Debug("index: duplicate file skipped", zap.String("pid", o.ID), zap.String("file_name", file.FileName))
					} else {
						o.log.Warn("index: save file failed", zap.Error(err), zap.String("pid", o.ID), zap.Int32("msg_id", msg.ID))
					}

					o.Failed++

					continue
				}

				o.Saved++
			}
		case <-ctx.Done():
			return
		case <-o.completedSignal:
			return
		}
	}
}
