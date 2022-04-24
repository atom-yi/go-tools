package yutil

// 判断字符串是否为空
func IsBlankStr(str string) bool {
	if len(str) == 0 {
		return true
	}
	for _, c := range str {
		if c != ' ' {
			return false
		}
	}
	return true
}

// 判断字符串是否不为空
func IsNotBlankStr(str string) bool {
	return !IsBlankStr(str)
}
