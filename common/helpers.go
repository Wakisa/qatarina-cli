package common

import "strings"

func DefaultIfEmpty(value string) string {
	if strings.TrimSpace(value) == "" {
		return "N/A"
	}
	return value
}
