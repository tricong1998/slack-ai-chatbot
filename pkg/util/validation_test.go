package util

import "testing"

func TestIsValidGoogleSheetLink(t *testing.T) {
	link := "https://docs.google.com/spreadsheets/u/0/d/1hbWo3suJYJYNIfWhj0oVT6E03qA0kOwhnIRwxMMNXaY/edit?gid=0&pli=1&authuser=0#gid=0"
	if !IsValidGoogleSheetLink(link) {
		t.Errorf("Expected %s to be a valid Google Sheet link", link)
	}
}
