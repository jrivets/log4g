package log4g

import (
	"github.com/jrivets/log4g/Godeps/_workspace/src/github.com/jrivets/gorivets"
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
)

type nameUtilsSuite struct {
	loggerName string
}

var _ = Suite(&nameUtilsSuite{})

func (s *nameUtilsSuite) TestAncestor(c *C) {
	c.Assert(ancestor("", ""), Equals, true)
	c.Assert(ancestor("a", "a"), Equals, true)
	c.Assert(ancestor("a.b", "a.b.c"), Equals, true)
	c.Assert(ancestor("a.b", "a.b.cd.e"), Equals, true)
	c.Assert(ancestor("a.b", "a.c.c"), Equals, false)
}

func (s *nameUtilsSuite) TestGetSetLogLevel(c *C) {
	ss, _ := gorivets.NewSortedSlice(2)
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a"}, ss), IsNil)

	ss.Add(&nameUtilsSuite{"b"})
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a"}, ss), IsNil)

	ss.Add(&nameUtilsSuite{""})
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a"}, ss).(*nameUtilsSuite).loggerName, Equals, "")

	ss.Add(&nameUtilsSuite{"a.b.c"})
	ss.Add(&nameUtilsSuite{"a.b"})
	ss.Add(&nameUtilsSuite{"a"})
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.c"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b.c")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.c.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b.c")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.b.c.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a.b.c")
	c.Assert(getNearestAncestor(&nameUtilsSuite{"a.bc.d"}, ss).(*nameUtilsSuite).loggerName, Equals, "a")
}

func (s *nameUtilsSuite) TestGetConfigParamName(c *C) {
	ctx, ok := getConfigParamName("abc", "context", nil)
	c.Assert(ok, Equals, false)

	ctx, ok = getConfigParamName("abc.asd.ab", "context", nil)
	c.Assert(ok, Equals, false)

	ctx, ok = getConfigParamName("context.test", "context", nil)
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "")

	ctx, ok = getConfigParamName("context..test", "context", nil)
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "")

	panicTest := gorivets.CheckPanic(func() { getConfigParamName("appender..test", "appender", isCorrectAppenderName) })
	c.Assert(panicTest, NotNil)

	ctx, ok = getConfigParamName("context...test", "context", nil)
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, ".")

	ctx, ok = getConfigParamName("context.test.text", "context", nil)
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "test")

	ctx, ok = getConfigParamName("context.a.b.c.test", "context", nil)
	c.Assert(ok, Equals, true)
	c.Assert(ctx, Equals, "a.b.c")
}

func (s *nameUtilsSuite) TestGetConfigParamAttribute(c *C) {
	c.Assert(getConfigParamAttribute("appender."), Equals, "")
	c.Assert(getConfigParamAttribute("appender.ROOT.level"), Equals, "level")
}

func (s *nameUtilsSuite) TestGroupConfigParams(c *C) {
	params := groupConfigParams(map[string]string{
		"context.ROOT.type": "123",
		"abc":               "def",
		"context.app.type":  "345",
		"context.ROOT.ttt":  "qqq",
	}, "context", nil)
	c.Assert(params["ROOT"]["type"], Equals, "123")
	c.Assert(params["ROOT"]["ttt"], Equals, "qqq")
	c.Assert(params["app"]["type"], Equals, "345")
	c.Assert(params["app"]["ttt"], Equals, "")
	c.Assert(params["abc"], IsNil)
}

func (s *nameUtilsSuite) TestCorrectAppenderName(c *C) {
	c.Assert(isCorrectAppenderName(""), Equals, false)
	c.Assert(isCorrectAppenderName("AbcL"), Equals, true)
	c.Assert(isCorrectAppenderName("abC1"), Equals, true)
	c.Assert(isCorrectAppenderName("2abC1"), Equals, false)
	c.Assert(isCorrectAppenderName("ad,CD"), Equals, false)
}

func (s *nameUtilsSuite) TestCorrectLoggerName(c *C) {
	c.Assert(isCorrectLoggerName(""), Equals, true)
	c.Assert(isCorrectLoggerName("a"), Equals, true)
	c.Assert(isCorrectLoggerName("a1"), Equals, true)
	c.Assert(isCorrectLoggerName("1a"), Equals, false)
	c.Assert(isCorrectLoggerName("a.b.c"), Equals, true)
	c.Assert(isCorrectLoggerName(".a"), Equals, false)
	c.Assert(isCorrectLoggerName("a."), Equals, false)
}

func (nus *nameUtilsSuite) name() string {
	return nus.loggerName
}

func (nus *nameUtilsSuite) Compare(other gorivets.Comparable) int {
	return compare(nus, other.(*nameUtilsSuite))
}
