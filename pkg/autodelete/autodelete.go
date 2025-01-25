package autodelete

import (
	"context"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const (
	dbFileName = "./autodelete.sqlite"
)

// Manager manages autodelete tasks.
type Manager struct {
	// Bot client to delete messages with.
	Bot *gotgbot.Bot
	// Database which stores messages.
	DB *sqlx.DB

	// duration isn't being added here for better runtime control.
}

// NewManager creates a new Manager from given bot.
func NewManager(bot *gotgbot.Bot) (*Manager, error) {
	db, err := sqlx.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS autodelete (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        chat_id INTEGER,
        message_id INTEGER,
        expiry_time DATETIME,
        UNIQUE(chat_id, message_id)
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Bot: bot,
		DB:  db,
	}, nil
}

const insertQuery = `INSERT INTO autodelete (chat_id, message_id, expiry_time) VALUES (:chat_id, :message_id, :expiry_time)
	ON CONFLICT(chat_id, message_id) DO UPDATE SET expiry_time=excluded.expiry_time;`

// Save adds a message to the autodelete database which will be deleted after duration.
func (m *Manager) Save(chatId, messageId int64, duration time.Duration) error {
	data := MessageData{
		ChatId:     chatId,
		MessageId:  messageId,
		ExpiryTime: time.Now().Add(duration),
	}

	_, err := m.DB.NamedExec(insertQuery, data)

	return err
}

// SaveMessage adds a message to the autodelete database which will be deleted after duration.
func (m *Manager) SaveMessage(msg *gotgbot.Message, duration time.Duration) error {
	data := MessageData{
		ChatId:     msg.Chat.Id,
		MessageId:  msg.MessageId,
		ExpiryTime: time.Now().Add(duration),
	}

	_, err := m.DB.NamedExec(insertQuery, data)

	return err
}

const deleteQuery = `DELETE FROM autodelete WHERE chat_id = ? AND message_id = ?`

// Remove removes a message from the database by its chatId & messageId.
func (m *Manager) Remove(chatId, messageId int64) error {
	_, err := m.DB.Exec(deleteQuery, chatId, messageId)
	return err
}

const selectQuery = `SELECT chat_id, message_id, expiry_time FROM autodelete WHERE expiry_time <= ?`

// Run starts the autodelete system which deletes expired messages every minute.
func (m *Manager) Run(ctx context.Context, log *zap.Logger) {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			var result []MessageData

			err := m.DB.Select(&result, selectQuery, time.Now())
			if err != nil {
				log.Warn("autodelete select query failed", zap.Error(err))
				break
			}

			for _, r := range result {
				_, err := m.Bot.DeleteMessage(r.ChatId, r.MessageId, nil)
				if err != nil {
					log.Info("autodelete message failed",
						zap.Int64("chat_id", r.ChatId),
						zap.Int64("message_id", r.MessageId),
						zap.Error(err),
					)
				}

				err = m.Remove(r.ChatId, r.MessageId)
				if err != nil {
					log.Warn("autodelete remove entry failed", zap.Error(err))
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// MessageData is a single row or entry in the autodelete database and hold data abou message to delete.
type MessageData struct {
	// Chat id where message is posted.
	ChatId int64 `db:"chat_id"`
	// Unique id of the message in the chat.
	MessageId int64 `db:"message_id"`
	// Time at which message will expire.
	ExpiryTime time.Time `db:"expiry_time"`
}
