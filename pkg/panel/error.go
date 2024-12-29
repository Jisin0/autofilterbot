package panel

import "fmt"

var _ error = PageNotFoundError{}

type PageNotFoundError struct {
	PageName string // name of the page or last route that wasn't found
}

func (e PageNotFoundError) Error() string {
	return fmt.Sprintf("page %s was not found", e.PageName)
}
