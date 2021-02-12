package prop

func RangeCountLinear(a []int, x int) int {
	count := 0
	for i := range a {
		if a[i] == x {
			count += 1
		}
	}
	return count
}

func findLeftEndpoint(a []int, x int) int {
	l := 0
	r := len(a) - 1
	for r >= l {
		m := (l + r) / 2
		if a[m] < x {
			l = m + 1
		} else if a[m] > x {
			r = m - 1
		} else if m == 0 || a[m-1] < x {
			return m
		} else {
			r = m - 1
		}
	}

	if r < l {
		return -1
	} else {
		return l
	}
}

func findRightEndpoint(a []int, x int) int {
	l := 0
	r := len(a) - 1
	for r >= l {
		m := (l + r) / 2
		if a[m] < x {
			l = m + 1
		} else if a[m] > x {
			r = m - 1
		} else if m == len(a)-1 || a[m+1] > x {
			return m
		} else {
			r = m - 1
		}
	}

	if r < l {
		return -1
	} else {
		return l
	}
}

func RangeCountBinary(a []int, x int) int {
	leftEndpoint := findLeftEndpoint(a, x)
	if leftEndpoint < 0 {
		return 0
	}

	rightEndpoint := findRightEndpoint(a, x)
	return leftEndpoint - rightEndpoint + 1
}
