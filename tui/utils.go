package tui

import (
	"regexp"
	"strings"
)

func splitBefore(s string, re *regexp.Regexp) (r []string) {
	re.ReplaceAllStringFunc(s, func(x string) string {
		s = strings.Replace(s, x, "::"+x, 1)
		return s
	})
	for _, x := range strings.Split(s, "::") {
		if x != "" {
			r = append(r, x)
		}
	}
	return
}
