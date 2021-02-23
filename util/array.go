package util

import (
	"strconv"
	"strings"
)

// Int64ArrayToString will covert inter array to string with separator `,`
func Int64ArrayToString(input []int64, sep string) string {
	b := make([]string, len(input))
	for i, v := range input {
		b[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(b, ",")
}
