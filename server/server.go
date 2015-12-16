package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	SockAddr   = "/tmp/metadesk.sock"
	ArgDelimit = "\n"
	MaxArgs    = 255
)

const (
	SuccessCode = iota
	ErrorCode
)

type Server struct {
	listener    net.Listener
	conn        net.Conn
	subscribers []*Subscriber
}

func NewServer() (server *Server, err error) {
	os.Remove(SockAddr)

	server = new(Server)
	server.subscribers = make([]*Subscriber, 0)
	server.listener, err = net.Listen("unix", SockAddr)

	return
}

func (this *Server) NextConn() error {
	if conn, err := this.listener.Accept(); err != nil {
		return errors.New("Failed to read connection from socket")
	} else {
		this.conn = conn
	}
	return nil
}

func (this *Server) ReadCommand() (args []string, err error) {
	reader := bufio.NewReader(this.conn)

	// read number of args
	argcStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, errors.New("Failed to read argument count")
	}
	argc, err := strconv.Atoi(argcStr[:len(argcStr)-1])
	if err != nil {
		return nil, errors.New("Non-integer argument count")
	}

	args = make([]string, argc)
	for i := range args {
		arg, err := reader.ReadString('\n')
		if err != nil {
			return nil, errors.New("Failed to read all arguments")
		}

		args[i] = arg[:len(arg)-1]
	}

	return args, nil
}

func (this *Server) WriteResponse(msg string, code int) error {
	writer := bufio.NewWriter(this.conn)
	var err error
	if len(msg) > 0 {
		_, err = writer.WriteString(fmt.Sprintf("%d\n%s\n", code, msg))
	} else {
		_, err = writer.WriteString(fmt.Sprintf("%d\n", code))
	}

	if err != nil {
		return errors.New("Failed to write response")
	}

	if writer.Flush() != nil {
		return errors.New("Failed to flush response")
	}

	if this.conn.Close() != nil {
		return errors.New("Failed to close connection")
	}

	return nil
}

func (this *Server) SubscribeConn(desktop *Desktop, formatter Formatter) {
	sub := &Subscriber{
		conn:      this.conn,
		writer:    bufio.NewWriter(this.conn),
		desktop:   desktop,
		formatter: formatter,
	}

	this.subscribers = append(this.subscribers, sub)
}

func (this *Server) UpdateSubscribers() error {
	for i, sub := range this.subscribers {
		if !sub.codeSupplied {
			_, err := sub.writer.WriteString(fmt.Sprintf("%d\n", SuccessCode))
			if err != nil {
				return errors.New("Failed to write success code")
			}

			sub.codeSupplied = true
		}

		var unsubbed bool
		outStr := sub.formatter.Format(sub.desktop) + "\n"
		if _, err := sub.writer.WriteString(outStr); err != nil {
			unsubbed = true
		}
		if sub.writer.Flush() != nil {
			unsubbed = true
		}

		if unsubbed {
			this.subscribers = append(this.subscribers[:i], this.subscribers[i+1:]...)
		}
	}

	return nil
}

func (this *Server) Close() {
	for _, s := range this.subscribers {
		s.conn.Close()
	}

	this.listener.Close()
}

type Subscriber struct {
	conn         net.Conn
	writer       *bufio.Writer
	desktop      *Desktop
	formatter    Formatter
	codeSupplied bool
}
