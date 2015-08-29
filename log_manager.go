package log4g

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
)

type logManager struct {
	config *logConfig
	rwLock sync.RWMutex
}

var lm *logManager = &logManager{config: newLogConfig()}

func (lm *logManager) getLogger(loggerName string) Logger {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lm.config.initIfNeeded()
	return lm.config.getLogger(loggerName)
}

func (lm *logManager) setLogLevel(loggerName string, level Level) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lm.config.initIfNeeded()
	lm.config.setLogLevel(level, loggerName)
}

func (lm *logManager) registerAppender(appenderFactory AppenderFactory) error {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	return lm.config.registerAppender(appenderFactory)
}

func (lm *logManager) shutdown() {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	lm.config.cleanUp()
	for _, af := range lm.config.appenderFactorys {
		af.Shutdown()
	}
}

func (lm *logManager) setPropsFromFile(configFileName string) error {
	f, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	props := map[string]string{}
	scanner := bufio.NewScanner(f)
	lineNum := 1
	for ; scanner.Scan(); lineNum++ {
		line := strings.TrimLeft(scanner.Text(), " ")
		if len(line) == 0 || strings.Index(line, "#") == 0 {
			continue
		}

		idx := strings.Index(line, "=")
		if idx < 0 {
			return errors.New("Unexpected value in line " + strconv.Itoa(lineNum) + ": \"" + line +
				"\" in the config file " + configFileName +
				". key-value pair should be specified in the form <key>=<value>")
		}

		props[strings.TrimRight(line[:idx], " ")] = line[idx+1:]
	}

	return lm.setNewProperties(props)
}

func (lm *logManager) setNewProperties(props map[string]string) (err error) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	defer func() {
		p := recover()
		if p == nil {
			return
		}
		err = errors.New(p.(string))
	}()
	config := newLogConfig()
	oldConfig := lm.config
	config.initWithParams(oldConfig, props)

	lm.config = config
	oldConfig.cleanUp()
	return
}

func (lm *logManager) setLogLevelName(level int, name string) bool {
	if level < 0 || level >= len(lm.config.levelNames) {
		return false
	}
	lm.config.levelNames[level] = name
	return true
}
