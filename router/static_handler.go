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
	"github.com/NYTimes/gziphandler"
	"io"
	log "k8s.io/klog/v2"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

// QueryOptions holds configuration for query service
type QueryOptions struct {
	// FileDir static file dir (/opt/solar-mesh/dist)
	FileDir string
	// BasePath is the prefix for all UI and API HTTP routes
	BasePath string
	// StaticAssets is the path for the static assets for the UI (https://github.com/uber/jaeger-ui)
	StaticAssets string
	// AdditionalHeaders
	AdditionalHeaders http.Header
	// MaxClockSkewAdjust is the maximum duration by which jaeger-query will adjust a span
	MaxClockSkewAdjust time.Duration
}

var (
	favoriteIcon    = "favicon.ico"
	staticRootFiles = []string{favoriteIcon}

	// The following patterns are searched and replaced in the index.html as a way of customizing the UI.
	basePathPattern = regexp.MustCompile(`<base href="/"`) // Note: tag is not closed

)

// RegisterStaticHandler adds handler for static assets to the router.
func RegisterStaticHandler(r *mux.Router, qOpts *QueryOptions) {
	staticHandler, err := NewStaticAssetsHandler(qOpts.StaticAssets, StaticAssetsHandlerOptions{
		BasePath: qOpts.BasePath,
	})

	if err != nil {
		log.Fatalf("Could not create static assets handler", err)
	}

	staticHandler.RegisterRoutes(r)
}

// StaticAssetsHandler handles static assets
type StaticAssetsHandler struct {
	options   StaticAssetsHandlerOptions
	indexHTML atomic.Value // stores []byte
	assetsFS  http.FileSystem
}

// StaticAssetsHandlerOptions defines options for NewStaticAssetsHandler
type StaticAssetsHandlerOptions struct {
	// FileDir static file dir (/opt/solar-mesh/dist)
	FileDir  string
	BasePath string
	IsGzip   bool
}

// NewStaticAssetsHandler returns a StaticAssetsHandler
func NewStaticAssetsHandler(staticAssetsRoot string, options StaticAssetsHandlerOptions) (*StaticAssetsHandler, error) {
	assetsFS := http.Dir(options.FileDir)
	if staticAssetsRoot != "" {
		assetsFS = http.Dir(staticAssetsRoot)
	}

	indexHTML, err := loadAndEnrichIndexHTML(assetsFS.Open, options)
	if err != nil {
		return nil, err
	}

	h := &StaticAssetsHandler{
		options:  options,
		assetsFS: assetsFS,
	}

	h.indexHTML.Store(indexHTML)

	return h, nil
}

func loadAndEnrichIndexHTML(open func(string) (http.File, error), options StaticAssetsHandlerOptions) ([]byte, error) {
	indexBytes, err := loadIndexHTML(open)
	if err != nil {
		return nil, fmt.Errorf("cannot load index.html: %w", err)
	}
	/*	// replace UI config
		if configObject, err := loadUIConfig(options.UIConfigPath); err != nil {
			return nil, err
		} else if configObject != nil {
			indexBytes = configObject.regexp.ReplaceAll(indexBytes, configObject.config)
		}*/
	// replace base path
	if options.BasePath == "" {
		options.BasePath = "/"
	}
	if options.BasePath != "/" {
		if !strings.HasPrefix(options.BasePath, "/") || strings.HasSuffix(options.BasePath, "/") {
			return nil, fmt.Errorf("invalid base path '%s'. Must start but not end with a slash '/', e.g. '/jaeger/ui'", options.BasePath)
		}
		indexBytes = basePathPattern.ReplaceAll(indexBytes, []byte(fmt.Sprintf(`<base href="%s/"`, options.BasePath)))
	}

	return indexBytes, nil
}

func loadIndexHTML(open func(string) (http.File, error)) ([]byte, error) {
	indexFile, err := open("/index.html")
	if err != nil {
		return nil, fmt.Errorf("cannot open index.html: %w", err)
	}
	defer indexFile.Close()
	indexBytes, err := io.ReadAll(indexFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read from index.html: %w", err)
	}
	return indexBytes, nil
}

// RegisterRoutes registers routes for this handler on the given router
func (sH *StaticAssetsHandler) RegisterRoutes(router *mux.Router) {
	router.NotFoundHandler = http.HandlerFunc(sH.notFoundHandler)
	router = router.PathPrefix(sH.options.BasePath).Subrouter()
	fileServer := http.FileServer(sH.assetsFS)
	if sH.options.BasePath != "/" {
		fileServer = http.StripPrefix(sH.options.BasePath+"/", fileServer)
	}
	if sH.options.IsGzip {
		fileServer = gziphandler.GzipHandler(fileServer)
	}
	// 这个地方会导致 notFoundHandler 失效
	router.PathPrefix("/static").Handler(fileServer)
	// 处理根路径下面的 favicon.ico 文件
	for _, file := range staticRootFiles {
		router.Path("/" + file).Handler(fileServer)
	}
}

func (sH *StaticAssetsHandler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(sH.indexHTML.Load().([]byte))
}
