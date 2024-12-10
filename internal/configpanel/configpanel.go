/*
Package configpanel creates a modular panel to
*/
package configpanel

import (
	"github.com/Jisin0/autofilterbot/internal/app"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

const (
	PathHome = "home"
)

// NewPanel intializes a new empty config panel.
func NewPanel(app *app.App) *Panel {
	return &Panel{
		App:   app,
		Pages: make([]*Page, 5),
	}
}

// Panel is the entrypoint to the configpanel to which pages can be added. Use the NewPanel function to create a new panel.
type Panel struct {
	// List of pages.
	Pages []*Page
	// Application instance.
	App *app.App
}

// AddPage adds a new page to the root panel.
func (p *Panel) AddPage(page *Page) *Page {
	p.Pages = append(p.Pages, page)
	return page
}

// HandleUpdate processes the update and runs it through the config panel.
func (p *Panel) HandleUpdate(ctx *ext.Context, bot *gotgbot.Bot) error {
	data := CallbackDataFromString(ctx.CallbackQuery.Data)

	_ = &Context{
		App:           p.App,
		Update:        ctx,
		CallbackQuery: ctx.CallbackQuery,
		CallbackData:  data,
	}

	if len(data.Path) < 2 || data.Path[1] == PathHome {

	}

	return nil
}
