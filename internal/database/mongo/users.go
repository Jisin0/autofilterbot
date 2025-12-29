package mongo

import (
	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SaveUser creates a new document in the user collection with the user id.
func (c *Client) SaveUser(userId int64) error {
	_, err := c.userCollection.InsertOne(c.ctx, model.User{UserId: userId})
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	return nil
}

// GetUser fetches a user from the database by id.
func (c *Client) GetUser(userId int64) (*model.User, error) {
	var u model.User

	res := c.userCollection.FindOne(c.ctx, idFilter(userId))
	if err := res.Err(); err != nil {
		return nil, err
	}

	err := res.Decode(&u)

	return &u, err
}

// DeleteUser deletes a user by their id.
func (c *Client) DeleteUser(userId int64) error {
	_, err := c.userCollection.DeleteOne(c.ctx, idFilter(userId))
	return err
}

// GetAllUsers return a cursor to loop over all users.
func (c *Client) GetAllUsers() (database.Cursor, error) {
	return c.userCollection.Find(c.ctx, bson.M{})
}

// idFilter creates a basic bson filter to find documents with matching _id.
func idFilter(id interface{}) bson.D {
	return bson.D{{Key: "_id", Value: id}}
}
