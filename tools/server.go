package tools

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/udayfs/rpseek/internal"
)

type Server struct {
	Host string
	Port int
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./web/index.html")
}

func jsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
	http.ServeFile(w, req, "./web/js/rpseek.js")
}

func searchHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(os.Stdout, req.Method, req.URL)
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

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/js/rpseek.js", jsHandler)

	server := &http.Server{Addr: addr}
	return server.Serve(ln)
}
