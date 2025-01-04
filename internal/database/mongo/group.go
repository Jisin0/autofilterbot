package mongo

func (c *Client) SaveGroup(id int64) error {
	_, err := c.groupCollection.InsertOne(c.ctx, idFilter(id))
	return err
}
