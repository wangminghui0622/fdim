package util

// HasDuplicate ƬǷظԪ??
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

// Contains ƬǷָԪ
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
