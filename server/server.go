package server

import (
	"bufio"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	isRunning atomic.Bool
	Listen    net.Listener
}

func Serve(port int) (*Server, error) {

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return &Server{}, err
	}

	s := new(Server)
	s.isRunning.Store(true)
	s.Listen = l

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

	writer := bufio.NewWriter(conn)

	response.WriteStatusLine(writer, response.Code200)

	h := response.GetDefaultHeaders(0)

	err := response.WriteHeaders(writer, h)
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}

	conn.Close()

}
