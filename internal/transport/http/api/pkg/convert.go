package pkg

import "strconv"

// ParseInt64 parses a string to int64.
func ParseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseUInt64 parses a string to uint64.
func ParseUInt64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
