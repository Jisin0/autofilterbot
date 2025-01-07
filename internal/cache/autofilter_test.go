package cache_test

import (
	"testing"
	"time"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/Jisin0/autofilterbot/internal/cache"
	"github.com/stretchr/testify/assert"
)

func TestAutofilter(t *testing.T) {
	assert := assert.New(t)

	data := &autofilter.SearchResult{
		Query:    "hello mom",
		FromUser: 69420,
		ChatID:   123456789,
		Files: []autofilter.Files{
			{{FileId: "QinsiSYA8ysa", FileName: "This Is A Cute Cate Video.mkv", FileType: "video", FileSize: 1 << 24}},
			{{FileId: "Hkwosmd_6dsn", FileName: "Deadpool.&.Wolverine.2024.Trailer.x264.AAC.mkv", FileSize: 132987239182}},
		},
	}

	c := cache.NewAutofilter(time.Minute * 1)

	err := c.Save(data)
	assert.NoError(err)

	res, err, _ := c.Get("hello mom", 123456789)
	assert.NoError(err)

	assert.Equal(data, res)
}
