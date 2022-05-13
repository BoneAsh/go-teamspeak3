package teamspeak3

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

const (
	DefaultConnectTimeout = time.Second * 10
	DefaultMsgPipeLength  = 10
	MaxBufferSize         = 10 << 20
	DefaultBufferSize     = 4096
)

// Protocol provide an interface to implement transport protocol
type Protocol interface {
	Connect(host string, port int, username string, password string) error
	Disconnect() error
	SetInput(content string) error
	GetOutputPipe() (<-chan string, error)
}

type Type int

const (
	Ssh Type = iota
)

var protocolMap = map[Type]Protocol{
	Ssh: &sshProtocol{},
}

func NewProtocol(t Type) (p Protocol, err error) {
	if protocol, ok := protocolMap[t]; ok {
		return protocol, nil
	} else {
		return nil, errors.New(fmt.Sprintf("protocol type(%d) is not support", t))
	}
}

type sshProtocol struct {
	Host      string
	ip        *net.IPAddr
	Port      int
	Username  string
	Password  string
	PublicKey string

	client  *ssh.Client
	session *ssh.Session

	stdinPipe  io.Writer
	stdoutPipe io.Reader

	msgOutPipe chan string

	keepAliveDuration time.Duration
	keepAliveData     string
}

func (s *sshProtocol) Connect(host string, port int, username string, password string) (err error) {
	s.Host = host
	s.Port = port
	s.Username = username
	s.Password = password
	config := &ssh.ClientConfig{
		User:            s.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         DefaultConnectTimeout,
	}
	if len(s.PublicKey) != 0 {
		if authMethod, err := s.publicKeyAuthMethod(); err != nil {
			return err
		} else {
			config.Auth = []ssh.AuthMethod{authMethod}
		}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.Password(s.Password)}
	}

	s.ip, err = net.ResolveIPAddr("ip", s.Host)
	if err != nil {
		return err
	}
	s.client, err = ssh.Dial("tcp", s.ip.String()+":"+strconv.Itoa(s.Port), config)
	if err != nil {
		return err
	}
	s.session, err = s.client.NewSession()
	if err != nil {
		_ = s.Disconnect()
		return err
	}
	err = s.session.RequestPty("linux", 1024, 1024, ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	})
	if err != nil {
		_ = s.Disconnect()
		return err
	}
	s.stdinPipe, err = s.session.StdinPipe()
	if err != nil {
		_ = s.Disconnect()
		return err
	}
	s.stdoutPipe, err = s.session.StdoutPipe()
	if err != nil {
		_ = s.Disconnect()
		return err
	}
	err = s.session.Shell()
	if err != nil {
		_ = s.Disconnect()
		return err
	}
	s.msgOutPipe = make(chan string, DefaultMsgPipeLength)
	go func() {
		err = s.session.Wait()
		if err != nil {
			_ = s.Disconnect()
			return
		}
	}()
	go s.outputWorker()
	return nil
}

func (s *sshProtocol) publicKeyAuthMethod() (authMethod ssh.AuthMethod, err error) {
	var key []byte
	var signer ssh.Signer
	key, err = ioutil.ReadFile(s.PublicKey)
	if err != nil {
		return nil, err
	}
	signer, err = ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

func (s *sshProtocol) Disconnect() (err error) {
	if s.session != nil {
		err = s.session.Close()
		if err != nil {
			return err
		}
	}
	if s.client != nil {
		err = s.client.Close()
		if err != nil {
			return err
		}
	}
	close(s.msgOutPipe)
	return nil
}

func (s *sshProtocol) outputWorker() {
	scanner := bufio.NewScanner(bufio.NewReader(s.stdoutPipe))
	scanner.Split(bufio.ScanLines)
	buffer := make([]byte, DefaultBufferSize)
	scanner.Buffer(buffer, MaxBufferSize)
	for scanner.Scan() {
		s.msgOutPipe <- scanner.Text()
	}
}

func (s *sshProtocol) SetInput(content string) (err error) {
	_, err = s.stdinPipe.Write([]byte(fmt.Sprintf("%v\n", content)))
	if err != nil {
		return
	}
	return nil
}

func (s *sshProtocol) GetOutputPipe() (outputChannel <-chan string, err error) {
	return s.msgOutPipe, err
}
