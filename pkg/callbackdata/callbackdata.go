/*
Package callbackdata allows working with data through data structures in a reliable way instead of working with strings.
*/
package callbackdata

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

//TODO: write tests

const (
	// Character that joins paths (colon)
	PathDelimiter = ':'
	// Character that joins arguments (underscore)
	ArgDelimiter = '_'
	// Divides two sections in the data namely path and arguments (pipe)
	SectionDelimiter = '|'
)

// FromString parses the string, usually the raw callback data, and returns a CallbackData object.
// the data should be a string in the format:
//
//	<path1>:<path2>:<path3...>|<arg1>_<arg2>_<arg3...>
func FromString(s string) CallbackData {
	sections := strings.SplitN(s, string(SectionDelimiter), 2)

	cB := CallbackData{}

	// section[0] should be path and must be present
	cB.Path = strings.Split(sections[0], string(PathDelimiter))

	if len(sections) > 1 {
		cB.Args = strings.Split(sections[1], string(ArgDelimiter))
	}

	return cB
}

// New creates a new empty callback data structure.
func New() *CallbackData {
	return &CallbackData{}
}

// CallbackData wraps the raw callback data which represents the path of the request with some useful methods.
type CallbackData struct {
	// Parsed path with index 0 being config
	Path []string
	// Optional arguments
	Args []string
}

// ToString stringifies the data to be used in buttons as callback data.
func (c *CallbackData) ToString() string {
	var b strings.Builder

	// write first path element which should be the root config
	// should always be set or will otherwise panic
	b.WriteString(c.Path[0])

	for _, p := range c.Path[1:] {
		b.WriteRune(PathDelimiter)
		b.WriteString(p)
	}

	if len(c.Args) > 0 {
		b.WriteRune(SectionDelimiter)
		b.WriteString(c.Args[0])

		for _, i := range c.Args[1:] {
			b.WriteRune(ArgDelimiter)
			b.WriteString(i)
		}
	}

	return b.String()
}

// AddArg appends an argument to the end of the callbackdata.
func (c *CallbackData) AddArg(val string) *CallbackData {
	c.Args = append(c.Args, val)
	return c
}

// AddArgs adds appends a list of arguments.
func (c *CallbackData) AddArgs(vals ...string) *CallbackData {
	c.Args = append(c.Args, vals...)
	return c
}

// AddPath adds a subpath to the end of existing paths.
func (c *CallbackData) AddPath(val string) *CallbackData {
	c.Path = append(c.Path, val)
	return c
}

// RemoveLastPath removes the last path in the callback data.
func (c *CallbackData) RemoveLastPath() *CallbackData {
	// if less than 2 paths then returned unchaged.
	if len(c.Path) < 2 {
		return c
	}

	c.Path = c.Path[:len(c.Path)-1]
	return c
}

// RemoveArgs removes all args from the callback data.
func (c *CallbackData) RemoveArgs() *CallbackData {
	c.Args = nil
	return c
}

// TODO: create BackButton bound method to generate back button to last route and implement at points of error
// BackOrCloseButton creates either a back button if applicable or a close button from the data.
func (c *CallbackData) BackOrCloseButton(userId ...int64) gotgbot.InlineKeyboardButton {
	if len(c.Path) <= 1 {
		return closeButton(userId...)
	} else {
		// nested page so add back button
		return gotgbot.InlineKeyboardButton{Text: "<- Back", CallbackData: c.RemoveArgs().RemoveLastPath().ToString()}
	}
}

func closeButton(userId ...int64) gotgbot.InlineKeyboardButton {
	data := New().AddPath("close")
	for _, u := range userId {
		data.AddArg(fmt.Sprint(u))
	}

	return gotgbot.InlineKeyboardButton{
		Text:         "𝖢𝗅𝗈𝗌𝖾",
		CallbackData: data.ToString(),
	}
}
