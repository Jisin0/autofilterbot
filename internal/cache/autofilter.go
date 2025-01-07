package cache

import (
	"errors"
	"time"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/Jisin0/autofilterbot/pkg/jsoncache"
)

// Autofilter manages the cache for autofilter results.
type Autofilter struct {
	cache *jsoncache.Cache
}

// Save saved the search results to a json file.
func (c *Autofilter) Save(data *autofilter.SearchResult) error {
	if data.UniqueId == "" {
		return errors.New("id is empty")
	}
	return c.cache.Save(data.UniqueId, *data)
}

// Get fetches the results from json cache.
//
// If the cache file doesnt exist or was expired a nil error with ok set to false will be returned.
// In case of any other error ok will be true and error will be set.
func (c *Autofilter) Get(uniqueId string) (*autofilter.SearchResult, error, bool) {
	var res autofilter.SearchResult

	err := c.cache.Load(uniqueId, &res)
	if err != nil {
		if err == jsoncache.ErrFileNotFound || err == jsoncache.ErrCacheDataExpired {
			return nil, nil, false
		}

		return nil, err, true
	}

	return &res, nil, true
}

func NewAutofilter(timeout time.Duration) *Autofilter {
	return &Autofilter{
		cache: jsoncache.NewCache(".results", timeout),
	}
}
