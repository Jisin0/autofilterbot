package configpanel

import "fmt"

var _ error = PageNotFoundError{}

type PageNotFoundError struct {
	pageName string // name of the page or last route that wasn't found
}

func (e PageNotFoundError) Error() string {
	return fmt.Sprintf("page %s was not found", e.pageName)
}
