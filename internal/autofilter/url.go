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

// URLDataFromString decodes a string encoded using URLData.Encode().
func URLDataFromString(input string) (URLData, error) {
	var d URLData

	bytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return d, err
	}

	s := string(bytes)

	split := strings.Split(s, "|")
	if len(split) < 4 {
		return d, fmt.Errorf("not enough arguments")
	}

	d.FileUniqueId = split[1]

	i, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return d, err
	}

	d.ChatId = i

	d.HasShortener = split[3] == "1"

	return d, nil
}
