package autofilter

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

// URLData wraps the data in a start url from autofilter results.
type URLData struct {
	// Id of the file.
	FileUniqueId string
	// Id of chat where query occurred.
	ChatId int64
	// Indicates whether shortened url should be given or actual file.
	HasShortener bool
}

// Encode converts data to a base64 string to be used in a start url.
func (d URLData) Encode() string {
	hasShortener := "0"
	if d.HasShortener {
		hasShortener = "1"
	}

	s := fmt.Sprintf("f|%s|%d|%s", d.FileUniqueId, d.ChatId, hasShortener)

	return base64.StdEncoding.EncodeToString([]byte(s))
}

// URLDataFromBase64String decodes the base64 encoded string and then parses it.
func URLDataFromBase64String(input string) (*URLData, error) {
	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return URLDataFromString(string(b))
}

// URLDataFromString decodes a string that is received as start data. The string must have been base64 decoded.
func URLDataFromString(input string) (*URLData, error) {

	split := strings.Split(input, "|")
	if len(split) < 4 {
		return nil, fmt.Errorf("not enough arguments")
	}

	i, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return nil, err
	}

	return &URLData{
		FileUniqueId: split[1],
		ChatId:       i,
		HasShortener: split[3] == "1",
	}, nil
}
