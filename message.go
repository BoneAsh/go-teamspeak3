package teamspeak3

type Message struct {
}

func NewMessage(content string) (m Message, err error) {
	m = Message{}
	return m, nil
}
