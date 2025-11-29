package mongo

import (
	"errors"

	"github.com/Jisin0/autofilterbot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

//NOTE: error checking is much looser in index db funcs

// NewIndexOperation inserts a new index operation into the database.
func (c *Client) NewIndexOperation(i *model.Index) error {
	_, err := c.opsCollection.InsertOne(c.ctx, i)
	return err
}

// UpdateIndexOperation updates an index operation.
// Returns a bool indication wether a match was found and errors.
func (c *Client) UpdateIndexOperation(pid string, vals map[string]interface{}) (bool, error) {
	r, err := c.opsCollection.UpdateOne(c.ctx, idFilter(pid), bson.M{"$set": bson.M(vals)})

	var ok bool

	if r != nil {
		ok = r.MatchedCount != 0
	}

	return ok, err
}

// GetIndexOperation fetches an index operation by it's id.
func (c *Client) GetIndexOperation(pid string) (*model.Index, error) {
	res := c.opsCollection.FindOne(c.ctx, idFilter(pid))
	if err := res.Err(); err != nil {
		return nil, err
	}

	var i model.Index

	err := res.Decode(&i)

	return &i, err
}

// GetAllIndexOperations fetches all active index operations.
func (c *Client) GetActiveIndexOperations() ([]*model.Index, error) {
	cursor, err := c.opsCollection.Find(c.ctx, bson.M{"is_paused": false})
	if err != nil {
		return nil, err
	}

	ops := make([]*model.Index, 0)
	errs := make([]error, 0)

	for cursor.Next(c.ctx) {
		var i model.Index

		e := cursor.Decode(&i)
		if e != nil {
			errs = append(errs, e)
			continue
		}

		ops = append(ops, &i)
	}

	return ops, errors.Join(errs...)
}

// DeleteOperation deletes an active operation by id.
func (c *Client) DeleteOperation(pid string) error {
	_, err := c.opsCollection.DeleteOne(c.ctx, idFilter(pid))
	return err
}
