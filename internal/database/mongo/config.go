package mongo

import "github.com/Jisin0/autofilterbot/internal/config"

//TODO: implement

func (c *Client) GetConfig(botId int64) (*config.Config, error) {
	return nil, nil
}

func (c *Client) UpdateConfig(botId int64, key string, value interface{}) error {
	return nil
}

func (c *Client) SaveConfig(botId int64, data *config.Config) error {
	return nil
}
