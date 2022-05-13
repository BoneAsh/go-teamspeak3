package teamspeak3

import "time"

type Client struct {
	protocol Protocol
	query    Query
}

func NewClient(t Type, host string, port int, username string, password string) (client *Client, err error) {
	client = &Client{}
	p, err := NewProtocol(t)
	if err != nil {
		return nil, err
	}
	client.protocol = p
	q, err := NewQuery(t)
	if err != nil {
		return nil, err
	}
	client.query = q
	err = client.protocol.Connect(host, port, username, password)
	if err != nil {
		return nil, err
	}
	err = client.query.Init(client.protocol, time.Second*200, "version", 2)
	if err != nil {
		return nil, err
	}
	return
}
