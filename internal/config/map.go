package config

// ToMap converts the contents of the struct into map so fields can be dynamically accessed.
func (c Config) ToMap() map[string]any {
	vals := make(map[string]any)

	if len(c.FsubChannels) != 0 {
		vals["fsub"] = c.FsubChannels
	}

	if c.MaxResults != 0 {
		vals["max_results"] = c.MaxResults
	}

	if c.MaxPages != 0 {
		vals["max_pages"] = c.MaxPages
	}

	if c.MaxPerPage != 0 {
		vals["max_per_page"] = c.MaxPerPage
	}

	if c.StartText != "" {
		vals["start_text"] = c.StartText
	}

	if len(c.StartButtons) != 0 {
		vals["start_buttons"] = c.StartButtons
	}

	if c.AboutText != "" {
		vals["about_text"] = c.AboutText
	}

	if len(c.AboutButtons) != 0 {
		vals["about_buttons"] = c.AboutButtons
	}

	if c.HelpText != "" {
		vals["help_text"] = c.HelpText
	}

	if len(c.HelpButtons) != 0 {
		vals["help_buttons"] = c.HelpButtons
	}

	if c.PrivacyText != "" {
		vals["privacy_text"] = c.PrivacyText
	}

	if len(c.PrivacyButtons) != 0 {
		vals["privacy_buttons"] = c.PrivacyButtons
	}

	if c.StatsText != "" {
		vals["stats_text"] = c.StatsText
	}

	if len(c.StatsButtons) != 0 {
		vals["stats_buttons"] = c.StatsButtons
	}

	return vals
}
