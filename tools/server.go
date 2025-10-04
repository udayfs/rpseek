package tools

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/udayfs/rpseek/internal"
)

const (
	MaxResponses = 10
)

type Server struct {
	Host      string
	Port      int
	IndexFile string
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./web/index.html")
}

func jsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
	http.ServeFile(w, req, "./web/js/rpseek.js")
}

func searchHandler(w http.ResponseWriter, req *http.Request, indexFilePath string) {
	var query internal.SearchQuery

	if err := json.NewDecoder(req.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokens := internal.Tokenize(query.Query)

	if res, err := internal.SearchDoc(indexFilePath, tokens); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.Header().Add("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(res[:MaxResponses])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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
	http.HandleFunc("/js/rpseek.js", jsHandler)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, s.IndexFile)
	})

	server := &http.Server{Addr: addr}
	return server.Serve(ln)
}
