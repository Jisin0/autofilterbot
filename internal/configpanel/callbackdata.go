package configpanel

import "strings"

const (
	// Character that joins paths (colon)
	PathDelimiter = ':'
	// Character that joins arguments (underscore)
	ArgDelimiter = '_'
	// Divides two sections in the data namely path and arguments (pipe)
	SectionDelimiter = '|'
)

// CallbackDataFromString parses the string, usually the raw callback data, and returns a CallbackData object.
// the data should be a string in the format:
//
//	<path1>:<path2>:<path3...>|<arg1>_<arg2>_<arg3...>
func CallbackDataFromString(s string) CallbackData {
	sections := strings.SplitN(s, string(SectionDelimiter), 1)

	newPath := CallbackData{
		Data: s,
	}

	// section[0] should be path and must be present
	newPath.Path = strings.Split(sections[0], string(PathDelimiter))

	if len(sections) > 1 {
		newPath.Args = strings.Split(sections[1], string(ArgDelimiter))
	}

	return newPath
}

// CallbackData wraps the raw callback data which represents the path of the request with some useful methods.
type CallbackData struct {
	// Raw data
	Data string
	// Parsed path with index 0 being config
	Path []string
	// Optional arguments
	Args []string
}

// String stringifies the data to be used in buttons as callback data.
func (c CallbackData) String() string {
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
func (c CallbackData) AddArg(val string) CallbackData {
	return CallbackData{
		Data: c.Data,
		Path: c.Path,
		Args: append(c.Args, val),
	}
}

// AddPath adds a subpath to the end of existing paths.
func (c CallbackData) AddPath(val string) CallbackData {
	return CallbackData{
		Data: c.Data,
		Path: append(c.Path, val),
		Args: c.Args,
	}
}
