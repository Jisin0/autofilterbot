package functions

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var _ error = (*FloodWaitError)(nil)

// FloodWaitError is a telegram rate limit error.
type FloodWaitError struct {
	Method   string
	Duration int64
}

func (f *FloodWaitError) Error() string {
	return fmt.Sprintf("429: unable to %s retry after %d", f.Method, f.Duration)
}

// AsFloodWait attemots to parse a telegram bot API floodwait error as a FloodWaitError.
func AsFloodWait(e error) (*FloodWaitError, bool) {
	rpc, ok := e.(*gotgbot.TelegramError)
	if !ok {
		return nil, false
	}

	if rpc.Code != 429 {
		return nil, false
	}

	if rpc.ResponseParams == nil || rpc.ResponseParams.RetryAfter == 0 {
		return nil, false
	}

	return &FloodWaitError{
		Method:   rpc.Method,
		Duration: rpc.ResponseParams.RetryAfter,
	}, true
}
