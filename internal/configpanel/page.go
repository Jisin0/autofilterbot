package configpanel

import "github.com/PaulSonOfLars/gotgbot/v2"

type OperationType int

const (
	OperationTypeDelete OperationType = iota + 1
	OperationTypeAdd
)

// Page is a single page in the config panel.
type Page struct {
	// Name shown on the panel button.
	DisplayName string
	// Name is the unique identifier for this page. It should be as short as possible to not affect hit the callback data limit.
	Name string
	// Function to call when page is hit. returns the text & buttons to edit and an error.
	// The back button will be appended and does not need to be returned.
	CallbackFunc func(ctx *Context) (string, [][]gotgbot.InlineKeyboardButton, error)
	// Additional subpages.
	SubPages []*Page
}

// findPage searches for a page with matching name.
func findPage(pages []*Page, name string) (*Page, bool) {
	for _, p := range pages {
		if p.Name == name {
			return p, true
		}
	}

	return nil, false
}
