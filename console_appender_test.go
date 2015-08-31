package log4g

import (
	. "github.com/jrivets/log4g/Godeps/_workspace/src/gopkg.in/check.v1"
	"time"
)

type cAppenderSuite struct {
	msg    string
	signal chan bool
}

var _ = Suite(&cAppenderSuite{})

func (cas *cAppenderSuite) Write(p []byte) (n int, err error) {
	cas.msg = string(p)
	cas.signal <- true
	return len(p), nil
}

func (s *cAppenderSuite) TestNewAppender(c *C) {
	a, err := caFactory.NewAppender(map[string]string{})
	c.Assert(a, IsNil)
	c.Assert(err, NotNil)

	a, err = caFactory.NewAppender(map[string]string{"abcd": "1234"})
	c.Assert(a, IsNil)
	c.Assert(err, NotNil)

	a, err = caFactory.NewAppender(map[string]string{"layout": "%c %p"})
	c.Assert(a, NotNil)
	c.Assert(err, IsNil)
}

func (s *cAppenderSuite) TestAppend(c *C) {
	s.signal = make(chan bool, 1)
	caFactory.out = s

	a, _ := caFactory.NewAppender(map[string]string{"layout": "[%d{15:04:05.000}] %p %c: %m"})
	testTime := time.Unix(123456, 0)
	expectedTime := testTime.Format("[15:04:05.000]")
	appended := a.Append(&LogEvent{FATAL, testTime, "a.b.c", "Hello Console!"})
	c.Assert(appended, Equals, true)
	<-s.signal
	c.Assert(s.msg, Equals, expectedTime+" FATAL a.b.c: Hello Console!\n")

	caFactory.Shutdown()
	appended = a.Append(&LogEvent{FATAL, time.Unix(0, 0), "a.b.c", "Never delivered"})
	c.Assert(appended, Equals, false)
}
