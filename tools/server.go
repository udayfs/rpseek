package tools

import (
	"fmt"
	"net"
	"net/http"

	"github.com/udayfs/rpseek/internal"
)

type Server struct {
	Host string
	Port int
}

func (s *Server) Serve() error {
	var err error

	if err = internal.ClearConsole(); err != nil {
		return err
	}

	addr := s.Host + fmt.Sprint(":", s.Port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("Server started listening for requests on http://%s/\n", addr)

	server := &http.Server{Addr: addr}
	return server.Serve(ln)
}
