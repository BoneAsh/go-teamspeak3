package teamspeak3

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	MinKeepAliveDuration = time.Second * 5
)

// Query provide an interface to implement standard io with protocol io
type Query interface {
	Init(protocol Protocol, keepAliveDuration time.Duration, keepAliveData string, keepAliveResponseLines int) error
	Request(content string) error
	GetResponsePipe() (<-chan string, error)
	Close() error
}

var queryMap = map[Type]Query{
	Ssh: &sshQuery{},
}

func NewQuery(t Type) (q Query, err error) {
	if query, ok := queryMap[t]; ok {
		return query, nil
	} else {
		return nil, errors.New(fmt.Sprintf("query type(%d) is not support", t))
	}
}

type sshQuery struct {
	protocol               Protocol
	requestPipe            chan string
	responsePipe           chan string
	keepAliveDuration      time.Duration
	keepAliveData          string
	keepAliveResponseLines int
	keepAliveInformPipe    chan struct{}
	stopPipe               chan struct{}
	lastRequest            string
}

func (s *sshQuery) Init(protocol Protocol, keepAliveDuration time.Duration, keepAliveData string, keepAliveResponseLines int) (err error) {
	if protocol == nil {
		return errors.New("protocol is nil")
	}
	s.protocol = protocol
	s.requestPipe = make(chan string, DefaultMsgPipeLength)
	s.responsePipe = make(chan string, DefaultMsgPipeLength*2)
	s.stopPipe = make(chan struct{}, 2)
	if keepAliveDuration < MinKeepAliveDuration {
		keepAliveDuration = MinKeepAliveDuration
	}
	s.keepAliveDuration = keepAliveDuration
	s.keepAliveData = keepAliveData
	s.keepAliveInformPipe = make(chan struct{}, 1)
	s.keepAliveResponseLines = keepAliveResponseLines + 1
	go s.requestWorker()
	go s.responseWorker()
	return nil
}

func (s *sshQuery) Request(content string) (err error) {
	if s.protocol == nil {
		return errors.New("protocol is nil")
	}
	s.lastRequest = content
	s.requestPipe <- content
	return nil
}

func (s *sshQuery) GetResponsePipe() (channel <-chan string, err error) {
	if s.responsePipe == nil {
		return nil, errors.New("response pipe is nil")
	}
	return s.responsePipe, nil
}

func (s *sshQuery) Close() (err error) {
	s.stopPipe <- struct{}{}
	s.stopPipe <- struct{}{}
	close(s.requestPipe)
	close(s.responsePipe)
	close(s.stopPipe)
	return nil
}

func (s *sshQuery) requestWorker() {
	for {
		select {
		case <-s.stopPipe:
			return
		case <-time.After(s.keepAliveDuration):
			s.keepAliveInformPipe <- struct{}{}
			err := s.protocol.SetInput(s.keepAliveData)
			if err != nil {
				return
			}
		case msg := <-s.requestPipe:
			err := s.protocol.SetInput(msg)
			if err != nil {
				return
			}
		}
	}
}

func (s *sshQuery) responseWorker() {
	outputPipe, err := s.protocol.GetOutputPipe()
	if err != nil {
		return
	}
	for {
		select {
		case <-s.stopPipe:
			return
		case response := <-outputPipe:
			response = ReplaceAnsiEscapeCode(response)
			// skip input
			if strings.HasSuffix(response, fmt.Sprintf("> %s%s", s.lastRequest, s.lastRequest)) {
				continue
			}
			s.responsePipe <- response
		case <-s.keepAliveInformPipe:
			// todo: add keep alive response judgement
			for i := 0; i < s.keepAliveResponseLines; i++ {
				<-outputPipe
			}
		}
	}
}
