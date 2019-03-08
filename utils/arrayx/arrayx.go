package arrayx

func FindStringInArray(obj string, target []string) bool {
	for _, item := range target {
		if obj == item {
			return true
		}
	}
	return false
}
