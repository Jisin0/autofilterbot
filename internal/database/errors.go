package database

import "fmt"

// FileAlreadyExistsError indicates a simliar or same file exists.
type FileAlreadyExistsError struct {
	FileName string
}

func (e FileAlreadyExistsError) Error() string {
	return fmt.Sprintf("file already exists: %s", e.FileName)
}
