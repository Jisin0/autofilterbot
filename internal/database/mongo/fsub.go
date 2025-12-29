package mongo

import "go.mongodb.org/mongo-driver/bson"

func (c *Client) SaveUserJoinRequest(userId, chatId int64) error {
	// if user isn't in db and this is executed no error is returned, the UpdateResult will have MatchedCount: 0.
	//TODO: save join request and add flag to indicate bot can't msg user to prevent edge cases
	_, err := c.userCollection.UpdateOne(
		c.ctx,
		idFilter(userId),
		bson.D{{Key: "$addToSet", Value: bson.D{{Key: "join_requests", Value: chatId}}}},
	)
	return err
}

func (c *Client) DeleteUserJoinRequest(userId, chatId int64) error {
	_, err := c.userCollection.UpdateOne(
		c.ctx,
		idFilter(userId),
		bson.D{{Key: "$pull", Value: bson.D{{Key: "join_requests", Value: chatId}}}},
	)
	return err
}
