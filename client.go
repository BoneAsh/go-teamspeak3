package teamspeak3

import (
	"errors"
	"strings"
	"time"
)

const TitleCheck = "TS3"

type Client struct {
	protocol     Protocol
	query        Query
	responsePipe <-chan string
	messagePipe  chan Message
	errorPipe    chan Error
	notifyPipe   chan Notify
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
		_ = client.Close()
		return nil, err
	}
	responsePipe, err := client.query.GetResponsePipe()
	client.messagePipe = make(chan Message, DefaultMsgPipeLength*2)
	client.errorPipe = make(chan Error, DefaultMsgPipeLength*2)
	client.notifyPipe = make(chan Notify, DefaultMsgPipeLength*4)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	client.responsePipe = responsePipe
	if title := <-responsePipe; title != TitleCheck {
		_ = client.Close()
		return nil, errors.New("title check failed")
	}
	<-responsePipe
	go client.Loop()
	return client, nil
}

func (c *Client) Close() (err error) {
	if c.protocol != nil {
		err = c.protocol.Disconnect()
		if err != nil {
			return err
		}
	}
	if c.query != nil {
		err = c.query.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Loop() {
	var builder strings.Builder
	for {
		r := <-c.responsePipe
		if strings.HasPrefix(r, "error") {
			c.errorPipe <- NewError(r)
		} else if strings.HasPrefix(r, "notify") {
			c.notifyPipe <- NewNotify(r)
		} else {
			builder.Reset()
			builder.WriteString(r)
			for {
				r = <-c.responsePipe
				if strings.HasPrefix(r, "error") {
					c.messagePipe <- NewMessage(builder.String())
					c.errorPipe <- NewError(r)
					break
				} else {
					builder.WriteString(r)
				}
			}
		}
	}
}
