package autofilter_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/stretchr/testify/assert"
)

func TestURLData(t *testing.T) {
	assert := assert.New(t)

	table := []autofilter.URLData{
		{
			FileId:       "SDKHjshkdja5645a==addkov_Askdj",
			ChatId:       -100123456789,
			HasShortener: true,
		},
		{
			FileId: "==AQSW5654is_88sdiJHSpPP",
			ChatId: 209812831,
		},
		{
			FileId:       "__89xczijhDSJsgjdJSHDAJs==",
			ChatId:       -87654321,
			HasShortener: true,
		},
	}

	for _, item := range table {
		t.Run(item.FileId, func(t *testing.T) {
			s := item.Encode()

			out, err := autofilter.URLDataFromString(s)
			assert.NoError(err)
			assert.Equal(item, out)
		})
	}
}
