package functions_test

import (
	"fmt"
	"testing"

	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/stretchr/testify/assert"
)

func TestParseMessageLink(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input          string
		expectedOutput *functions.MessageLink
	}{
		{
			input: "https://t.me/durov/12",
			expectedOutput: &functions.MessageLink{
				Username:  "durov",
				MessageId: 12,
			},
		},
		{
			input: "t.me/c/2235842999/1",
			expectedOutput: &functions.MessageLink{
				ChatId:    -1002235842999,
				MessageId: 1,
			},
		},
		{
			input: "t.me/MyUsername/123456",
			expectedOutput: &functions.MessageLink{
				Username:  "MyUsername",
				MessageId: 123456,
			},
		},
		{
			input: "https://t.me/c/1234567890/9876",
			expectedOutput: &functions.MessageLink{
				ChatId:    -1001234567890,
				MessageId: 9876,
			},
		},
	}

	for _, item := range table {
		t.Run(fmt.Sprintf(
			"%d_%s_%d",
			item.expectedOutput.ChatId,
			item.expectedOutput.Username,
			item.expectedOutput.MessageId,
		), func(t *testing.T) {
			r, err := functions.ParseMessageLink(item.input)

			assert.NoError(err)
			assert.Equal(item.expectedOutput, r)
		})
	}
}
