/*
Modified version of gotgbot filters/command to handle multiple commands
*/
package exthandlers

import (
	"strings"
	"unicode/utf8"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// ensure Commands implements ext.Handler
var _ ext.Handler = Commands{}

// Commands is the go-to handler for setting up Commands in your bot. By default, it will use telegram-native commands
// that start with a forward-slash (/), but it can be customised to react to any message starting with a character.
//
// For example, a command handler on "help" with the triggers []rune("/!,") would trigger for "/help", "!help", or ",help".
type Commands struct {
	Triggers     []rune
	AllowEdited  bool
	AllowChannel bool
	Commands     []string // set to a lowercase value for case-insensitivity
	Response     handlers.Response
}

// NewCommand creates a new case-insensitive command.
// By default, commands do not work on edited messages, or channel posts. These can be enabled by setting the
// AllowEdited and AllowChannel fields respectively.
func NewCommands(c []string, r handlers.Response) Commands {
	return Commands{
		Triggers:     []rune{'/'},
		AllowEdited:  false,
		AllowChannel: false,
		Commands:     toLowerSlice(c),
		Response:     r,
	}
}

// SetAllowEdited Enables edited messages for this handler.
func (c Commands) SetAllowEdited(allow bool) Commands {
	c.AllowEdited = allow
	return c
}

// SetAllowChannel Enables channel messages for this handler.
func (c Commands) SetAllowChannel(allow bool) Commands {
	c.AllowChannel = allow
	return c
}

// SetTriggers sets the list of triggers to be used with this command.
func (c Commands) SetTriggers(triggers []rune) Commands {
	c.Triggers = triggers
	return c
}

func (c Commands) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.Message != nil {
		if ctx.Message.GetText() == "" {
			return false
		}
		return c.checkMessage(b, ctx.Message)
	}

	// if no edits and message is edited
	if c.AllowEdited && ctx.EditedMessage != nil {
		if ctx.EditedMessage.GetText() == "" {
			return false
		}
		return c.checkMessage(b, ctx.EditedMessage)
	}
	// if no channel and message is channel message
	if c.AllowChannel && ctx.ChannelPost != nil {
		if ctx.ChannelPost.GetText() == "" {
			return false
		}
		return c.checkMessage(b, ctx.ChannelPost)
	}
	// if no channel, no edits, and post is edited
	if c.AllowChannel && c.AllowEdited && ctx.EditedChannelPost != nil {
		if ctx.EditedChannelPost.GetText() == "" {
			return false
		}
		return c.checkMessage(b, ctx.EditedChannelPost)
	}

	return false
}

func (c Commands) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	return c.Response(b, ctx)
}

func (c Commands) Name() string {
	return "commands_" + c.Commands[0] // not ideal but meh
}

func (c Commands) checkMessage(b *gotgbot.Bot, msg *gotgbot.Message) bool {
	ents := msg.GetEntities()
	if len(ents) != 0 && ents[0].Offset == 0 && ents[0].Type != "bot_command" {
		return false
	}

	text := msg.GetText()

	var cmd string
	for _, t := range c.Triggers {
		if r, _ := utf8.DecodeRuneInString(text); r != t {
			continue
		}

		split := strings.Split(strings.ToLower(strings.Fields(text)[0]), "@")
		if len(split) > 1 && !strings.EqualFold(split[1], b.User.Username) {
			return false
		}
		cmd = split[0][1:]
		break
	}

	if cmd == "" {
		return false
	}

	return contains(c.Commands, cmd)
}

// toLowerSlice comverts a slice of strings to lower case.
func toLowerSlice(v []string) (out []string) {
	for _, s := range v {
		out = append(out, strings.TrimSpace((strings.ToLower(s))))
	}
	return
}

// contains reports whether value s is in slice a.
func contains(a []string, v string) bool {
	for _, s := range a {
		if s == v {
			return true
		}
	}
	return false
}
