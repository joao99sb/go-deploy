package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func startNewReverseProxy(target *url.URL) *httputil.ReverseProxy {
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	return reverseProxy
}

func ProxyRequestHandler(p *Proxy, endpoint string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		if p.CurrentUrl == nil {
			http.Error(w, "Proxy não configurado", http.StatusInternalServerError)
			return
		}

		reverseProxy := startNewReverseProxy(p.CurrentUrl)

		r.URL.Host = p.CurrentUrl.Host
		r.URL.Scheme = p.CurrentUrl.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = p.CurrentUrl.Host
		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, endpoint)

		ctx, cancel := context.WithTimeout(r.Context(), GetTimeOut())
		defer cancel()

		// Cria uma nova requisição com o contexto que inclui o timeout
		newReq := r.WithContext(ctx)
		// Cria um ResponseRecorder para capturar a resposta
		rr := httptest.NewRecorder()

		reverseProxy.ServeHTTP(rr, newReq)

		if rr.Code == http.StatusBadGateway {
			Error503Message(w)
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

func Error503Message(w http.ResponseWriter) {

	content, err := os.ReadFile("html/503.html")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo HTML:", err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintf(w, "%s", content)
}
