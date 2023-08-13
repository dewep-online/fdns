package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var dot = byte('.')

func Level(s string, level int) string {
	if level == 0 {
		return "."
	}
	max := len(s) - 1
	count, pos := 0, 0
	if s[max] == dot {
		max--
	}

	for i := max; i >= 0; i-- {
		if s[i] == dot {
			count++
			if count == level {
				pos = i + 1
				break
			}
		}
	}
	return s[pos:]
}

func CountLevels(s string) int {
	var count int
	lastIndex := len(s) - 1
	for i := 0; i < len(s); i++ {
		if s[i] == dot {
			count++
			if i == lastIndex {
				count--
			}
		}
	}
	return count
}

var domainRex = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

func Normalize(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	if !domainRex.MatchString(domain) {
		return "", fmt.Errorf("invalid domain")
	}
	domain = strings.TrimRight(domain, ".")
	domain = strings.ToLower(domain)
	return domain + ".", nil
}
