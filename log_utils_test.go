package log4g

import (
	"github.com/jrivets/log4g/Godeps/_workspace/src/github.com/jrivets/go-common/collections"
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
	ss, _ := collections.NewSortedSlice(2)
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

	ok = checkPanic(func() { getConfigParamName("appender..test", "appender", isCorrectAppenderName) })
	c.Assert(ok, Equals, true)

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

func (s *nameUtilsSuite) TestParseInt64(c *C) {
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

func (s *nameUtilsSuite) TestParseInt(c *C) {
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

func (s *nameUtilsSuite) TestParseBool(c *C) {
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

func (nus *nameUtilsSuite) name() string {
	return nus.loggerName
}

func (nus *nameUtilsSuite) Compare(other collections.Comparator) int {
	return compare(nus, other.(*nameUtilsSuite))
}

func checkPanic(f func()) (result bool) {
	defer func() {
		result = recover() != nil
	}()
	f()
	return
}
