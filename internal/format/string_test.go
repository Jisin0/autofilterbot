package format_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/format"
	"github.com/stretchr/testify/assert"
)

func TestKeyValueFormat(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		template string
		values   map[string]string

		expectedOutput string
	}{
		{
			template: "Hey {key} I'm not famous!",
			values: map[string]string{
				"bad": "value",
				"key": "mom",
			},

			expectedOutput: "Hey mom I'm not famous!",
		},
		{
			template: "This has no values",
			values: map[string]string{
				"foo": "bar",
			},

			expectedOutput: "This has no values",
		},
		{
			template:       "Should be unchanged",
			expectedOutput: "Should be unchanged",
		},
	}

	for _, item := range table {
		assert.Equal(item.expectedOutput, format.KeyValueFormat(item.template, item.values))
	}
}
