package handlers_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/handlers"
	"github.com/stretchr/testify/assert"
)

func TestFormatString(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input  string
		values map[string]string
		output string
	}{
		{
			input: "Hello {var}",
			values: map[string]string{
				"var": "mom",
				"foo": "bar",
			},
			output: "Hello mom",
		},
		{
			input: "Hello {var}",
			values: map[string]string{
				"var": "World",
			},
			output: "Hello World",
		},
		{
			input:  "Hello {var}",
			output: "Hello {var}",
		},
	}

	for _, item := range table {
		t.Run(item.output, func(t *testing.T) {
			o := handlers.FormatString(item.input, item.values)
			assert.Equal(item.output, o)
		})
	}
}
