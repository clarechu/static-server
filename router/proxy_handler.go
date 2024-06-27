package router

import (
	"errors"
	"fmt"
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	if trip.StatusCode != 200 {
		data, err := io.ReadAll(trip.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
			return
		}
		w.WriteHeader(trip.StatusCode)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, errors.New(string(data)))))
		return
	}
	data, err := io.ReadAll(trip.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, err.Error())))
		return
	}
	for k, vs := range trip.Header {
		for _, v := range vs {
			w.Header().Set(k, v)
		}
	}
	w.Write(data)
}
