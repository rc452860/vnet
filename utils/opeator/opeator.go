package opeator

// IntIn implement python in operator
func IntIn(src int, to []int) bool {
	for _, item := range to {
		if src == item {
			return true
		}
	}
	return false
}
