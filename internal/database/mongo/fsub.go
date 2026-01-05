package mongo

import (
	"github.com/Jisin0/autofilterbot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var upsert = true

func (c *Client) SaveUserJoinRequest(userId, chatId int64) error {
	_, err := c.joinRequestsCollection.UpdateOne(
		c.ctx,
		idFilter(userId),
		bson.D{{Key: "$addToSet", Value: bson.D{{Key: "join_requests", Value: chatId}}}},
		&options.UpdateOptions{Upsert: &upsert},
	)
	return err
}

func (c *Client) DeleteUserJoinRequest(userId, chatId int64) error {
	_, err := c.joinRequestsCollection.UpdateOne(
		c.ctx,
		idFilter(userId),
		bson.D{{Key: "$pull", Value: bson.D{{Key: "join_requests", Value: chatId}}}},
		&options.UpdateOptions{Upsert: &upsert},
	)
	return err
}

// GetUser fetches a user's join requests from the database by id.
func (c *Client) GetUserJoinRequests(userId int64) (*model.User, error) {
	var u model.User

	res := c.joinRequestsCollection.FindOne(c.ctx, idFilter(userId))
	if err := res.Err(); err != nil {
		return nil, err
	}

	err := res.Decode(&u)

	return &u, err
}
