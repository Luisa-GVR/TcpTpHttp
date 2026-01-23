package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

var badRequestHTML = []byte("<html>\n  <head>\n    <title>400 Bad Request</title>\n  </head>\n  <body>\n    <h1>Bad Request</h1>\n    <p>Your request honestly kinda sucked.</p>\n  </body>\n</html>")

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
	rw := response.NewWriter(conn)

	if err != nil {
		headers := response.GetDefaultHeaders(len(badRequestHTML))
		headers.Set("content-type", "text/html")

		rw.WriteStatusLine(response.Code400)
		rw.WriteHeaders(headers)
		rw.WriteBody([]byte(badRequestHTML))
		return
	}

	s.handler(rw, req)

}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

// WriteHandlerError Not being used, should be removed later
func WriteHandlerError(writer io.Writer, herr *HandlerError) error {
	if herr == nil {
		return nil
	}

	rw := response.NewWriter(writer)

	if err := rw.WriteStatusLine(herr.StatusCode); err != nil {
		return err
	}

	body := herr.Message
	headers := response.GetDefaultHeaders(len(body))

	if err := rw.WriteHeaders(headers); err != nil {
		return err
	}
	if _, err := writer.Write([]byte(body)); err != nil {
		return err
	}

	return nil
}
