package object

// convertRange 对 i 进行转换，使其满足 0 <= i <= n
func convertRange(i, n int) int {
	if i < 0 {
		i += n
	}
	if i < 0 {
		return 0
	} else if i > n {
		return n
	} else {
		return i
	}
}
