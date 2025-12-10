package mongo

import (
	"context"
	"sync"

	"github.com/Jisin0/autofilterbot/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Ensure *MultiCursor implements database.Cursor.
var _ database.Cursor = (*MultiCursor)(nil)

// MultiCursor orchestrates mongodb queries to multiple collections as a single virtual colllection implementing database.Cursor.
type MultiCursor struct {
	mu  sync.Mutex
	log *zap.Logger

	// current cursor that is being queried and data being fetched.
	currentCursor *mongo.Cursor
	// collections that have not been queried yet.
	remainingCollections []*mongo.Collection
	// filter to query collections.
	filter interface{}
	// find options
	opts []*options.FindOptions
}

// Next loads the next document into the current cursor or quries the next collection if available. It returns false if all collections were exhausted.
func (c *MultiCursor) Next(ctx context.Context) bool {
	if c.currentCursor == nil { // should never happen in a perfect world
		c.log.Warn("multicursor: next: current cursor is nil", zap.Any("query", c.filter))
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.currentCursor.Next(ctx) {
		return true
	}

	if len(c.remainingCollections) == 0 {
		return false
	}

	var (
		res *mongo.Cursor
		err error
	)

	for i, col := range c.remainingCollections {
		res, err = col.Find(ctx, c.filter, c.opts...) // should this ctx be used? maybe pass ctx to MultiCursor from Find and use the same
		if err != nil {
			c.log.Debug("multicursor: next: find operation failed", zap.Error(err))
			continue
		}

		if res.Next(ctx) { // If cursor provides a document update currentCursor and cut remainingCollections
			if err := c.currentCursor.Close(ctx); err != nil { // close current cursor
				c.log.Error("multicursor: next: failed to close current cursor", zap.Error(err))
			}

			c.currentCursor = res

			if len(c.remainingCollections) > i+1 {
				c.remainingCollections = c.remainingCollections[i+1:]
			} else if len(c.remainingCollections) == i+1 {
				c.remainingCollections = nil
			}

			return true
		}
	}

	// all collections are exhausted so empty remainingCollections just in case and return false
	c.remainingCollections = nil

	return false
}

// Decode unmarshals the current document into the value pointed to by v.
// The value v must be a pointer to a struct or map.
func (c *MultiCursor) Decode(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.currentCursor.Decode(v)
}

// Close closes the cursor, releasing any resources associated with it.
// It should be called after the cursor is no longer needed.
func (c *MultiCursor) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.currentCursor.Close(ctx)
}
