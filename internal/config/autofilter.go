package config

const (
	defaultMaxResults = 50
	defaultMaxPerPage = 10
	defaultMaxPages   = 5
)

func (c *Config) GetMaxResults() int {
	if c.MaxResults != 0 {
		return c.MaxResults
	}

	return defaultMaxResults
}

func (c *Config) GetMaxPerPage() int {
	if c.MaxPerPage != 0 {
		return c.MaxPerPage
	}

	return defaultMaxPerPage
}

func (c *Config) GetMaxPages() int {
	if c.MaxResults != 0 {
		return c.MaxResults
	}

	return defaultMaxPages
}

func (c *Config) GetResultTemplate() string {
	if c.ResultTemplate != "" {
		return c.ResultTemplate
	}

	return `<b>🗣️ Hᴇʏ <tg-spoiler>{mention}</tg-spoiler> !
ʜᴇʀᴇ's ᴡʜᴀᴛ ɪ ғᴏᴜɴᴅ ғᴏʀ </b> <code>{query}</code> 👇

{warn}`
}

func (c *Config) GetNoResultText() string {
	if c.NoResultText != "" {
		return c.NoResultText
	}

	return `<b>❌ I'ᴍ sᴏʀʀʏ {mention} I ᴄᴏᴜʟᴅɴ'ᴛ ғɪɴᴅ ᴀɴʏ ʀᴇsᴜʟᴛs ғᴏʀ <code>{query}</code>
<blockquote>♻️ ᴘʟᴇᴀsᴇ ᴅᴏᴜʙʟᴇ ᴄʜᴇᴄᴋ ᴛʜᴇ sᴘᴇʟʟɪɴɢ ᴏғ ᴛʜᴇ ғɪʟᴇ ʏᴏᴜ ᴡᴀɴᴛ ᴏʀ ᴄʟɪᴄᴋ ᴏɴ ᴛʜᴇ ʙᴜᴛᴛᴏɴ ʙᴇʟᴏᴡ ᴛᴏ ɢᴏᴏɢʟᴇ ғᴏʀ sᴜɢɢᴇsᴛɪᴏɴs 👇</blockquote></b>`
}

func (c *Config) GetButtonTemplate() string {
	if c.ButtonTemplate != "" {
		return c.ButtonTemplate
	}

	return "📂 {file_size} {file_name}"
}

func (c *Config) GetSizeButton() bool {
	return c.SizeButton
}
