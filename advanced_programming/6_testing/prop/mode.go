package prop

func Mode(a []int) int {
	best := 0
	counts := make(map[int]int)
	for _, x := range a {
		counts[x] += 1
		if best == 0 || counts[x] > counts[best] {
			best = x
		}
	}
	return best
}
