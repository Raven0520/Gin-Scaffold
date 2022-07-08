package util

// InSliceString 判断字符串是否在切片中
func InSliceString(needle string, slices []string) bool {
	for _, value := range slices {
		if value == needle {
			return true
		}
	}
	return false
}

// InSliceInt8 判断 数字是否在切片中
func InSliceInt8(needle int8, slices []int8) bool {
	for _, n := range slices {
		if n == needle {
			return true
		}
	}
	return false
}
