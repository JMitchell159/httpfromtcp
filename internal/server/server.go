package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/JMitchell159/httpfromtcp/internal/response"
)

type Server struct {
	Listener net.Listener
	State    atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		Listener: l,
	}
	server.State.Store(true)

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	s.State.Store(false)
	return nil
}

func (s *Server) listen() {
	for s.State.Load() {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Println(err)
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response.WriteStatusLine(conn, 200)
	h := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, h)
	conn.Close()
}
