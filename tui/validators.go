package tui

import "regexp"

func valIsBlank(s string) (string, bool) {
	if s == "" {
		return "Field must not be blank", true
	}
	return "", false
}

func valIsntLcAndNum(s string) (string, bool) {
	lcAndNum := regexp.MustCompile("^[a-z0-9]*$")
	if !lcAndNum.MatchString(s) {
		return "May only include numbers and lowercase letters", true
	}
	return "", false
}

func valIsntLcNumUnderDash(s string) (string, bool) {
	lcNumUnderScoreDash := regexp.MustCompile("^[a-z0-9_-]*$")
	if !lcNumUnderScoreDash.MatchString(s) {
		return "May only include numbers, lowercase letters, underscores, and dashes", true
	}
	return "", false
}

func valIsntNum(s string) (string, bool) {
	num := regexp.MustCompile("^[0-9]*$")
	if !num.MatchString(s) {
		return "May only include numbers", true
	}
	return "", false
}

func valMultiSlash(s string) (string, bool) {
	slashes := regexp.MustCompile("/{2,}").MatchString(s)
	if slashes {
		return "Can only have one slash - multifractions not supported", true
	}
	return "", false
}

func valIsntTimeUnit(s string) (string, bool) {
	acceptedChars := regexp.MustCompile("^[0-9ywdhms]*$")
	if !acceptedChars.MatchString(s) {
		return "May only include numbers and y, w, d, h, m, and s", true
	}
	return "", false
}

func valRepeatLetters(s string) (string, bool) {
	repChar := regexp.MustCompile("[a-zA-Z]{2,}")
	if repChar.MatchString(s) {
		return "Cannot include multiple characters in a row", true
	}
	return "", false
}

func valStartsWithChar(s string) (string, bool) {
	startsWChar := regexp.MustCompile("^[a-zA-Z]")
	if startsWChar.MatchString(s) {
		return "Cannot start with a character", true
	}
	return "", false
}

func valNoTimeUnit(s string) (string, bool) {
	timeChars := regexp.MustCompile("[ywdhms]")
	if !timeChars.MatchString(s) {
		return "Needs some unit of time (y, w, d, h, m, s)", true
	}
	return "", false
}

// TODO

// 10uM/mg/mg should fail
