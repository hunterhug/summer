package stool

import (
	"strconv"
)

func StringToInt(s string) int {
	if i, err := strconv.ParseInt(s, 10, 32); err == nil {
		return int(i)
	}
	return 0
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToString(i float64) string {
	return strconv.FormatFloat(i, 'f', -1, 64)
}

func StringToInt64(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	return 0
}

func SInt64(input string) int64 {
	result, _ := strconv.ParseInt(input, 10, 64)
	return result
}

func StringToFloat64(s string) float64 {
	if i, err := strconv.ParseFloat(s, 64); err == nil {
		return i
	}
	return 0.00
}

func StringToFloat32(s string) float32 {
	if i, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(i)
	}
	return 0.00
}

func StringToBool(s string) bool {
	if i, err := strconv.ParseBool(s); err == nil {
		return i
	}
	return false
}
