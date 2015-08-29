package log4g

import "time"

// Level type represents logging level as an integer value which lies in [0..70] range.
// A level with lowest value has higher priority than a level with highest value.
// That means that if level X is set as maximum for logging, only messages with
// levels X1 <= X will be processed.
type Level int

// Predefined log levels constants. Users can define own ones or overwrite
// the predefined via configuration
const levelStep = 10
const (
	FATAL Level = levelStep*iota + levelStep
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
	ALL
)

// Logger interface provides various methods for sending messages to log4g
type Logger interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
	Log(level Level, args ...interface{})
	Logf(level Level, fstr string, args ...interface{})
	Logp(level Level, payload interface{})
}

// LogEvent is DTO, bearing a log message between the log4g components. This
// object contains information about the message which eventually should be
// delivered to the log storage(s) (final destinations of the message)
// through one or many log appenders
type LogEvent struct {
	Level      Level
	Timestamp  time.Time
	LoggerName string
	Payload    interface{}
}

// Appender is an interface for a log endpoint. Different log storages can be
// connected to the library by implementing the interface
type Appender interface {
	Append(event *LogEvent) bool
	// should be called every time when the instance is not going to be used anymore
	Shutdown()
}

// The factory allows to create an appender instances
type AppenderFactory interface {
	// Appender name
	Name() string
	NewAppender(map[string]string) (Appender, error)
	Shutdown()
}

// SetLogLevelName allows to associate level with its name. All messages with
// the level, which have been emitted after this settings, will appear with the
// provided name.
// returns false if the level is out of the acceptable range, or true if the
// name is applied.
func SetLogLevelName(level Level, name string) bool {
	return lm.setLogLevelName(int(level), name)
}

// GetLogger returns pointer to the Logger object for specified logger name.
// The function will always return the same pointer for the same logger's name
// regardless of log4g configuration or other settings
func GetLogger(loggerName string) Logger {
	return lm.getLogger(loggerName)
}

// SetLogLevel allows to set specified level for the provided logger name.
// Logger name is a string which should start from a letter, can contain letters [A-Za-z],
// digits[0-9] and dots '.' symbols. The name cannot have '.' as a last symbol
// in the logger's name.
//
// log4g uses tree-based structures to represent logger name model: Every logger name
// can be considered like a node in the tree where "" name is root of the tree.
// '.' has special meaning to separate the name nodes.
//
// So, for example, the logger name "FileSystem" can be represented as following
// tree:
//
//   "" (ROOT element)
//    \
//     \
//   FileSystem
//
// And if 2 more loggers are defined "FileSystem.ext3" and "FileSystem.ntfs"
// then the tree would be represented like:
//
//   "" (ROOT element)
//    \
//     \
//   FileSystem
//      /\
//     /  \
//   ext3 ntfs
//
// setting log level for "FileSystem" can also set the level for "FileSystem.ext3"
// and "FileSystem.ntfs", with one exception - if their names was not set before
// explicitly by the call of via configuration. In other words a logger inherits
// its log level from an ancestor in the logger name tree if its own log level
// was not set before!
func SetLogLevel(loggerName string, level Level) {
	lm.setLogLevel(loggerName, level)
}

// RegisterAppender allows to register an appender implementation in log4g.
// All appenders should register themself calling the function from init() or
// by calling this function directly.
// The function returns error if another factory has been registered for the
// same name before the call
// Parameters:
//		appenderFactory - a factory object which allows to create new instances of
//    	the appender type.
func RegisterAppender(appenderFactory AppenderFactory) error {
	return lm.registerAppender(appenderFactory)
}

// ConfigF reads log4g configuration properties from text file, which name is provided in
// configFileName parameter.
func ConfigF(configFileName string) error {
	return lm.setPropsFromFile(configFileName)
}

// Config allows to configure log4g by properties provided in the key:value form
func Config(props map[string]string) error {
	return lm.setNewProperties(props)
}

// Should be called to shutdown log subsystem properly. It will notify all logContexts and wait
// while all go routines that deliver messages to appenders are over. Calling this method could
// be essential to finalize some appenders and release their resources properly
func Shutdown() {
	lm.shutdown()
}
