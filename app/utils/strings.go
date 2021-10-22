package utils

func StringSliceIn(element string, sliceInfo []string) bool {
	for _, item := range sliceInfo {
		if  item == element {
			return true
		}
	}
	return false
}