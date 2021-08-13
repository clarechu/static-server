package router

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewServer(root *Root) *Server {
	//文件浏览
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	/*	if root.Path == "/" {
			r.PathPrefix("/").Handler(http.FileServer(http.Dir(root.FileDir)))
		} else {
			r.PathPrefix(root.Path).HandlerFunc(IndexHandler(root.Index))
			r.PathPrefix("/").Handler(http.FileServer(http.Dir(root.FileDir)))
		}*/
	spa := spaHandler{staticPath: root.FileDir, indexPath: root.Index, rootPath: root.Path}
	r.PathPrefix(root.Path).Handler(spa)
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

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	s, err := ioutil.ReadFile("")
	if err != nil {
		return
	}
	w.Write(s)
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	rootPath   string
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = strings.Replace(path, h.rootPath, "", 1)
	r.URL.Path = path
	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		path = filepath.Join(h.staticPath, h.indexPath)
		http.ServeFile(w, r, path)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
