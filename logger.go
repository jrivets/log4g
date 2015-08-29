package log4g

import (
	"fmt"
	"time"
)

type logger struct {
	loggerName string
	lls        *logLevelSetting
	lctx       *logContext
	logLevel   Level
}

func (l *logger) Fatal(args ...interface{}) {
	l.Log(FATAL, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Log(ERROR, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.Log(WARN, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Log(INFO, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Log(DEBUG, args...)
}

func (l *logger) Trace(args ...interface{}) {
	l.Log(TRACE, args...)
}

func (l *logger) Log(level Level, args ...interface{}) {
	if l.logLevel < level {
		return
	}
	l.logInternal(level, fmt.Sprint(args...))
}

func (l *logger) Logf(level Level, fstr string, args ...interface{}) {
	if l.logLevel < level {
		return
	}
	msg := fstr
	if len(args) > 0 {
		msg = fmt.Sprintf(fstr, args...)
	}
	l.logInternal(level, msg)
}

func (l *logger) Logp(level Level, payload interface{}) {
	if l.logLevel < level {
		return
	}
	l.logInternal(level, payload)
}

func (l *logger) logInternal(level Level, payload interface{}) {
	l.lctx.log(&LogEvent{level, time.Now(), l.loggerName, payload})
}

func (l *logger) setLogLevelSetting(lls *logLevelSetting) {
	l.lls = lls
	l.logLevel = lls.level
}

func (l *logger) setLogContext(lctx *logContext) {
	l.lctx = lctx
}

// Apply new LogLevelSetting to all appropriate loggers
func applyNewLevelToLoggers(lls *logLevelSetting, loggers map[string]*logger) {
	for _, l := range loggers {
		if !ancestor(lls.loggerName, l.loggerName) {
			continue
		}
		if ancestor(l.lls.loggerName, lls.loggerName) {
			l.setLogLevelSetting(lls)
		}
	}
}
