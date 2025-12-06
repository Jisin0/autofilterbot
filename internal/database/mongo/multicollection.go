package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Jisin0/autofilterbot/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	collectionUpdateDuration = 10 * time.Minute // duration every collection update should be run
)

// MultiCollection wraps a number of collections to create a virtual collection that can be effectively queried like a regular mongo collection.
type MultiCollection struct {
	// storageCollection is the current collection that will save new documents.
	// It is the collection with least documents and is periodically updated by a background job.
	storageCollection *mongo.Collection
	// Index number of current storage collection.
	storageCollectionIndex int
	// all file storage collections with collection from primary database at index 0.
	allCollections []*mongo.Collection

	log *zap.Logger
}

// NewMultiCollection creates a new multi collection and sets the collection at given index as current storage.
func NewMultiCollection(allCollections []*mongo.Collection, index int, log *zap.Logger) *MultiCollection {
	return &MultiCollection{
		storageCollection:      allCollections[index],
		storageCollectionIndex: index,
		allCollections:         allCollections,
		log:                    log,
	}
}

// InsertOne inserts a single document to the current storage collection.
func (c *MultiCollection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.storageCollection.InsertOne(ctx, document)
}

// Find executes a find command and returns a Cursor over the matching documents in the virtual collection.
func (c *MultiCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (database.Cursor, error) {
	cursor := &MultiCursor{
		log:    c.log,
		filter: filter,
		opts:   opts,
	}

	var (
		res *mongo.Cursor
		err error
	)

	for i, col := range c.allCollections {
		res, err = col.Find(ctx, filter, opts...)
		if err != nil { // this does not mean no documents are found. If filter was not matched, an empty cursor will be returned.
			continue
		}

		cursor.currentCursor = res

		if len(c.allCollections) > 1 {
			cursor.remainingCollections = c.allCollections[i+1 : len(c.allCollections)-1]
		}
	}

	return cursor, err
}

// FindOne finds a single document in any collection that matches given filter.
func (c *MultiCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	var r *mongo.SingleResult

	for _, col := range c.allCollections {
		r = col.FindOne(ctx, filter, opts...)
		if r.Err() == nil {
			return r
		}
	}

	return r
}

// DeleteMany deletes all documents matching given filter in any collection.
func (c *MultiCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	var (
		result    mongo.DeleteResult
		allErrors []error
	)

	for _, col := range c.allCollections {
		r, err := col.DeleteMany(ctx, filter, opts...)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		result.DeletedCount += r.DeletedCount
	}

	return &result, errors.Join(allErrors...)
}

// DeleteOne deletes the first document matching the filter in any collection.
// If the filter does not match any documents, the operation will succeed and a DeleteResult with a DeletedCount of 0 will be returned.
func (c *MultiCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	for i, col := range c.allCollections {
		res, err := col.DeleteOne(ctx, filter, opts...)
		if err != nil {
			c.log.Error("multicollection: deleteone failed", zap.Int("index", i), zap.Error(err))
			continue
		}

		if res.DeletedCount > 0 {
			return res, err
		}
	}

	return &mongo.DeleteResult{}, nil
}

// UpdateMany updates all documents that match filter in every collection.
// If the filter does not match any documents, the operation will succeed and an UpdateResult with a MatchedCount of 0 will be returned.
// The result will not contain UpsertedId field.
func (c *MultiCollection) UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	var (
		result    mongo.UpdateResult
		allErrors []error
	)

	for _, col := range c.allCollections {
		res, err := col.UpdateMany(ctx, filter, update, opts...)
		if err != nil {
			allErrors = append(allErrors, err)
		}

		if res != nil {
			result.MatchedCount += res.MatchedCount
			result.ModifiedCount += res.ModifiedCount
			result.UpsertedCount += res.UpsertedCount
		}
	}

	return &result, errors.Join(allErrors...)
}

// UpdateOne updates the first document matching the filter in any collection.
// If the filter does not match any documents, the operation will succeed and a UpdateResult with a ModifiedCount of 0 will be returned.
func (c *MultiCollection) UpdateOne(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	for i, col := range c.allCollections {
		res, err := col.UpdateOne(ctx, filter, update, opts...)
		if err != nil {
			c.log.Error("multicollection updateone failed", zap.Int("index", i), zap.Error(err))
			continue
		}

		if res.ModifiedCount > 0 {
			return res, err
		}
	}

	return &mongo.UpdateResult{}, nil
}

// EstimatedDocumentCount executes a count command and returns an estimate of the total number of documents in all collections using collection metadata.
//
// An error in any collectino will end with the accumulated total and error being returned immediately.
func (c *MultiCollection) EstimatedDocumentCount(ctx context.Context, opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	var total int64

	for _, col := range c.allCollections {
		n, err := col.EstimatedDocumentCount(ctx, opts...)
		if err != nil {
			return total, err
		}

		total += n
	}

	return total, nil
}

// RunCollectionUpdater is a background job that ensures storageCollection is set to the collection with least documents stored.
//
// WARNING: The document count of the collection does not essentially represent the storage usage of the database but the logic depends on the assumption that files will be by far the heaviest collection.
func (c *MultiCollection) RunCollectionUpdater(ctx context.Context, log *zap.Logger) {
	if len(c.allCollections) == 1 { // hopefully isnt 0 lol
		return
	}

	log.Debug("mongo collection updater job started")

	ticker := time.NewTicker(collectionUpdateDuration)

	for {
		select {
		case <-ticker.C:
			currentCount, err := c.storageCollection.EstimatedDocumentCount(ctx)
			if err != nil {
				log.Error("failed to get document count for current storage collection", zap.Error(err))
				continue
			}

			var (
				smallestDocumentCount      = currentCount
				smallestDocumentCollection = c.storageCollection

				i   int
				col *mongo.Collection
			)

			for i, col = range c.allCollections {
				if col == c.storageCollection { // skip current storage collection
					continue
				}

				count, err := col.EstimatedDocumentCount(ctx)
				if err != nil {
					log.Error("failed to get document count for collection", zap.Int("index", i), zap.Error(err))
					continue
				}

				if count < smallestDocumentCount {
					smallestDocumentCount = count
					smallestDocumentCollection = col
				}
			}

			// If smallest collection is different from current then update
			if smallestDocumentCollection != c.storageCollection {
				c.storageCollection = smallestDocumentCollection
				c.storageCollectionIndex = i

				log.Debug("multicollection: updated storage collection", zap.Int("index", i))
			}
		case <-ctx.Done():
			return
		}
	}
}

// SetStorageCollection sets the collection with given index for storing new files.
func (c *MultiCollection) SetStorageCollection(index int) error {
	if len(c.allCollections) <= index {
		return fmt.Errorf("multicolllection: setstorage: index %d out of range with length %d", index, len(c.allCollections))
	}

	c.storageCollection = c.allCollections[index]
	c.storageCollectionIndex = index

	return nil
}
