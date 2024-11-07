package util

import "regexp"

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidGoogleSheetLink(link string) bool {
	linkRegex := regexp.MustCompile(`^https://docs\.google\.com/spreadsheets/(?:u/\d+/)?d/[a-zA-Z0-9_-]+/edit.*$`)
	return linkRegex.MatchString(link)
}
