package utils

// DerefInt64 解引用int64指针，如果为nil返回0
func DerefInt64(val *int64) int64 {
	if val == nil {
		return 0
	}
	return *val
}

// DerefInt 解引用int指针，如果为nil返回0
func DerefInt(val *int) int {
	if val == nil {
		return 0
	}
	return *val
}

// DerefString 解引用string指针，如果为nil返回空字符串
func DerefString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

// DerefFloat64 解引用float64指针，如果为nil返回0.0
func DerefFloat64(val *float64) float64 {
	if val == nil {
		return 0.0
	}
	return *val
}
