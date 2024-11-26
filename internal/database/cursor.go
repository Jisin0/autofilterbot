package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure *mongo.Cursor implements the interface.
var _ Cursor = (*mongo.Cursor)(nil)

// Cursor implements pagination to adapt for various databases.
type Cursor interface {
	// Next loads the next document from the cursor.
	Next(context.Context) bool
	// Decode unmarshals the current document into the value pointed to by v.
	// The value v must be a pointer to a struct or map.
	Decode(v interface{}) error
	// Close closes the cursor, releasing any resources associated with it.
	// It should be called after the cursor is no longer needed.
	Close(ctx context.Context) error
}
