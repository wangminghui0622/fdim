package util

// CalculateOffset 计算分页偏移�?
// PageNumber: 页码（从1开始）
// ShowNumber: 每页显示数量
func CalculateOffset(pageNumber, showNumber int32) int32 {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if showNumber < 1 {
		showNumber = 10
	}
	return (pageNumber - 1) * showNumber
}

// CalculateLimit 计算分页限制数量
func CalculateLimit(showNumber int32) int32 {
	if showNumber < 1 {
		return 10
	}
	return showNumber
}
