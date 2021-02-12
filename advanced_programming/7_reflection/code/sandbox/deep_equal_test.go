package sandbox

import (
	"testing"
)

func TestDeepEqual(t *testing.T) {
	type movie struct {
		year    int
		title   string
		sequels []*movie
	}

	for _, test := range []struct {
		a, b   interface{}
		expect bool
	}{
		{2, 5, false},
		{[]int{2, 3, 4}, []int{2, 3, 4}, true},
		{movie{1995, "Hackers", nil}, movie{1997, "Hackers", nil}, false},
		{map[string]int{"hello": 2, "foo": 3}, map[string]int{"hello": 2, "foo": 3}, true},
	} {
		if DeepEqual(test.a, test.b) != test.expect {
			t.Errorf("Wrong answer for DeepEqual(%v, %v)", test.a, test.b)
		}
	}
}
