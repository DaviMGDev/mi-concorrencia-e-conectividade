package utils

import "strconv"

var count uint64 = 0

func Count() string {
	count++
	return strconv.FormatUint(count, 10)
}
