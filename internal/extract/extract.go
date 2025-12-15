package extract

import (
	"fmt"
	"regexp"
)

func ByRegex(b []byte, pattern string, group int) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("compile regex: %w", err)
	}

	m := re.FindSubmatch(b)
	if m == nil {
		return "", fmt.Errorf("pattern did not match")
	}

	if group < 0 || group >= len(m) {
		return "", fmt.Errorf("group %d out of range (matches=%d)", group, len(m)-1)
	}

	return string(m[group]), nil
}