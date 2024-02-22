package proxy

import (
	"net/http"
	"net/http/httputil"
)

// ReverseTarget reverse proxies to target server
func ReverseTarget(reverseProxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reverseProxy.ServeHTTP(w, r)
	}
}
