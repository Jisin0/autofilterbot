package autofilter_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/stretchr/testify/assert"
)

func TestURLData(t *testing.T) {
	assert := assert.New(t)

	table := []*autofilter.URLData{
		{
			FileUniqueId: "SDKHjshkdja5645a==addkov_Askdj",
			ChatId:       -100123456789,
			HasShortener: true,
		},
		{
			FileUniqueId: "==AQSW5654is_88sdiJHSpPP",
			ChatId:       209812831,
		},
		{
			FileUniqueId: "__89xczijhDSJsgjdJSHDAJs==",
			ChatId:       -87654321,
			HasShortener: true,
		},
	}

	for _, item := range table {
		t.Run(item.FileUniqueId, func(t *testing.T) {
			s := item.Encode()

			out, err := autofilter.URLDataFromBase64String(s)
			assert.NoError(err)
			assert.Equal(item, out)
		})
	}
}
