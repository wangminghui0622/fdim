package util

// CalculateOffset ҳƫ??
// PageNumber: ҳ루1ʼ
// ShowNumber: ÿҳʾ
func CalculateOffset(pageNumber, showNumber int32) int32 {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if showNumber < 1 {
		showNumber = 10
	}
	return (pageNumber - 1) * showNumber
}

// CalculateLimit ҳ
func CalculateLimit(showNumber int32) int32 {
	if showNumber < 1 {
		return 10
	}
	return showNumber
}
