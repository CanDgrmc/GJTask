package lib

func ExistsIntValue(array []int, needle int) bool {
	for _, i := range array {
		if i == needle {
			return true
		}
	}
	return false
}

func GetIndexOf(array []int, needle int) int {
	for index, i := range array {
		if i == needle {
			return index
		}
	}
	return -1
}
