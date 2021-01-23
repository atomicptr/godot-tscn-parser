package parser

func intArrayContains(n int, arr []int) bool {
	for _, e := range arr {
		if e == n {
			return true
		}
	}
	return false
}
