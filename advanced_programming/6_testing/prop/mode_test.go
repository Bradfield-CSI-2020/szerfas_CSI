package prop

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

type testCase struct {
	a []int
	m int
}

/*
func TestModeUnit(t *testing.T) {
	for _, test := range []testCase{
		{
			a: []int{1, 1, 2, 3},
			m: 1,
		},
		{
			a: []int{1, 2, 3, 4, 4},
			m: 4,
		},
		{
			a: []int{1, 2, 3, 4},
			m: 1,
		},
	} {
		actual := Mode(test.a)
		if actual != test.m {
			t.Errorf("Mode(%v) should be %v, got %v", test.a, test.m, actual)
		}
	}
}
*/

func TestMode(t *testing.T) {
	properties := gopter.NewProperties(nil) // nil = default config (e.g. 100 runs)

	properties.Property("mode is most frequently occurring element", prop.ForAll(
		func(a []int) string {
			m := Mode(a)

			// For each value in the input slice, keep track of
			// how many times it appears
			counts := make(map[int]int)
			for _, x := range a {
				counts[x] += 1
			}

			// For each value that appears in the input slice,
			// it better not appear more frequently than the "mode"
			// (which is by definition the one that appears the MOST
			// times)
			for x := range counts {
				if counts[x] > counts[m] {
					return fmt.Sprintf("Mode(%v) != %v: %v is more frequent", a, m, x)
				}
			}

			return ""
		},

		gen.SliceOf(gen.IntRange(0, 10)),
	))

	properties.TestingRun(t)
}
