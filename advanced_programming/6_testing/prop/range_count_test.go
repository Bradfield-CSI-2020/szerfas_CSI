package prop

import (
	"fmt"
	"sort"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestRangeCount(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("binary range count matches linear", prop.ForAll(
		func(a []int, x int) string {
			sort.Ints(a)
			linearCount := RangeCountLinear(a, x)
			binaryCount := RangeCountBinary(a, x)
			if linearCount != binaryCount {
				return fmt.Sprintf(
					"RangeCountBinary(%v, %v): expected %d, got %d", a, x, linearCount, binaryCount)
			}
			return ""
		},
		gen.SliceOf(gen.IntRange(0, 10)),
		gen.IntRange(-20, 20),
	))

	properties.TestingRun(t)
}
