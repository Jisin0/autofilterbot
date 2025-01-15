// Package database holds interfaces and types used for mongo or couchbase storage.
package database

import (
	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/internal/model"
)

const (
	CollectionNameUsers      = "Users"
	CollectionNameFiles      = "Files"
	CollectionNameConfigs    = "Configs"
	CollectionNameOperations = "Operations"
	CollectionNameGroups     = "Groups"

	DefaultDatabaseName = "AutoFilterBot"
)

type Database interface {
	// Shutdown gracefully closes the database.
	Shutdown() error

	// SaveUser saves the id of a user to the database if it does not exist.
	SaveUser(userId int64) error
	// GetUser gets a user from the database using their id.
	GetUser(userId int64) (*model.User, error)
	// DeleteUser deletes a user from the database. This could be because the user has blocked the bot.
	DeleteUser(userId int64) error
	// SaveUserJoinRequest saves the chat id to which a user has sent a join request.
	// The join request is not saved if the user is not saved in the database.
	SaveUserJoinRequest(userId, chatId int64) error
	// DeletUserJoinRequest deletes the chat from the join requests list.
	DeleteUserJoinRequest(userId, chatId int64) error
	// GetUsers returns a cursor to loop through all saved users.
	GetAllUsers() (Cursor, error)

	// SaveFile saves a file to the database and returns a FileAlreadyExistsError if the file already exists.
	// The file can be a duplicate if it has the same file_id or file_name-file_size combination.
	SaveFile(f *model.File) error
	// SaveFiles saves multiple files to the database and returns a list of errors.
	SaveFiles(files ...*model.File) []error
	// GetFile fetches a file from the database using its unique_id.
	GetFile(fileId string) (*model.File, error)
	// DeleteFile deletes a file from the database using its unique_id.
	DeleteFile(fileId string) error
	// SearchFiles searches for files in the database by their name. The query should be sanitized first.
	SearchFiles(query string) (Cursor, error)

	// SaveGroup inserts a group id into the database to keep track of them.
	SaveGroup(id int64) error

	// GetConfig fetches the bot configs from the database.
	GetConfig(botId int64) (*config.Config, error)
	// UpdateConfig updates a single element of config.
	UpdateConfig(botId int64, key string, value interface{}) error
	// SaveConfig saves the config struct. Useful for importing configs.
	SaveConfig(botId int64, data *config.Config) error
	// ResetConfig removes a config field, resetting it to it's default value.
	ResetConfig(botId int64, key string) error

	// Stats gets the database usage statistics.
	Stats() (*model.Stats, error)
}
