package log4g

import "github.com/jrivets/log4g/Godeps/_workspace/src/github.com/jrivets/go-common/collections"

type logLevelSetting struct {
	loggerName string
	level      Level
}

/**
 * Stores level setting for the provided loggerName
 * Params:
 *		loggerName - should be eligable normalized logger name
 */
func setLogLevel(level Level, loggerName string, logLevels *collections.SortedSlice) *logLevelSetting {
	if level < 0 {
		return nil
	}
	var lls *logLevelSetting = &logLevelSetting{loggerName, level}
	idx, found := logLevels.Find(lls)
	if found {
		lls = logLevels.At(idx).(*logLevelSetting)
		lls.level = level
	} else {
		logLevels.Add(lls)
	}
	return lls
}

func getLogLevelSetting(loggerName string, logLevels *collections.SortedSlice) *logLevelSetting {
	lProvider := getNearestAncestor(&logLevelSetting{loggerName: loggerName}, logLevels)
	if lProvider == nil {
		return nil
	}
	return lProvider.(*logLevelSetting)
}

// logNameProvider implementation
func (lls *logLevelSetting) name() string {
	return lls.loggerName
}

// Comparator implementation
func (lls *logLevelSetting) Compare(other collections.Comparator) int {
	return compare(lls, other.(*logLevelSetting))
}
