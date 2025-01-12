package jsoncache_test

import (
	"testing"
	"time"

	"github.com/Jisin0/autofilterbot/pkg/jsoncache"
	"github.com/stretchr/testify/assert"
)

type MockData struct {
	Id         string
	FieldStr   string  `json:"field_str,omitempty"`
	FieldInt   int     `json:"field_int,omitempty"`
	FieldSlice []int64 `json:"field_slice,omitempty"`
}

func TestJsonCache(t *testing.T) {
	assert := assert.New(t)

	c := jsoncache.NewCache(".test", time.Minute*5)

	table := []MockData{
		{
			Id:         "t1",
			FieldStr:   "foo",
			FieldSlice: []int64{420, 69},
		},
		{
			Id:       "t2",
			FieldInt: 22,
			FieldStr: "hawk2",
		},
		{
			Id:         "t3",
			FieldSlice: []int64{0},
		},
	}

	for _, item := range table {
		assert.NoError(c.Save(item.Id, item))
		var d MockData
		assert.NoError(c.Load(item.Id, &d))
		assert.Equal(item, d)
	}

	assert.NoError(c.Close())
}
