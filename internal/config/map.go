package config

const (
	FieldNameFsub             = "fsub"
	FieldNameMaxResults       = "max_results"
	FiledNameMaxPages         = "max_pages"
	FieldNameMaxPerPage       = "max_per_page"
	FieldNameStart            = "start"
	FieldNameAbout            = "about"
	FieldNameHelp             = "help"
	FieldNamePrivacy          = "privacy"
	FieldNameStats            = "stats"
	FieldNameShortener        = "shortener"
	FieldNameNoResultText     = "no_result_text"
	FieldNameResultTemplate   = "af_template"
	FieldNameButtonTemplate   = "btn_template"
	FieldNameFdetailsTemplate = "fdetails_template"
	FieldNameSizeButton       = "size_btn"
	FieldNameAutodeleteTime   = "autodel_time"
)

// ToMap converts the contents of the struct into map so fields can be dynamically accessed.
func (c *Config) ToMap() map[string]any {
	if c.cachedMap == nil {
		c.RefreshMap()
	}

	return c.cachedMap
}

func (c *Config) toMap() map[string]any {
	vals := make(map[string]any)

	vals[FieldNameFsub] = c.GetFsubChannels()
	vals[FieldNameMaxResults] = c.GetMaxResults()
	vals[FiledNameMaxPages] = c.GetMaxPages()
	vals[FieldNameMaxPerPage] = c.GetMaxPerPage()

	// all message values are saved by prefix appended with _text and _buttons for text and markup
	vals[FieldNameStart] = c.GetStartMessage("")
	vals[FieldNameAbout] = c.GetAboutMessage()
	vals[FieldNameHelp] = c.GetHelpMessage()
	vals[FieldNamePrivacy] = c.GetPrivacyMessage()
	vals[FieldNameStats] = c.GetStatsMessage()

	vals[FieldNameShortener] = c.GetShortener()
	vals[FieldNameNoResultText] = c.GetNoResultText()
	vals[FieldNameResultTemplate] = c.GetResultTemplate()
	vals[FieldNameButtonTemplate] = c.GetButtonTemplate()
	vals[FieldNameFdetailsTemplate] = c.GetFileDetailsTemplate()
	vals[FieldNameSizeButton] = c.GetSizeButton()
	vals[FieldNameAutodeleteTime] = c.GetAutodeleteTime()

	return vals
}

// RefreshMap refreshes the value of the cached map.
func (c *Config) RefreshMap() {
	c.cachedMap = c.toMap()
}
