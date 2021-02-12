package hello_test

import (
	"fmt"
	"testing"
	"io"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{
	dir string
}

var _ = Suite(&MySuite{})

func (s *MySuite) TestHelloWorld(c *C) {
	c.Assert(42, Equals, "42")
	c.Assert(io.ErrClosedPipe, ErrorMatches, "io: .*on closed pipe")
	c.Check(42, Equals, 42)
}

func (s *MySuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
	// Use s.dir to prepare some data.
}

func (s *MySuite) TestWithDir(c *C) {
	// Use the data in s.dir in the test.
	fmt.Printf("s.dir: %s\n", s.dir)
	c.Assert(s.dir, NotNil)
}