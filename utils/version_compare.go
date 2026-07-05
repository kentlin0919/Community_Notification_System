package utils

import (
	"strconv"
	"strings"
)

// CompareVersions 比較兩個以點分隔的版本字串；a<b 回傳 -1，相等回傳 0，a>b 回傳 1
func CompareVersions(a, b string) int {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}

	for i := 0; i < maxLen; i++ {
		numA := versionPartToInt(partsA, i)
		numB := versionPartToInt(partsB, i)
		if numA < numB {
			return -1
		}
		if numA > numB {
			return 1
		}
	}
	return 0
}

func versionPartToInt(parts []string, index int) int {
	if index >= len(parts) {
		return 0
	}
	value, err := strconv.Atoi(parts[index])
	if err != nil {
		return 0
	}
	return value
}
