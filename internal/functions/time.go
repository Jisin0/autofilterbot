package functions

import "time"

// FormatUnixTimestamp converts a unix timestamp into a string in the dd/mm/yyyy format.
func FormatUnixTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("02/01/2006")
}
