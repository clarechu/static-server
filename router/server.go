package router

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer(root *Root) *Server {
	//文件浏览
	r := mux.NewRouter()
	if root.Path == "/" {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(root.FileDir)))
	} else {
		r.PathPrefix(root.Path).HandlerFunc(IndexHandler(root.Index))
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(root.FileDir)))
	}
	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    fmt.Sprintf(":%d", root.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return &Server{
		sv: srv,
	}
}

func (s *Server) Run() {
	log.Printf("Starting up http-server, serving ./dist")
	log.Printf("Available on:")
	log.Printf("   http://127.0.0.1%s", s.sv.Addr)
	log.Printf("Hit CTRL-C to stop the server")
	log.Fatal(s.sv.ListenAndServe())
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}
