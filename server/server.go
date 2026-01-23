package server

import (
	"bufio"
	"bytes"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	isRunning atomic.Bool
	Listen    net.Listener
	handler   Handler
}

func Serve(port int, handler Handler) (*Server, error) {

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return &Server{}, err
	}

	s := new(Server)
	s.isRunning.Store(true)
	s.Listen = l
	s.handler = handler

	go func() {
		s.listen()
	}()

	return s, nil
}

func (s *Server) Close() error {
	s.isRunning.Store(false)

	err := s.Listen.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) listen() {
	for {
		if !s.isRunning.Load() {
			return
		}

		conn, err := s.Listen.Accept()
		if err != nil {

			if s.isRunning.Load() {
				log.Println(err)
				continue
			} else {
				return
			}

		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {

	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}
	var body bytes.Buffer

	herr := s.handler(&body, req)

	if herr == nil {
		writer := bufio.NewWriter(conn)

		newHeaders := response.GetDefaultHeaders(len(body.Bytes()))

		err = response.WriteStatusLine(writer, response.Code200)
		if err != nil {
			return
		}

		err = response.WriteHeaders(writer, newHeaders)
		if err != nil {
			return
		}

		_, err = writer.Write(body.Bytes())
		if err != nil {
			return
		}
		err = writer.Flush()
		if err != nil {
			return
		}

	} else {
		writer := bufio.NewWriter(conn)
		err2 := WriteHandlerError(writer, herr)
		if err2 != nil {
			writer.Write([]byte(err2.Error()))
		}
		writer.Flush()
	}

}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(writer io.Writer, herr *HandlerError) error {
	if herr == nil {
		return nil
	}

	if err := response.WriteStatusLine(writer, herr.StatusCode); err != nil {
		return err
	}

	body := herr.Message
	headers := response.GetDefaultHeaders(len(body))

	if err := response.WriteHeaders(writer, headers); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(body)); err != nil {
		return err
	}

	return nil
}
