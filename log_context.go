package log4g

import (
	"errors"
	"github.com/jrivets/log4g/Godeps/_workspace/src/github.com/jrivets/gorivets"
	"strconv"
)

type logContext struct {
	loggerName string
	appenders  []Appender
	inherited  bool
	blocking   bool
	eventsCh   chan *Event
	controlCh  chan bool
}

func newLogContext(loggerName string, appenders []Appender, inherited, blocking bool, bufSize int) (*logContext, error) {
	if bufSize <= 0 {
		return nil, errors.New("Cannot create channel with non-positive size=" + strconv.Itoa(bufSize))
	}

	if appenders == nil || len(appenders) == 0 {
		return nil, errors.New("At least one appender should be in Log Context")
	}

	eventsCh := make(chan *Event, bufSize)
	controlCh := make(chan bool, 1)
	lc := &logContext{loggerName, appenders, inherited, blocking, eventsCh, controlCh}

	go func() {
		defer onStop(controlCh)
		for {
			le, ok := <-eventsCh
			if !ok {
				break
			}
			lc.onEvent(le)
		}
	}()
	return lc, nil
}

// Processing go routine finalizer
func onStop(controlCh chan bool) {
	controlCh <- true
	close(controlCh)
}

func getLogLevelContext(loggerName string, logContexts *gorivets.SortedSlice) *logContext {
	lProvider := getNearestAncestor(&logContext{loggerName: loggerName}, logContexts)
	if lProvider == nil {
		return nil
	}
	return lProvider.(*logContext)
}

// log() function sends the logEvent to all the logContext appenders.
// It returns true if the logEvent was sent and false if the context is shut down or
// the context is non-blocking (allows to lost log messages in case of overflow)
func (lc *logContext) log(le *Event) (result bool) {
	// Channel can be already closed, so end quietly
	result = false
	defer gorivets.EndQuietly()

	if lc.blocking {
		lc.eventsCh <- le
		return true
	}

	select {
	case lc.eventsCh <- le:
		return true
	default:
	}
	return false
}

// Called from processing go routine
func (lc *logContext) onEvent(le *Event) {
	appenders := lc.appenders
	if len(appenders) == 1 {
		appenders[0].Append(le)
		return
	}
	for _, a := range appenders {
		a.Append(le)
	}
}

func (lc *logContext) shutdown() {
	close(lc.eventsCh)
	<-lc.controlCh
}

// logNameProvider implementation
func (lc *logContext) name() string {
	return lc.loggerName
}

// Comparable implementation
func (lc *logContext) Compare(other gorivets.Comparable) int {
	return compare(lc, other.(*logContext))
}
