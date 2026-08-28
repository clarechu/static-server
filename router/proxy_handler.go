package router

import (
	"encoding/json"
	"io"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ProxyRouter struct {
	defaultTransport *http.Transport
	URI              *url.URL
}

func NewProxyRouter(host string) *ProxyRouter {
	defaultTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	uri, err := url.Parse(host)
	if err != nil {
		panic(err)
	}
	if uri.Scheme == "" {
		uri.Scheme = "http"
	}
	return &ProxyRouter{
		URI:              uri,
		defaultTransport: defaultTransport,
	}
}

func (ro *ProxyRouter) Proxy(w http.ResponseWriter, r *http.Request) {
	klog.Infof("begin proxy request: %s, proxy host: %s", r.RequestURI, ro.URI.Host)

	r.URL.Host = ro.URI.Host
	r.URL.Scheme = ro.URI.Scheme
	trip, err := ro.defaultTransport.RoundTrip(r)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if trip.StatusCode != 200 {
		data, err := io.ReadAll(trip.Body)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSONError(w, trip.StatusCode, string(data))
		return
	}
	data, err := io.ReadAll(trip.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	for k, vs := range trip.Header {
		for _, v := range vs {
			w.Header().Set(k, v)
		}
	}
	w.Write(data)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	resp := map[string]string{"message": message}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		klog.Errorf("failed to encode error response: %v", err)
	}
}
