/*
Package autofilter contains types and methods to help geat autofilter results easily.
*/
package autofilter

import (
	"context"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/model"
)

type FilesFromCursorOptions interface {
	GetMaxResults() int
	GetMaxPages() int
	GetMaxPerPage() int
}

// FilesFromCursor loops through a cursor and outputs an array of files.
func FilesFromCursor(ctx context.Context, c database.Cursor, opts FilesFromCursorOptions) ([][]model.File, error) {
	var (
		totalCount   int
		totalButtons = make([][]model.File, 0, opts.GetMaxResults())
	)

	for i := 0; i < opts.GetMaxPages(); i++ {
		row := make([]model.File, 0, opts.GetMaxPerPage())

		for j := 0; j < opts.GetMaxPerPage(); j++ {
			if !c.Next(ctx) {
				return totalButtons, nil
			}

			var f model.File

			err := c.Decode(&f)
			if err != nil {
				return totalButtons, err
			}

			row = append(row, f)
		}

		totalButtons = append(totalButtons, row)

		if totalCount >= opts.GetMaxResults() { // only checks after completing a page. should this behaviour be changed?
			break
		}
	}

	return totalButtons, nil
}
