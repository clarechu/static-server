package router

import (
	"github.com/gorilla/mux"
	log "k8s.io/klog/v2"
	"net/http"
	"strings"
)

//LogMiddleware  automatically sets the Access-Control-Allow-Methods response header
func LogMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.RequestURI, "/api") {
				log.V(3).Infof("METHOD: %+v, URL: %+v, REMOTE IP:%+v", req.Method, req.RequestURI, req.RemoteAddr)

			} else {
				log.Infof("METHOD: %+v, URL: %+v, REMOTE IP:%+v", req.Method, req.RequestURI, req.RemoteAddr)
			}
			next.ServeHTTP(w, req)
		})
	}
}

var allMethods = []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions,
	http.MethodDelete, http.MethodConnect}
var allAllowHeaders = []string{
	"X-PINGOTHER", "Content-Type",
}

//CORSMethodMiddleware  automatically sets the Access-Control-Allow-Methods response header
func CORSMethodMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allMethods, ","))
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allAllowHeaders, ","))
			next.ServeHTTP(w, req)
		})
	}
}
