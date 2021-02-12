package prop

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestMaxProduct(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("MaxProduct matches slow implementation", prop.ForAll(
		func(a []int32) string {
			expect := MaxProductSlow(a)
			actual := MaxProduct(a)
			if expect != actual {
				return fmt.Sprintf("MaxProduct(%v): expected %v, actual %v", a, expect, actual)
			}
			return ""
		},
		gen.SliceOf(gen.Int32()).SuchThat(func(a []int32) bool {
			return len(a) >= 2
		}),
	))

	properties.TestingRun(t)
}
