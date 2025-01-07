package functions_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/stretchr/testify/assert"
)

func TestFileSizeToString(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input          int64
		expectedOutput string
	}{
		{
			input:          512,
			expectedOutput: "512 B",
		},
		{
			input:          1000420,
			expectedOutput: "976.97 KB",
		}, {
			input:          302098567,
			expectedOutput: "288.10 MB",
		}, {
			input:          2894036420,
			expectedOutput: "2.70 GB",
		},
	}

	for _, item := range table {
		t.Run(item.expectedOutput, func(t *testing.T) {
			assert.Equal(item.expectedOutput, functions.FileSizeToString(item.input))
		})
	}
}
