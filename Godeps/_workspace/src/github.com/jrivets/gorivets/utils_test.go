package gorivets

import (
	"errors"
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
)

type utilsSuite struct {
}

var _ = Suite(&utilsSuite{})

func (s *utilsSuite) TestNoPanic(c *C) {
	defer EndQuietly()
	panic("")
}

func (s *utilsSuite) TestMin(c *C) {
	c.Assert(Min(-1, 0), Equals, -1)
	c.Assert(Min(1, 0), Equals, 0)
	c.Assert(Min(2, 2), Equals, 2)
	c.Assert(Min(10, 100), Equals, 10)
}

func (s *utilsSuite) TestParseInt64(c *C) {
	_, err := ParseInt64("123", 1, 220, 0)
	c.Assert(err, NotNil)
	_, err = ParseInt64("123", 1, 220, 221)
	c.Assert(err, NotNil)

	v, err := ParseInt64("123", 1, 220, 10)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, int64(123))

	v, err = ParseInt64("", 1, 220, 10)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, int64(10))

	v, err = ParseInt64("k", 1, 220, 10)
	c.Assert(err, NotNil)

	v, err = ParseInt64("1k", 1, 220, 10)
	c.Assert(err, NotNil)

	v, err = ParseInt64("1k", 1, 2200, 10)
	c.Assert(v, Equals, int64(1000))

	v, err = ParseInt64("1Mb", 1, 2200000, 10)
	c.Assert(v, Equals, int64(1000000))

	v, err = ParseInt64("1MiB", 1, 2200000, 10)
	c.Assert(v, Equals, int64(1024*1024))
	c.Assert(int(v), Equals, 1024*1024)
}

func (s *utilsSuite) TestParseInt(c *C) {
	_, err := ParseInt("123", 1, 220, 0)
	c.Assert(err, NotNil)
	_, err = ParseInt("123", 1, 220, 221)
	c.Assert(err, NotNil)

	v, err := ParseInt("123", 1, 220, 10)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, 123)

	v, err = ParseInt("", 1, 220, 10)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, 10)

	v, err = ParseInt("1k", 1, 2200, 10)
	c.Assert(v, Equals, 1000)
}

func (s *utilsSuite) TestParseBool(c *C) {
	_, err := ParseBool("123", true)
	c.Assert(err, NotNil)

	v, err := ParseBool("true", true)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, true)

	v, err = ParseBool("false", true)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, false)

	v, err = ParseBool("", false)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, false)

	v, err = ParseBool("", true)
	c.Assert(err, IsNil)
	c.Assert(v, Equals, true)
}

func (s *utilsSuite) TestAssertNotNull(c *C) {
	c.Assert(CheckPanic(func() { AssertNoError(nil) }), Equals, false)
	c.Assert(CheckPanic(func() { AssertNoError(errors.New("ddd")) }), Equals, true)
}
