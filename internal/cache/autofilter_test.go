package cache_test

import (
	"testing"
	"time"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/Jisin0/autofilterbot/internal/cache"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAutofilter(t *testing.T) {
	assert := assert.New(t)

	data := &autofilter.SearchResult{
		UniqueId: "ssabmud",
		Query:    "hello mom",
		FromUser: 69420,
		ChatID:   123456789,
		Files: []autofilter.Files{
			{{File: model.File{FileId: "QinsiSYA8ysa", FileName: "This Is A Cute Cate Video.mkv", FileType: "video", FileSize: 1 << 24}}},
			{{File: model.File{FileId: "Hkwosmd_6dsn", FileName: "Deadpool.&.Wolverine.2024.Trailer.x264.AAC.mkv", FileSize: 132987239182}}},
		},
	}

	c := cache.NewAutofilter(time.Minute * 1)

	err := c.Save(data)
	assert.NoError(err)

	res, _, err := c.Get("ssabmud")
	assert.NoError(err)

	assert.Equal(data, res)
}
