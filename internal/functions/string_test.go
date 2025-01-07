package functions_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/stretchr/testify/assert"
)

func TestRemoveSymbols(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "-foo #(2024) 🩷 bar",
			expectedOutput: "foo 2024 bar",
		},
		{
			input:          "hello mom",
			expectedOutput: "hello mom",
		},
		{
			input:          "!@#$%^&*()}:?<>",
			expectedOutput: "",
		},
	}

	for _, item := range table {
		t.Run(item.expectedOutput, func(t *testing.T) {
			assert.Equal(item.expectedOutput, functions.RemoveSymbols(item.input))
		})
	}
}

func TestRemoveExtension(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input          string
		expectedOutput string
	}{
		{
			input:          "Just a file name.mkv",
			expectedOutput: "Just a file name",
		},
		{
			input:          "No extension",
			expectedOutput: "No extension",
		},
		{
			input:          "short extension.py",
			expectedOutput: "short extension",
		},
	}

	for _, item := range table {
		t.Run(item.expectedOutput, func(t *testing.T) {
			assert.Equal(item.expectedOutput, functions.RemoveExtension(item.input))
		})
	}
}
