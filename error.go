package teamspeak3

type Error struct {
	Id  int
	Msg string
}

func NewError(content string) (e Error) {
	e = Error{}
	return e
}
