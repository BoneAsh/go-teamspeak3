package teamspeak3

import (
	"regexp"
	"strconv"
	"strings"
)

func ReplaceAnsiEscapeCode(text string) string {
	r, _ := regexp.Compile("\\x1b\\[[;\\d]*[A-Za-z]")
	return r.ReplaceAllString(text, "")
}

func DecodeResponse(text string) (res map[string]interface{}) {
	res = make(map[string]interface{})
	textSplits := strings.Split(text, " ")
	for _, textSplit := range textSplits {
		if strings.ContainsRune(textSplit, '=') {
			sliceSplit := strings.SplitN(textSplit, "=", 2)
			if strings.ContainsRune(sliceSplit[1], '|') {
				sliceSplitSplit := strings.Split(sliceSplit[1], "|")
				var item []map[string]interface{}
				for _, s := range sliceSplitSplit {
					item = append(item, DecodeResponse(s))
				}
				res[sliceSplit[0]] = item
			} else {
				if i, err := strconv.Atoi(sliceSplit[1]); err != nil {
					res[sliceSplit[0]] = sliceSplit[1]
				} else {
					res[sliceSplit[0]] = i
				}

			}
		} else {
			res[textSplit] = nil
		}
	}
	return res
}
