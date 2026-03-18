package util

// HasDuplicate 检查切片中是否有重复元�?
func HasDuplicate(slice []string) bool {
	seen := make(map[string]bool)
	for _, s := range slice {
		if seen[s] {
			return true
		}
		seen[s] = true
	}
	return false
}

// Contains 检查切片中是否包含指定元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
