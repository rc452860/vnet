package utils

func StringArrayContain(slice []string, ele string) bool {
	for _, i := range slice {
		if i == ele {
			return true
		}
	}
	return false
}
