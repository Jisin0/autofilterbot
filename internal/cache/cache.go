/*
Package cache manages all json caches used in the project.
*/
package cache

import "time"

const (
	defualtAutofilterTimeout = time.Minute * 15
)

// Cache wraps all json cache helper types into a struct.
type Cache struct {
	Autofilter *Autofilter
}

// NewCache initalizes and creates a new cache structure.
func NewCache() *Cache {
	return &Cache{
		Autofilter: NewAutofilter(defualtAutofilterTimeout),
	}
}
