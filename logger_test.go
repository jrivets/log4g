package log4g

import (
	"time"

	. "gopkg.in/check.v1"
)

type loggerSuite struct {
	loggerName string
}

var _ = Suite(&loggerSuite{})

func (s *loggerSuite) TestApplyNewLevelToLoggers(c *C) {
	rootLLS := &logLevelSetting{rootLoggerName, INFO}

	loggers := make(map[string]*logger)
	loggers["a"] = &logger{"a", rootLLS, nil, INFO}
	loggers["a.b"] = &logger{"a.b", rootLLS, nil, INFO}
	loggers["a.b.c"] = &logger{"a.b.c", rootLLS, nil, INFO}
	loggers["a.b.c.d"] = &logger{"a.b.c.d", rootLLS, nil, INFO}

	applyNewLevelToLoggers(&logLevelSetting{"a.b", DEBUG}, loggers)
	c.Assert(loggers["a"].logLevel, Equals, INFO)
	c.Assert(loggers["a.b"].logLevel, Equals, DEBUG)
	c.Assert(loggers["a.b.c"].logLevel, Equals, DEBUG)
	c.Assert(loggers["a.b.c.d"].logLevel, Equals, DEBUG)

	applyNewLevelToLoggers(&logLevelSetting{"a.b.c", TRACE}, loggers)
	applyNewLevelToLoggers(&logLevelSetting{"a.b", ERROR}, loggers)
	c.Assert(loggers["a"].logLevel, Equals, INFO)
	c.Assert(loggers["a.b"].logLevel, Equals, ERROR)
	c.Assert(loggers["a.b.c"].logLevel, Equals, TRACE)
	c.Assert(loggers["a.b.c.d"].logLevel, Equals, TRACE)
}

func (s *loggerSuite) TestLog(c *C) {
	lctx := &logContext{eventsCh: make(chan *Event, 1)}
	l := wrap(&logger{"a", nil, lctx, INFO})
	l.Log(INFO, "Hello")
	go waitThenClose(500, lctx)
	le, ok := <-lctx.eventsCh

	c.Assert(ok, Equals, true)
	c.Assert(le.Payload.(string), Equals, "Hello")
	c.Assert(le.Level, Equals, INFO)
	c.Assert(le.LoggerName, Equals, "a")
}

func (s *loggerSuite) TestLogDisabled(c *C) {
	lctx := &logContext{eventsCh: make(chan *Event, 1)}
	l := wrap(&logger{"a", nil, lctx, INFO})
	l.Log(DEBUG, "Hello")
	go waitThenClose(50, lctx)
	_, ok := <-lctx.eventsCh
	c.Assert(ok, Equals, false)
}

func (s *loggerSuite) TestLogf(c *C) {
	lctx := &logContext{eventsCh: make(chan *Event, 2)}
	l := wrap(&logger{"a", nil, lctx, INFO})
	l.Logf(INFO, "Hello")
	l.Logf(INFO, "Hello %s", "World!")
	c.Assert(l.GetName(), Equals, "a")
	go waitThenClose(500, lctx)
	le, ok := <-lctx.eventsCh
	c.Assert(ok, Equals, true)
	c.Assert(le.Payload.(string), Equals, "Hello")

	le, ok = <-lctx.eventsCh
	c.Assert(ok, Equals, true)
	c.Assert(le.Payload.(string), Equals, "Hello World!")
}

func (s *loggerSuite) TestLogp(c *C) {
	lctx := &logContext{eventsCh: make(chan *Event, 1)}
	l := wrap(&logger{"a", nil, lctx, INFO})
	l.Logp(INFO, lctx)
	go waitThenClose(500, lctx)
	le, ok := <-lctx.eventsCh
	c.Assert(ok, Equals, true)
	c.Assert(le.Payload, Equals, lctx)
}

func (s *loggerSuite) TestMessages(c *C) {
	lctx := &logContext{eventsCh: make(chan *Event, 10)}
	l := wrap(&logger{"a", nil, lctx, TRACE})
	l.Info(INFO)
	l.Warn(WARN)
	l.Debug(DEBUG)
	l.Trace(TRACE)
	l.Error(ERROR)
	l.Fatal(FATAL)
	go waitThenClose(500, lctx)
	c.Assert((<-lctx.eventsCh).Level, Equals, INFO)
	c.Assert((<-lctx.eventsCh).Level, Equals, WARN)
	c.Assert((<-lctx.eventsCh).Level, Equals, DEBUG)
	c.Assert((<-lctx.eventsCh).Level, Equals, TRACE)
	c.Assert((<-lctx.eventsCh).Level, Equals, ERROR)
	c.Assert((<-lctx.eventsCh).Level, Equals, FATAL)
}

func (s *loggerSuite) TestLevelsOrder(c *C) {
	ok := FATAL < ERROR && ERROR < WARN && WARN < INFO && INFO < DEBUG && DEBUG < TRACE
	c.Assert(ok, Equals, true)
}

func (s *loggerSuite) TestWithId(c *C) {
	lctx := &logContext{eventsCh: make(chan *Event, 10)}
	l := wrap(&logger{"a", nil, lctx, TRACE})
	c.Assert(l.loggerId, IsNil)
	l2 := l.WithId(123).(*log_wrapper)
	c.Assert(l2.loggerId, Equals, 123)

	l3 := l2.WithId(123).(*log_wrapper)
	c.Assert(l3, Equals, l2)
}

func waitThenClose(mSec time.Duration, lctx *logContext) {
	time.Sleep(time.Millisecond * mSec)
	close(lctx.eventsCh)
}
