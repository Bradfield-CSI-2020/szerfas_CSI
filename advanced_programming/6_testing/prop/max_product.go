package prop

import (
	"math"
	"sort"
)

type Int32Slice []int32

func (a Int32Slice) Len() int           { return len(a) }
func (a Int32Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Int32Slice) Less(i, j int) bool { return a[i] < a[j] }

func MaxProductSlow(a []int32) int64 {
	n := len(a)
	if n < 2 {
		panic("Input must have at least 2 elements")
	}

	var ans int64 = math.MinInt64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p := int64(a[i]) * int64(a[j])
			if p > ans {
				ans = p
			}
		}
	}
	return ans
}

func MaxProduct(a []int32) int64 {
	n := len(a)
	if n < 2 {
		panic("Input must have at least 2 elements")
	}

	sort.Sort(Int32Slice(a))
	return int64(a[n-1]) * int64(a[n-2])
}
