package mongo

import (
	"context"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Client) SaveFile(f *model.File) error {
	// Find any with matching file_id
	if res := c.fileCollection.FindOne(c.ctx, fileIdFilter(f.FileId)); res.Err() != mongo.ErrNoDocuments {
		return database.FileAlreadyExistsError{FileName: f.FileName}
	}

	// Find a document that starts with the same file_name and is within a 100 byte range of file_size
	duplicateFilter := bson.D{
		{Key: "file_name", Value: bson.D{{Key: "$regex", Value: "^" + f.FileName}}},
		{Key: "file_size", Value: bson.D{
			{Key: "$gte", Value: f.FileSize - 100},
			{Key: "$lte", Value: f.FileSize + 100},
		}},
	}
	if res := c.fileCollection.FindOne(c.ctx, duplicateFilter); res.Err() != mongo.ErrNoDocuments {
		return database.FileAlreadyExistsError{FileName: f.FileName}
	}

	_, err := c.fileCollection.InsertOne(c.ctx, f)
	return err
}

func (c *Client) SaveFiles(files ...*model.File) []error {
	var errs []error
	for _, f := range files {
		if err := c.SaveFile(f); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (c *Client) GetFile(fileId string) (*model.File, error) {
	res := c.fileCollection.FindOne(c.ctx, fileIdFilter(fileId))
	if err := res.Err(); err != nil {
		return nil, err
	}

	var f model.File

	res.Decode(&f)

	return &f, nil
}

func (c *Client) DeleteFile(fileId string) error {
	_, err := c.fileCollection.DeleteOne(c.ctx, fileIdFilter(fileId))
	return err
}

func (c *Client) SearchFiles(query string) (database.Cursor, error) {
	pattern := `(?i)(\b|[\.\+\-_])` + strings.ReplaceAll(query, " ", `.*[\s\.\+\-_]`) + `(\b|[\.\+\-_])`
	pipeline := bson.D{{Key: "file_name", Value: bson.D{{Key: "$regex", Value: pattern}}}}

	return c.fileCollection.Find(context.TODO(), pipeline, options.Find().SetSort(bson.M{"time": -1}).SetLimit(50))
}

// fileIdFilter creates a bson filter to match by file_id.
func fileIdFilter(id string) bson.D {
	return bson.D{{Key: "file_id", Value: id}}
}
