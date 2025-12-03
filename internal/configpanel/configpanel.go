/*
Package configpanel handles the /settings command
*/
package configpanel

import (
	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/pkg/panel"
	"go.uber.org/zap"
)

const (
	OperationDelete = "del"
	OperationSet    = "set"
	OperationReset  = "reset"
)

type AppPreview interface {
	GetDB() database.Database
	GetConfig() *config.Config
	GetLog() *zap.Logger
	RefreshConfig()
	GetAdditionalCollectionCount() int
}

// CreatePanel creates the bot's configpanel and adds all pages.
func CreatePanel(app AppPreview) *panel.Panel {
	p := panel.NewPanel()

	p.AddPage(panel.NewPage("sizebtn", "Size Button").WithCallbackFunc(BoolField(app, config.FieldNameSizeButton)))
	p.AddPage(panel.NewPage("autodel", "Auto Delete").WithCallbackFunc(TimeField(app, config.FieldNameAutodeleteTime, []int{5, 10, 15, 20, 30, 45})))
	p.AddPage(panel.NewPage("filedel", "File AutoDelete").WithCallbackFunc(TimeField(app, config.FieldNameFileAutoDelete, []int{5, 10, 15, 20, 30, 45})))

	dbPage := panel.NewPage("db", "Database").WithContent("📂 Configure Database Settings from the Options Below.")
	dbPage.NewSubPage("coll", "File Database").WithCallbackFunc(IntField(app, config.FieldNameCollectionIndex, IntFieldOpts{
		Range:       &IntRange{Start: 0, End: app.GetAdditionalCollectionCount()},
		Description: "Collection/Database to Store Files. 0 is your Main Database.",
	}))
	dbPage.NewSubPage("updater", "Auto Collection Updater").WithCallbackFunc(BoolField(
		app,
		config.FieldNameCollectionUpdater,
		"The Auto Collection-Updater periodically runs to change the database used to store files, to the database with least files.\n\nWhen Enabled, the Collection set from Config Panel will bee Ignored\n\nNOTE: Application must be restarted for changes to take effect.\n\n",
	))

	p.AddPage(dbPage)

	return p
}
