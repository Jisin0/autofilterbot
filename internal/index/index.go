package index

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/amarnathcjd/gogram/telegram"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Default API Credentials from tgx, generating your own and setting as env vars is recommended
const (
	DefaultAppID   = 21724
	DefaultAppHash = "3e0cb5efcd52300aec5994fdfc5bdc16"
)

// Operation handles and manages the index operation.
type Operation struct {
	mu sync.Mutex

	*model.Index

	db  database.Database
	log *zap.Logger
	bot *gotgbot.Bot

	// for accurate calculation of ETA

	startTime        time.Time // time at which this intance of operation was started/resumed
	startMessageID   int64     // msg id at which this intance of operation started/resumed
	mtprotoChannelID int64

	cancelFunc      context.CancelFunc
	completedSignal chan byte // notifies goroutines linked to the operation of completion
}

// NewOperation creates a new index operation and context to pass to *Operation.Run.
func (m *Manager) NewOperation(ctx context.Context, i *model.Index, db database.Database, log *zap.Logger, b *gotgbot.Bot) (context.Context, *Operation) {
	ctx2, cancel := context.WithCancel(ctx)
	return ctx2, &Operation{
		Index:           i,
		db:              db,
		log:             log,
		bot:             b,
		cancelFunc:      cancel,
		completedSignal: make(chan byte),
	}
}

const (
	defaultBatchSize = 200

	progressUpdateSeconds = 10 // number of seconds after which progress msg should be updated
)

// Run starts the index operation from the given CurrentMessageID until EndMessageID.
func (o *Operation) run(ctx context.Context) {
	// updates the progress msg to a basic start msg, also doubles as a check if the msg exists and can be edited
	startText := fmt.Sprintf("index from %d at %d to %d ...", o.StartMessageID, o.CurrentMessageID, o.EndMessageID)

	if o.CurrentMessageID == o.StartMessageID {
		startText = "Starting " + startText
	} else {
		startText = "Resuming " + startText
	}

	progressM, err := o.bot.SendMessage(o.ProgressMessageChatID, startText, &gotgbot.SendMessageOpts{})
	if err != nil {
		o.log.Error(err.Error(), zap.String("pid", o.ID), zap.Int64("message_id", progressM.MessageId), zap.Int64("chat_id", o.ProgressMessageChatID))
		o.bot.SendMessage(o.ProgressMessageChatID, fmt.Sprintf("🛑 Index Stopped: Unable to Update Progress Message: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{o.ResumeButton()}}},
		})

		return
	}

	//TODO: refactor error msg code

	var (
		appID   = DefaultAppID
		appHash = DefaultAppHash
	)

	if s := os.Getenv("APP_ID"); s != "" {
		id, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			o.log.Debug("index: failed to parse app id from environment", zap.Error(err), zap.String("val", s))
		} else {
			if s := os.Getenv("APP_HASH"); s != "" {
				appID = int(id)
				appHash = s
			} else {
				o.log.Warn("index: app id is set but app hash is empty, operation starting using deafult credentials")
			}
		}
	}

	c, err := telegram.NewClient(telegram.ClientConfig{
		AppID:         int32(appID),
		AppHash:       appHash,
		NoUpdates:     true,
		MemorySession: true,
		LogLevel:      telegram.LogError,
		DisableCache:  true,
	})
	if err != nil {
		o.log.Error(fmt.Sprintf("index: create client failed: %v", err), zap.String("pid", o.ID))
		o.bot.SendMessage(o.ProgressMessageChatID, fmt.Sprintf("🛑 Index Stopped: Unable to Create Client: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{o.ResumeButton()}}},
		})

		return
	}

	err = c.LoginBot(o.bot.Token)
	if err != nil {
		o.log.Error(fmt.Sprintf("index: login bot failed: %v", err), zap.String("pid", o.ID))
		o.bot.SendMessage(o.ProgressMessageChatID, fmt.Sprintf("🛑 Index Stopped: Unable to Login Bot: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{o.ResumeButton()}}},
		})

		return
	}

	_, err = c.GetMe()
	if err != nil {
		o.log.Error(fmt.Sprintf("index: getme failed: %v", err), zap.String("pid", o.ID))
		o.bot.SendMessage(o.ProgressMessageChatID, fmt.Sprintf("🛑 Index Stopped: Unable to Invoke Method: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{o.ResumeButton()}}},
		})

		return
	}

	inputChannel, err := getChat(c, o.ChannelID)
	if err != nil {
		o.log.Error(fmt.Sprintf("index: getchat failed: %v", err), zap.String("pid", o.ID), zap.Int64("tdlib_id", o.ChannelID))
		o.bot.SendMessage(o.ProgressMessageChatID, fmt.Sprintf("🛑 Index Stopped: Unable to Get Chat: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{o.ResumeButton()}}},
		})

		return
	}

	msgChan := make(chan []telegram.Message, defaultBatchSize) // allows to queue one full chunk, may cause inaccuracy in saved value in completed progress

	go o.MessageProcessor(ctx, msgChan)

	updateTicker := time.NewTicker(time.Second * time.Duration(progressUpdateSeconds))

	o.startTime = time.Now()
	o.startMessageID = o.CurrentMessageID
	o.mtprotoChannelID = inputChannel.ChannelID

	// updates progress msg and sync db, dettatched from index operation for real time updates
	// ticker may need to be adjusted in case of msg edit floods
	go func() {
		for {
			select {
			case <-o.completedSignal: // operation complete
				return
			case <-ctx.Done(): // user cancel
				return
			case <-updateTicker.C:
				o.pushToDB()

				progressBuilder := o.buildProgressMessage()

				progressBuilder.WriteString("\n<b>Index in Progress ⚡️</b>")

				_, _, err := progressM.EditText(o.bot, progressBuilder.String(), &gotgbot.EditMessageTextOpts{
					ParseMode: gotgbot.ParseModeHTML,
					ReplyMarkup: gotgbot.InlineKeyboardMarkup{
						InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
							{o.PauseButton(), o.CancelButton()},
						},
					},
					ChatId:    progressM.Chat.Id,
					MessageId: progressM.MessageId,
				})
				if err != nil {
					o.log.Debug(fmt.Sprintf("index: failed to update progress message: %v", err), zap.String("pid", o.ID), zap.Int64("message_id", progressM.MessageId))
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// operation paused either by the user or application quitting, not operation completion
			o.pushToDB()

			progressBuilder := o.buildProgressMessage()

			progressBuilder.WriteString("\n<b>Index Operation Paused ▶️</b>")

			_, _, err := progressM.EditText(o.bot, progressBuilder.String(), &gotgbot.EditMessageTextOpts{
				ParseMode: gotgbot.ParseModeHTML,
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{
					InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
						{o.ResumeButton(), o.ModifyButton(), o.CancelButton()},
					},
				},
			})
			if err != nil {
				o.log.Warn(fmt.Sprintf("index: failed to update paused message: %v", err), zap.String("pid", o.ID), zap.Int64("message_id", progressM.MessageId))
			}

			return
		default:
			// check if end reached
			if o.CurrentMessageID >= o.EndMessageID {
				o.completedSignal <- 1

				b := o.buildProgressMessage()
				b.WriteString("\n<b>Index Operation Completed 🎉</b>")

				_, _, err := progressM.EditText(o.bot, b.String(), &gotgbot.EditMessageTextOpts{
					ParseMode: gotgbot.ParseModeHTML,
				})
				if err != nil {
					o.log.Debug("index: failed to update progress to success message", zap.Error(err), zap.String("pid", o.ID))
				}

				err = o.db.DeleteOperation(o.ID)
				if err != nil {
					o.log.Warn(fmt.Sprintf("index: delete operation failed: %v", err), zap.String("pid", o.ID))
				}

				return
			}

			messageChunk := o.inputMessageSlice()

			rawMsgs, err := c.ChannelsGetMessages(inputChannel, messageChunk)
			if err != nil {
				s, ok, e := ParseMtProtoFloodwait(err)
				if e != nil {
					o.log.Error(
						fmt.Sprintf("index: parse floodwait error failed: %v", e),
						zap.String("pid", o.ID),
						zap.String("api_error", err.Error()),
						zap.Int64("channel_id", o.ChannelID),
					)

					continue
				}

				if !ok {
					o.log.Warn(
						fmt.Sprintf("index: getmessages failed: %v", err),
						zap.String("pid", o.ID),
						zap.Int64("current_msg_id", o.CurrentMessageID),
						zap.Int64("channel_id", o.ChannelID),
					)
					//does not continue(skip incrementing and retry) so process moves onto next batch
				}

				if s != 0 {
					time.Sleep(time.Second * time.Duration(s))
					continue
				}

				o.log.Error(fmt.Sprintf("index: getmessages failed: %v", err), zap.String("pid", o.ID), zap.Int64("channel_id", o.ChannelID))
				o.bot.SendMessage(o.ProgressMessageChatID, fmt.Sprintf("🛑 Index Stopped: Unable to Get messages: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{
					ParseMode:   gotgbot.ParseModeHTML,
					ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{o.ResumeButton()}}},
				})

				return
			}

			msgs := make([]telegram.Message, 0)

			switch m := rawMsgs.(type) {
			case *telegram.MessagesChannelMessages:
				msgs = m.Messages
			case *telegram.MessagesMessagesObj:
				msgs = m.Messages
			case *telegram.MessagesMessagesSlice:
				msgs = m.Messages
			case *telegram.MessagesMessagesNotModified:

			}

			msgChan <- msgs
			o.CurrentMessageID = int64(messageChunk[len(messageChunk)-1].(*telegram.InputMessageID).ID) // sets current id to that of last msg in chunk
		}
	}
}

// pushToDB updates the progress of the operation in the database. Errors are output to logger.
func (o *Operation) pushToDB() {
	update := map[string]interface{}{
		"current": o.CurrentMessageID,
		"saved":   o.Saved,
		"failed":  o.Failed,
	}

	_, err := o.db.UpdateIndexOperation(o.ID, update)
	if err != nil {
		o.log.Error(fmt.Sprintf("index: failed to update db values %v", err), zap.String("pid", o.ID))
	}
}

// getChat fetches a channel and it's access hash from it's botapi/tdlib id.
func getChat(client *telegram.Client, id int64) (*telegram.InputChannelObj, error) {
	rawChats, err := client.ChannelsGetChannels([]telegram.InputChannel{&telegram.InputChannelObj{ChannelID: TDLibChannelIDToPlain(id), AccessHash: 0}})
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	var chats []telegram.Chat

	switch c := rawChats.(type) {
	case *telegram.MessagesChatsObj:
		chats = c.Chats
	case *telegram.MessagesChatsSlice:
		chats = c.Chats
	default:
		return nil, errors.New("unknown chats type")
	}

	if len(chats) == 0 {
		return nil, errors.New("chats list is empty")
	}

	switch c := chats[0].(type) {
	case *telegram.Channel:
		return &telegram.InputChannelObj{ChannelID: c.ID, AccessHash: c.AccessHash}, nil
	case *telegram.ChatObj:
		return &telegram.InputChannelObj{ChannelID: c.ID, AccessHash: 0}, nil // probably should just skip regular chats
	case *telegram.ChannelForbidden:
		return nil, errors.New("channel forbidden")
	case *telegram.ChatEmpty:
		return nil, errors.New("chat empty")
	case *telegram.ChatForbidden:
		return nil, errors.New("chat forbidden")
	default:
		return nil, errors.New(fmt.Sprintf("unknown chat type: %T", c))
	}
}
