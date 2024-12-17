/*
Package configpanel creates a modular panel to
*/
package configpanel

import (
	"fmt"

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
	//TODO: handle result and error
	handleUpdate(p, ctx, bot)

	//TODO: create close button after markup length check

	return nil
}

// handleUpdate processes the update and returns the text and buttons to edit the message with.
// Returns a [PageNotFoundError] if page was not found.
func handleUpdate(p *Panel, update *ext.Context, bot *gotgbot.Bot) (string, [][]gotgbot.InlineKeyboardButton, error) {
	data := CallbackDataFromString(update.CallbackQuery.Data)

	ctx := &Context{
		App:           p.App,
		Bot:           bot,
		Update:        update,
		CallbackQuery: update.CallbackQuery,
		CallbackData:  data,
	}

	if len(data.Path) < 2 || data.Path[1] == PathHome {
		//TODO: edit with homepage and buttons
	}

	rootPage, ok := findPage(p.Pages, data.Path[1])
	if !ok {
		return "", nil, PageNotFoundError{pageName: data.Path[1]}
	}

	if len(data.Path) == 2 { // if only one subroute i.e display root page
		//TODO: handle pls
	}

	currentPage := rootPage

	for _, subRoutes := range data.Path[2:] {
		nextPage, ok := findPage(currentPage.SubPages, subRoutes)
		if !ok {
			//TODO: return not found msg
		}

		currentPage = nextPage
	}

	// If page has no subpages call CallbackFunc and return result
	// Otherwise return page content and generated markup

	if len(currentPage.SubPages) == 0 {
		return currentPage.CallbackFunc(ctx)
	}

	var (
		content = currentPage.Content
		markup  = buttonsFromPages(ctx.CallbackData, currentPage.SubPages)
	)

	if content == "" {
		content = fmt.Sprintf("<b>%s</b>\n\nUse the options below to configure this field 👇", currentPage.DisplayName)
	}

	return content, markup, nil
}
