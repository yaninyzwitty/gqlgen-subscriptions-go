package helpers

import "strconv"

func UintToString(id64 uint64) string {
	return strconv.FormatUint(id64, 10)
}
