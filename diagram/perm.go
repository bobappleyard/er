package diagram

func firstPerm(n int) []int {
	p := make([]int, n)
	for i := range p {
		p[i] = i
	}
	return p
}

func nextPerm(p []int) []int {
	i := findLast(1, len(p), func(i int) bool { return p[i] > p[i-1] })
	if i == -1 {
		return nil
	}
	j := findLast(i, len(p), func(j int) bool { return p[j] > p[i-1] })
	p[j], p[i-1] = p[i-1], p[j]
	reverse(p[i:])
	return p
}

func findLast(start, end int, f func(int) bool) int {
	for i := end - 1; i >= start; i-- {
		if f(i) {
			return i
		}
	}
	return -1
}

func reverse(p []int) {
	for i := 0; i < len(p)/2; i++ {
		p[i], p[len(p)-i-1] = p[len(p)-i-1], p[i]
	}
}
