/*
Package autofilter contains types and methods to help work with autofilter results.
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
func FilesFromCursor(ctx context.Context, c database.Cursor, opts FilesFromCursorOptions) ([]Files, error) {
	var (
		totalCount int
		finished   bool
		totalFiles = make([]Files, 0, opts.GetMaxResults())
	)

	for i := 0; i < opts.GetMaxPages(); i++ {
		row := make([]File, 0, opts.GetMaxPerPage())

		for j := 0; j < opts.GetMaxPerPage(); j++ {
			if !c.Next(ctx) {
				finished = true
				break
			}

			var f model.File

			err := c.Decode(&f)
			if err != nil {
				return totalFiles, err
			}

			row = append(row, File{File: f})
		}

		if len(row) != 0 {
			totalFiles = append(totalFiles, row)
		}

		if finished {
			return totalFiles, nil
		}

		if totalCount >= opts.GetMaxResults() { // only checks after completing a page. should this behaviour be changed?
			break
		}
	}

	return totalFiles, nil
}
