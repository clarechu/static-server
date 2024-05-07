// Copyright (c) 2021 The static-server Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	log "k8s.io/klog/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type StaticServerConfig struct {
	Routers []*StaticRouter `yaml:"routers,omitempty"`
}

type StaticRouter struct {
	FileDir    string `yaml:"file_dir,omitempty" json:"file_dir,omitempty"`
	Path       string `yaml:"path,omitempty" json:"path,omitempty"`
	PublicPath string `json:"public_path,omitempty" yaml:"public_path,omitempty"`
	IsGzip     bool   `json:"is_gzip,omitempty" yaml:"is_gzip,omitempty"`
}

func NewServer(root *Root) *Server {
	//文件浏览
	r := mux.NewRouter()
	if root.Config == "" {
		staticAssetsHandler, err := NewStaticAssetsHandler("", StaticAssetsHandlerOptions{
			FileDir:  root.FileDir,
			BasePath: root.PublicPath,
			IsGzip:   root.IsGzip,
		})
		if err != nil {
			log.Warningf("new static assets handler :%v", err)
		} else {
			staticAssetsHandler.RegisterRoutes(r)
		}
	} else {
		config := getConfig(root.Config)
		for _, router := range config.Routers {
			staticAssetsHandler, err := NewStaticAssetsHandler("", StaticAssetsHandlerOptions{
				FileDir:  router.FileDir,
				BasePath: router.PublicPath,
				IsGzip:   router.IsGzip,
			})
			if err != nil {
				log.Warningf("new static assets handler :%v", err)
			} else {
				staticAssetsHandler.StaticRegisterRoutes(r)
			}
		}
	}

	addHTTPMiddleware(r)
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

func getConfig(path string) *StaticServerConfig {
	file, err := os.ReadFile(path)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	server := &StaticServerConfig{}
	err = yaml.Unmarshal(file, server)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return server
}

func addHTTPMiddleware(router *mux.Router) {
	router.Use(CORSMethodMiddleware(router))
	//router.Use(LogMiddleware(router))
}

func (s *Server) Run() {
	log.V(0).Info("Starting up http-server, serving ./dist")
	log.V(0).Info("Available on:")
	log.V(0).Infof("   http://127.0.0.1%s", s.sv.Addr)
	log.V(0).Infof("Hit CTRL-C to stop the server")
	log.Fatal(s.sv.ListenAndServe())
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
	path := r.URL.Path
	path = strings.Replace(path, h.rootPath, "", 1)
	r.URL.Path = path
	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		path = path + ".html"
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			http.ServeFile(w, r, path)
			return
		} else if err == nil {
			http.ServeFile(w, r, path)
			return
		}
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
