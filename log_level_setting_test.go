package log4g

import (
	"github.com/jrivets/log4g/Godeps/_workspace/src/github.com/jrivets/go-common/collections"
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type logLevelSettingSuite struct {
}

var _ = Suite(&logLevelSettingSuite{})

func (s *logLevelSettingSuite) TestSetLogLevel(c *C) {
	ss, _ := collections.NewSortedSlice(2)
	setLogLevel(INFO, "a.b", ss)
	setLogLevel(INFO, "a", ss)
	setLogLevel(INFO, "b", ss)

	c.Assert(ss.Len(), Equals, 3)
	c.Assert(ss.At(0).(*logLevelSetting).loggerName, Equals, "a")
	c.Assert(ss.At(1).(*logLevelSetting).loggerName, Equals, "a.b")
	c.Assert(ss.At(2).(*logLevelSetting).loggerName, Equals, "b")
}

func (s *logLevelSettingSuite) TestGetSetLogLevel(c *C) {
	ss, _ := collections.NewSortedSlice(2)
	c.Assert(getLogLevelSetting("a", ss), IsNil)

	setLogLevel(INFO, "b", ss)
	c.Assert(getLogLevelSetting("a", ss), IsNil)

	setLogLevel(INFO, "", ss)
	c.Assert(getLogLevelSetting("a", ss).loggerName, Equals, "")

	setLogLevel(INFO, "b.d", ss)
	setLogLevel(INFO, "b.d.e", ss)
	setLogLevel(INFO, "b.d.e.g", ss)
	c.Assert(getLogLevelSetting("b.d.a", ss).loggerName, Equals, "b.d")
	c.Assert(getLogLevelSetting("b.d.e.f", ss).loggerName, Equals, "b.d.e")
	c.Assert(getLogLevelSetting("b.d.e.g", ss).loggerName, Equals, "b.d.e.g")
}
