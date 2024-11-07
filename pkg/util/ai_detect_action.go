package util

import (
	"regexp"
	"strings"
)

func DetectAction(text string) string {
	actionPattern := regexp.MustCompile(`(?i)action:\s*([\w_]+)[\s\)]*$`)
	match := actionPattern.FindStringSubmatch(strings.TrimSpace(text))

	// If a match is found, return the action name
	if len(match) > 1 {
		return match[1]
	}

	// If no action is detected, return an empty string
	return ""
}
