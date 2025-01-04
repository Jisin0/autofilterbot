// Package mongo implements database.Database using mongodb.
package mongo

import (
	"context"

	"github.com/Jisin0/autofilterbot/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ensure *Client implements database.Database
var _ database.Database = (*Client)(nil)

// Client implements database.Database using mongodb
type Client struct {
	// userCollections stores data about users of the bot.
	userCollection *mongo.Collection
	// fileCollection stores all saved files
	fileCollection *mongo.Collection
	// configCollection stores settings configuration of the bot.
	configCollection *mongo.Collection
	// groupCollection contains data about group chats.
	groupCollection *mongo.Collection

	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
}

// NewClient creates a new client and connect to mongodb.
//
// - ctx: context that will be further used for every db query.
// - mongodbUri: primary database uri.
// - databaseName: name of database.
// - collectionName: name of file or media collection.
// - extraURLs: optional. Additional mongodb urls for storing files.
func NewClient(ctx context.Context, mongodbUri, databaseName, collectionName string, extraURLs ...string) (*Client, error) { //TODO: implement multi collection with updater
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
	if err != nil {
		return nil, err
	}

	if databaseName == "" {
		databaseName = database.DefaultDatabaseName
	}

	if collectionName == "" {
		collectionName = database.CollectionNameFiles
	}

	dataBase := mongoClient.Database(databaseName)

	fcol := dataBase.Collection(database.CollectionNameFiles)
	fcol.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "file_name", Value: "text"}, {Key: "time", Value: 1}}})

	client := &Client{
		ctx:              ctx,
		client:           mongoClient,
		db:               dataBase,
		userCollection:   dataBase.Collection(database.CollectionNameUsers),
		fileCollection:   fcol,
		configCollection: dataBase.Collection(database.CollectionNameConfigs),
		groupCollection:  dataBase.Collection(database.CollectionNameGroups),
	}

	return client, nil
}

func (c *Client) Shutdown() error {
	return c.client.Disconnect(context.Background()) // main ctx may already have been cancelled when this is called
}
