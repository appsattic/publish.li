package main

import "strings"

var chars = "abcdefghijklmnopqrstuvwxyz"
var twitterChars = chars + "_"
var facebookChars = chars + "."
var githubChars = chars + "-"
var instagramChars = chars

func isValidTwitterHandle(handle string) bool {
	return strings.Trim(strings.ToLower(handle), twitterChars) == ""
}

func isValidFacebookHandle(handle string) bool {
	return strings.Trim(strings.ToLower(handle), facebookChars) == ""
}

func isValidGitHubHandle(handle string) bool {
	return strings.Trim(strings.ToLower(handle), githubChars) == ""
}

func isValidInstagramHandle(handle string) bool {
	return strings.Trim(strings.ToLower(handle), instagramChars) == ""
}
