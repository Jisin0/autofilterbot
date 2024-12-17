package configpanel

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type OperationType int

const (
	OperationTypeDelete OperationType = iota + 1
	OperationTypeAdd
)

// Page is a single page in the config panel.
type Page struct {
	// Name shown on the panel button.
	DisplayName string
	// Content is the text content of the page.
	Content string
	// Function to dynamically generate content for the page instead of the static Content.
	ContentGenerator ContentGenerator
	// Name is the unique identifier for this page. It should be as short as possible to not affect the callback data limit.
	Name string
	// Function to call when page is hit. returns the text & buttons to edit and an error.
	// The back button will be appended and does not need to be returned.
	CallbackFunc func(ctx *Context) (string, [][]gotgbot.InlineKeyboardButton, error)
	// Additional subpages.
	SubPages []*Page
}

// WithContentGenerator sets the ContentGenerator field of the page.
func (p *Page) WithContentGenerator(g ContentGenerator) *Page {
	p.ContentGenerator = g
	return p
}

// GetContent returns the COntent for the page or the default content if empty.
func (p *Page) GetContent() string {
	if p.ContentGenerator != nil {
		return p.ContentGenerator()
	}

	if p.Content != "" {
		return p.Content
	}

	return fmt.Sprintf("<b>%s</b>\n\nUse the options below to configure this field 👇", p.DisplayName)
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
