package database

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// FileAlreadyExistsError indicates a similar or same file exists.
type FileAlreadyExistsError struct {
	FileName string
}

func (e FileAlreadyExistsError) Error() string {
	return fmt.Sprintf("file already exists: %s", e.FileName)
}

// IsNoDocumentsError reports wether the error is mongo no documnet found error for SingleResult methods.
func IsNoDocumentsError(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}
