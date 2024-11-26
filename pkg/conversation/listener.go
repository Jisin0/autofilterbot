package conversation

import (
	"slices"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

// Listener contains data to filter incoming updates and handle them.
type Listener struct {
	filter      filters.Message
	messageChan chan *gotgbot.Message
}

// NewListener creates a new listener with given filter and message chan.
func NewListener(f filters.Message, c chan *gotgbot.Message) *Listener {
	return &Listener{
		filter:      f,
		messageChan: c,
	}
}

// ListenerArray is a thread-safe interface for working with an array of listener.
type ListenerArray struct {
	mu   sync.RWMutex
	list []*Listener
}

// NewListenerArray creates a new empty ListenerArray.
func NewListenerArray() *ListenerArray {
	return &ListenerArray{list: make([]*Listener, 0)}
}

// Add adds a listener to the listener.
func (ls *ListenerArray) Add(l *Listener) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.list = append(ls.list, l)
}

// FindMatchAndDelete finds a listener with a matching filter and deletes it, returning the listener.
func (ls *ListenerArray) FindMatchAndDelete(m *gotgbot.Message) (*Listener, bool) {
	for i, l := range ls.list {
		if !l.filter(m) {
			continue
		}

		ls.mu.Lock()
		defer ls.mu.Unlock()

		ls.list = slices.Delete(ls.list, i, i+1)
		return l, true
	}

	return nil, false
}

// Delete deletes listener at index i.
func (ls *ListenerArray) Delete(i int) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.list = slices.Delete(ls.list, i, i+1)
}
