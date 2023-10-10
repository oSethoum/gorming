package parser

import "strings"

func cleanString(s string, parts ...string) string {
	for _, part := range parts {
		s = strings.ReplaceAll(s, part, "")
	}
	return s
}
