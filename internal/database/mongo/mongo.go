// Package mongo implements database.Database using mongodb.
package mongo

import (
	"context"
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Ensure *Client implements database.Database
var _ database.Database = (*Client)(nil)

// Client implements database.Database using mongodb
type Client struct {
	// userCollections stores data about users of the bot.
	userCollection *mongo.Collection
	// joinRequestCollection stores join requests sent by users.
	joinRequestsCollection *mongo.Collection
	// fileCollection stores all saved files.
	fileCollection *MultiCollection
	// configCollection stores settings configuration of the bot.
	configCollection *mongo.Collection
	// groupCollection contains data about group chats.
	groupCollection *mongo.Collection
	// Collection of long operations like index.
	opsCollection *mongo.Collection

	ctx    context.Context
	client *mongo.Client
	db     *mongo.Database
}

// NewClientOpts provides optional parameters to NewClient().
type NewClientOpts struct {
	// Name of the dabase within the cluster. Defaults to database.DefaultDatabaseName.
	DatabaseName string
	// Name of the collection where files are stored. Defaults to database.DefaultCollectionNameFiles.
	FilesCollectionName string
	// Additional database urls aside from the primary db, used to store files.
	AdditionalURLs []string
	// Index of the file collection to use for storage, defaults to 0. Can be updated from config panel.
	MultiCollectionIndex int
}

// NewClient creates a new client and connect to mongodb.
//
// - ctx: context that will be further used for every db query.
// - mongodbUri: primary database uri.
// - databaseName: name of database.
// - collectionName: name of file or media collection.
// - extraURLs: optional. Additional mongodb urls for storing files.
func NewClient(ctx context.Context, mongodbUri string, log *zap.Logger, opts ...NewClientOpts) (*Client, error) { //TODO: implement multi collection with updater
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
	if err != nil {
		return nil, err
	}

	var clientOpts NewClientOpts
	if len(opts) != 0 {
		clientOpts = opts[0]
	}

	databaseName := database.DefaultDatabaseName
	if clientOpts.DatabaseName != "" {
		databaseName = clientOpts.DatabaseName
	}

	collectionName := database.CollectionNameFiles
	if clientOpts.FilesCollectionName != "" {
		collectionName = clientOpts.FilesCollectionName
	}

	dataBase := mongoClient.Database(databaseName)
	primaryFileCollection := dataBase.Collection(collectionName)

	fileCollections := []*mongo.Collection{primaryFileCollection}

	for i, url := range clientOpts.AdditionalURLs {
		c, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			log.Warn("mongo: newclient: failed to connect to additional database", zap.Int("num", i+1))
			continue
		}

		fileCollections = append(fileCollections, c.Database(databaseName).Collection(collectionName))
	}

	fileCollection := NewMultiCollection(fileCollections, clientOpts.MultiCollectionIndex, log)

	primaryFileCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{Key: "file_name", Value: "text"}, {Key: "time", Value: 1}}})

	client := &Client{
		ctx:                    ctx,
		client:                 mongoClient,
		db:                     dataBase,
		userCollection:         dataBase.Collection(database.CollectionNameUsers),
		fileCollection:         fileCollection,
		configCollection:       dataBase.Collection(database.CollectionNameConfigs),
		groupCollection:        dataBase.Collection(database.CollectionNameGroups),
		opsCollection:          dataBase.Collection(database.CollectionNameOperations),
		joinRequestsCollection: dataBase.Collection(database.CollectionNameJoinRequests),
	}

	return client, nil
}

func (c *Client) Shutdown() error {
	return c.client.Disconnect(context.Background()) // main ctx may already have been cancelled when this is called
}

// fileCounts generates a better visual list of file collection counts, when implementing fmt.Stringer, used for stats.
type fileCounts []int64

func (f fileCounts) String() string {
	if len(f) == 0 {
		return "No Collections Found"
	}

	if len(f) == 1 {
		return fmt.Sprint(f[0])
	}

	var s string

	for i, n := range f {
		s += fmt.Sprintf("\n├┄┄Collection %d: %d", i, n)
	}

	return s
}

func (c *Client) Stats() (*model.Stats, error) {
	users, err := c.userCollection.EstimatedDocumentCount(c.ctx)
	if err != nil {
		return nil, err
	}

	groups, err := c.groupCollection.EstimatedDocumentCount(c.ctx)
	if err != nil {
		return nil, err
	}

	var files []int64

	for _, coll := range c.fileCollection.allCollections {
		n, err := coll.EstimatedDocumentCount(c.ctx)
		if err != nil {
			return nil, err
		}

		files = append(files, n)
	}

	return &model.Stats{
		Users:  users,
		Groups: groups,
		Files:  fileCounts(files),
	}, nil
}

func (c *Client) GetName() string {
	return "MongoDB Atlas"
}

// UpdateStorageCollection updates the storage collection to given index.
func (c *Client) UpdateStorageCollection(index int) error {
	return c.fileCollection.SetStorageCollection(index)
}

// RunCollectionUpdater is a background job that ensures storageCollection is set to the collection with least documents stored.
//
// WARNING: The document count of the collection does not essentially represent the storage usage of the database but the logic depends on the assumption that files will be by far the heaviest collection.
func (c *Client) RunCollectionUpdater(ctx context.Context, log *zap.Logger) {
	go c.fileCollection.RunCollectionUpdater(ctx, log)
}
