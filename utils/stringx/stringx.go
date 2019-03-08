package stringx

import "strings"

// IsDigit judement data is any of 1234567890
func IsDigit(data string) bool {
	if len(data) != 1 {
		return false
	}
	if strings.IndexAny(data, "1234567890") != -1 {
		return true
	}
	return false
}
