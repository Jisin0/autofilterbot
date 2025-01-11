package panel

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type CallbackFunc func(ctx *Context) (string, [][]gotgbot.InlineKeyboardButton, error)

// Page is a single page in the config panel.
type Page struct {
	// Name is the unique identifier for this page. It should be as short as possible to not affect the callback data limit.
	Name string
	// Name shown on the panel button.
	DisplayName string
	// Content is the text content of the page.
	Content string
	// Function to dynamically generate content for the page instead of the static Content.
	ContentGenerator ContentGenerator
	// Function to call when page is hit. returns the text & buttons to edit and an error.
	// The back button will be appended and does not need to be returned.
	CallbackFunc CallbackFunc
	// Additional subpages. CallbackFunc becomes obsolete if set.
	SubPages []*Page
}

// WithContentGenerator sets the ContentGenerator field of the page.
func (p *Page) WithContentGenerator(g ContentGenerator) *Page {
	p.ContentGenerator = g
	return p
}

// WithContent sets the Content field of the page.
func (p *Page) WithContent(val string) *Page {
	p.Content = val
	return p
}

// WithContent sets the Content field of the page.
func (p *Page) WithCallbackFunc(val CallbackFunc) *Page {
	p.CallbackFunc = val
	return p
}

// AddSubPage adds a new sub page.
func (p *Page) AddSubPage(page *Page) *Page {
	p.SubPages = append(p.SubPages, page)
	return page
}

// NewPage creates a new empty subpage with given name and displayName.
func (p *Page) NewSubPage(name, displayName string) *Page {
	newPage := &Page{
		Name:        name,
		DisplayName: displayName,
	}
	p.SubPages = append(p.SubPages, newPage)

	return newPage
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

// NewPage creates a new empty page with name and display name shown on button.
//
// NOTE: name should be short and meaningful as it is put in callback data, displayName is shown as button label.
func NewPage(name, displayName string) *Page {
	return &Page{
		Name:        name,
		DisplayName: displayName,
	}
}
