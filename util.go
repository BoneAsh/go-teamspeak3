package teamspeak3

import "regexp"

func ReplaceAnsiEscapeCode(text string) string {
	r, _ := regexp.Compile("\\x1b\\[[;\\d]*[A-Za-z]")
	return r.ReplaceAllString(text, "")
}
