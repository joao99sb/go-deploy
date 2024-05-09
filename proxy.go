package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}

func ProxyRequestHandler(proxy *httputil.ReverseProxy, url *url.URL, endpoint string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url.Host
		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint)

		ctx, cancel := context.WithTimeout(r.Context(), GetTimeOut())
		defer cancel()

		// Cria uma nova requisição com o contexto que inclui o timeout
		newReq := r.WithContext(ctx)
		// Cria um ResponseRecorder para capturar a resposta
		rr := httptest.NewRecorder()

		proxy.ServeHTTP(rr, newReq)

		if rr.Code == http.StatusBadGateway {
			http.Error(w, "O servidor de destino está inacessível", http.StatusServiceUnavailable)
			return
		}
		for k, vv := range rr.Header() {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(rr.Code)

		io.Copy(w, rr.Body)
	}
}
