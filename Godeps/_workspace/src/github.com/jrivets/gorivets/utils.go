package gorivets

import (
	"errors"
	"strconv"
	"strings"
)

// The Comparator interface is used in collections to compare element if some ordering is needed
// please see `SortedSlice` as an example
type Comparator interface {
	Compare(Comparator) int
}

// Returns minimal value of two integers provided
func Min(a, b int) int {
	if a < b {
		return a
	} else if b < a {
		return b
	}
	return a
}

// Calls recover() to consume panic if it happened, it is recomended to be used with defer:
//
// 	func TestNoPanic() {
//		defer EndQuietly()
// 		panic()
//	}
func EndQuietly() {
	recover()
}

// Parses boolean value removing leading or trailing spaces if present
func ParseBool(value string, defaultValue bool) (bool, error) {
	value = strings.ToLower(strings.Trim(value, " "))
	if value == "" {
		return defaultValue, nil
	}

	return strconv.ParseBool(value)
}

// Parses string to int value, see ParseInt64
func ParseInt(value string, min, max, defaultValue int) (int, error) {
	res, err := ParseInt64(value, int64(min), int64(max), int64(defaultValue))
	return int(res), err
}

// ParseInt64 tries to convert value to int64, or returns default if the value is empty string.
// The value can look like "12Kb" which means 12*1000, or "12Kib" which means 12*1024. The following
// suffixes are supported for scaling by 1000 each: "kb", "mb", "gb", "tb", "pb", or "k", "m", "g", "t", "p".
// For 1024 scale the following suffixes are supported: "kib", "mib", "gib", "tib", "pib"
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
