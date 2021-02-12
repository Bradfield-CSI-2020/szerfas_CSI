package sandbox

import (
	"fmt"
	"testing"
)

func TestJSON(t *testing.T) {
	type movie struct {
		year            int
		title, director string
		awards          []string
	}

	m := movie{
		1995,
		"Hackers",
		"???",
		[]string{"best movie of all time", "hack the planet prize"},
	}

	fmt.Println(ToJSON(m))
}
