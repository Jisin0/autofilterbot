package functions

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var nonAlphaNumericRegex = regexp.MustCompile(`[^\w\s]+`)

// RemoveSymbols returns a copy of the string will all non alpha-numeric characters removed.
func RemoveSymbols(input string) string {
	// removes all symbols using regex and then splits into fields and rejoins to remove unnecessary whitespaces
	return strings.Join(strings.Fields(nonAlphaNumericRegex.ReplaceAllString(input, " ")), " ")
}

// RemoveExtension removes the extension from a file name if any.
func RemoveExtension(input string) string {
	if input == "" {
		return ""
	}

	index := strings.LastIndex(input, ".")
	if (len(input) - index) <= 4 { // if last index of . is within 4 character range of end of string then cut around it
		input = input[:index]
	}

	return input
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// RandString creates a randomly generated string of given length.
func RandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
