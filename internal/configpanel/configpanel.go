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
}

// CreatePanel creates the bot's configpanel and adds all pages.
func CreatePanel(app AppPreview) *panel.Panel {
	p := panel.NewPanel()

	p.AddPage(panel.NewPage("sizebtn", "Size Button").WithCallbackFunc(BoolField(app, config.FieldNameSizeButton)))
	p.AddPage(panel.NewPage("autodel", "Auto Delete").WithCallbackFunc(TimeField(app, config.FieldNameAutodeleteTime, []int{5, 10, 15, 20, 30, 45})))
	p.AddPage(panel.NewPage("filedel", "File AutoDelete").WithCallbackFunc(TimeField(app, config.FieldNameFileAutoDelete, []int{5, 10, 15, 20, 30, 45})))

	return p
}
