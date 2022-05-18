package teamspeak3

import (
	"github.com/mitchellh/mapstructure"
	"strings"
)

type Error struct {
	Id  int
	Msg string
}

func (e *Error) Decode(content map[string]interface{}) (err error) {
	return mapstructure.Decode(content, &e)
}

func NewError(content string) (e Error, err error) {
	contentSplits := strings.SplitN(content, " ", 2)
	e = Error{}
	err = e.Decode(DecodeResponse(contentSplits[1]))
	if err != nil {
		return e, err
	}
	return e, nil
}
