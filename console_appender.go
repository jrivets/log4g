package log4g

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const consoleAppenderName = "log4g/consoleAppender"

// layout - appender setting to specify format of the log event to message
// transformation
const CAParamLayout = "layout"

type consoleAppender struct {
	layoutTemplate LayoutTemplate
}

type consoleAppenderFactory struct {
	msgChannel chan string
	out        io.Writer
}

var caFactory *consoleAppenderFactory

func init() {
	caFactory = &consoleAppenderFactory{make(chan string, 1000), os.Stdout}

	err := RegisterAppender(caFactory)
	if err != nil {
		close(caFactory.msgChannel)
		fmt.Println("It is impossible to register console appender: ", err)
		panic(err)
	}
	go func() {
		for {
			str, ok := <-caFactory.msgChannel
			if !ok {
				break
			}
			fmt.Fprint(caFactory.out, str, "\n")
		}
	}()
}

func (*consoleAppenderFactory) Name() string {
	return consoleAppenderName
}

func (caf *consoleAppenderFactory) NewAppender(params map[string]string) (Appender, error) {
	layout, ok := params[CAParamLayout]
	if !ok || len(layout) == 0 {
		return nil, errors.New("Cannot create console appender without specified layout")
	}

	layoutTemplate, err := ParseLayout(layout)
	if err != nil {
		return nil, errors.New("Cannot create console appender: " + err.Error())
	}

	return &consoleAppender{layoutTemplate}, nil
}

func (caf *consoleAppenderFactory) Shutdown() {
	close(caf.msgChannel)
}

// Appender interface implementation
func (cAppender *consoleAppender) Append(event *LogEvent) (ok bool) {
	ok = false
	defer EndQuietly()
	msg := ToLogMessage(event, cAppender.layoutTemplate)
	caFactory.msgChannel <- msg
	ok = true
	return ok
}

func (cAppender *consoleAppender) Shutdown() {
	// Nothing should be done for the console appender
}
