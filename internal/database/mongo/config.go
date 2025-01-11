package mongo

import (
	"errors"

	"github.com/Jisin0/autofilterbot/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	boolTrue = true
)

func (c *Client) GetConfig(botId int64) (*config.Config, error) {
	r := &config.Config{}

	res := c.configCollection.FindOne(c.ctx, idFilter(botId))
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return r, nil
		}

		return r, err
	}

	if err := res.Decode(r); err != nil {
		return r, err
	}

	return r, nil
}

func (c *Client) UpdateConfig(botId int64, key string, value interface{}) error {
	_, err := c.configCollection.UpdateOne(c.ctx, idFilter(botId), bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}, &options.UpdateOptions{Upsert: &boolTrue})
	return err
}

func (c *Client) SaveConfig(botId int64, data *config.Config) error {
	_, err := c.configCollection.InsertOne(c.ctx, *data)
	return err
}

func (c *Client) ResetConfig(botId int64, key string) error {
	_, err := c.configCollection.UpdateOne(c.ctx, idFilter(botId), bson.D{{Key: "$unset", Value: bson.D{{Key: key, Value: ""}}}})
	return err
}
