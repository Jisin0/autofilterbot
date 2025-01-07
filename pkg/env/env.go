package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Int64s parses the environment variable as a slice of int64, ignoring any errors.
func Int64s(name string) (l []int64) {
	s := os.Getenv(name)
	if s == "" {
		return
	}

	for _, str := range strings.Fields(s) {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			fmt.Println(err) // bad practice blah blah
			continue
		}

		l = append(l, i)
	}

	return
}
