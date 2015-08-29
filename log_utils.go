package log4g

import (
	"errors"
	"github.com/jrivets/log4g/Godeps/_workspace/src/github.com/jrivets/go-common/collections"
	"regexp"
	"strconv"
	"strings"
)

const maxInt64 = 1<<63 - 1

type logNameProvider interface {
	name() string
}

func compare(n1, n2 logNameProvider) int {
	switch {
	case n1.name() == n2.name():
		return 0
	case n1.name() < n2.name():
		return -1
	}
	return 1
}

// loggerName cannot start/end from spaces and dots
func normalizeLogName(name string) string {
	return strings.Trim(name, ". ")
}

/**
 * Checks whether the checkedName is ancestor for the loggerName or not
 * The name checkedName is ancestor for the loggerName if:
 *	- checkedName == loggerName
 *  - loggerName == checkedName.<some name here>
 * 	- checkedName == rootLoggerName
 */
func ancestor(checkedName, loggerName string) bool {
	if checkedName == loggerName || checkedName == rootLoggerName {
		return true
	}

	lenc := len(checkedName)
	lenl := len(loggerName)
	if strings.HasPrefix(loggerName, checkedName) && lenl > lenc && loggerName[lenc] == '.' {
		return true
	}
	return false
}

func getNearestAncestor(comparator collections.Comparator, names *collections.SortedSlice) logNameProvider {
	if names.Len() == 0 {
		return nil
	}
	nProvider := comparator.(logNameProvider)
	for idx := Min(names.Len()-1, names.GetInsertPos(nProvider.(collections.Comparator))); idx >= 0; idx-- {
		candidate := names.At(idx).(logNameProvider)
		if ancestor(candidate.name(), nProvider.name()) {
			return candidate
		}
	}
	return nil
}

// Gets the name of a parameter provided in the form: <prefix>.<name>.<attribute>
func getConfigParamName(param, prefix string, checker func(string) bool) (string, bool) {
	pr := prefix + "."
	if !strings.HasPrefix(param, pr) {
		return "", false
	}

	start := len(pr)
	end := strings.LastIndex(param, ".")
	if start == end+1 {
		return "", true
	}

	paramName := param[start:end]
	if checker != nil && !checker(paramName) {
		panic("Unacceptable param value \"" + paramName + "\" for " + prefix + " setting.")
	}

	return paramName, true
}

// Gets the attribute of a parameter provided in the form: <prefix>.<name>.<attribute>
func getConfigParamAttribute(param string) string {
	end := strings.LastIndex(param, ".")
	if end == len(param)-1 {
		return ""
	}
	return param[end+1:]
}

// Groups params with the prefix by their names into a map of maps, where the second
// map defines parameters for the particular key value (param name) from the first map
func groupConfigParams(params map[string]string, prefix string, checker func(string) bool) map[string]map[string]string {
	resultMap := make(map[string]map[string]string)
	for k, v := range params {
		name, ok := getConfigParamName(k, prefix, checker)
		if !ok {
			continue
		}
		attribute := getConfigParamAttribute(k)

		m, ok := resultMap[name]
		if !ok {
			m = make(map[string]string)
			resultMap[name] = m
		}
		m[attribute] = v
	}
	return resultMap
}

func isCorrectAppenderName(appenderName string) bool {
	matched, err := regexp.MatchString("^[A-Za-z][A-Za-z0-9.]+$", appenderName)
	if !matched || err != nil {
		return false
	}
	return true
}

func isCorrectLoggerName(loggerName string) bool {
	if loggerName == "" {
		return true
	}
	matched, err := regexp.MatchString("^[A-Za-z]+([A-Za-z0-9.]*[A-Za-z0-9]+)*$", loggerName)
	if !matched || err != nil {
		return false
	}
	return true
}

func ParseBool(value string, defaultValue bool) (bool, error) {
	value = strings.ToLower(strings.Trim(value, " "))
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseBool(value)
}

func ParseInt(value string, min, max, defaultValue int) (int, error) {
	res, err := ParseInt64(value, int64(min), int64(max), int64(defaultValue))
	return int(res), err
}

// ParseInt tries to convert value to int64, or returns default if the value is empty string
func ParseInt64(value string, min, max, defaultValue int64) (int64, error) {
	if defaultValue < min || defaultValue > max || max < min {
		return 0, errors.New("Inconsistent arguments provided min=" + strconv.FormatInt(min, 10) +
			", max=" + strconv.FormatInt(max, 10) + ", defaultVelue=" + strconv.FormatInt(defaultValue, 10))
	}
	value = strings.ToLower(strings.Trim(value, " "))
	if value == "" {
		return defaultValue, nil
	}

	value, scale := parseSuffixVsScale(value, []string{"kb", "mb", "gb", "tb", "pb"}, 1000)
	if scale == 1 {
		value, scale = parseSuffixVsScale(value, []string{"k", "m", "g", "t", "p"}, 1000)
		if scale == 1 {
			value, scale = parseSuffixVsScale(value, []string{"kib", "mib", "gib", "tib", "pib"}, 1024)
		}
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	val := int64(intValue) * scale

	if min > val || max < val {
		return 0, errors.New("Value should be in the range [" + strconv.FormatInt(min, 10) + ".." + strconv.FormatInt(max, 10) + "]")
	}

	return val, nil
}

func parseSuffixVsScale(value string, suffixes []string, scale int64) (string, int64) {
	idx, str := getSuffix(value, suffixes)
	if idx < 0 {
		return value, 1
	}
	val := scale
	for ; idx > 0; idx-- {
		val *= scale
	}
	return value[:len(value)-len(str)], val
}

func getSuffix(value string, suffixes []string) (int, string) {
	for idx, sfx := range suffixes {
		if strings.HasSuffix(value, sfx) {
			return idx, sfx
		}
	}
	return -1, ""
}

// Utility methods
func Min(a, b int) int {
	if a < b {
		return a
	} else if b < a {
		return b
	}
	return a
}

func EndQuietly() {
	recover()
}
