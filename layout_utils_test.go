package log4g

import (
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
	"time"
)

type layoutUtilsSuite struct {
}

var _ = Suite(&layoutUtilsSuite{})

func (s *layoutUtilsSuite) TestParseLayoutFail(c *C) {
	t, err := ParseLayout("%A")
	c.Assert(t, IsNil)
	c.Assert(err, NotNil)

	t, err = ParseLayout("%d(123)")
	c.Assert(t, IsNil)
	c.Assert(err, NotNil)

	t, err = ParseLayout("%dd{1234}")
	c.Assert(t, IsNil)
	c.Assert(err, NotNil)

	t, err = ParseLayout("%d(123)")
	c.Assert(t, IsNil)
	c.Assert(err, NotNil)

	t, err = ParseLayout("%d{123")
	c.Assert(t, IsNil)
	c.Assert(err, NotNil)

	t, err = ParseLayout("%")
	c.Assert(t, IsNil)
	c.Assert(err, NotNil)
}

func (s *layoutUtilsSuite) TestParseLayout(c *C) {
	t, _ := ParseLayout(" ")
	c.Assert(len(t), Equals, 1)
	c.Assert(t[0].value, Equals, " ")
	c.Assert(t[0].pieceType, Equals, lpText)

	t, _ = ParseLayout("Text")
	c.Assert(len(t), Equals, 1)
	c.Assert(t[0].value, Equals, "Text")
	c.Assert(t[0].pieceType, Equals, lpText)

	t, _ = ParseLayout("%c%p%d{123}%m%%")
	c.Assert(len(t), Equals, 5)
	c.Assert(t[0].pieceType, Equals, lpLoggerName)
	c.Assert(t[1].pieceType, Equals, lpLogLevel)
	c.Assert(t[2].value, Equals, "123")
	c.Assert(t[2].pieceType, Equals, lpDate)
	c.Assert(t[3].pieceType, Equals, lpMessage)
	c.Assert(t[4].pieceType, Equals, lpText)

	t, _ = ParseLayout("%c %p %d{123} %m %% %% ")
	c.Assert(len(t), Equals, 10)
}

func (s *layoutUtilsSuite) TestToLogMessage(c *C) {
	t, _ := ParseLayout("[%d{01-02 15:04:05.000}] %p %c: %%%m")
	le := &LogEvent{FATAL, time.Unix(123456, 0), "a.b.c", "The Message"}
	c.Assert(ToLogMessage(le, t), Equals, "[01-02 02:17:36.000] FATAL a.b.c: %The Message")
}
