package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TIMEOUT_DEFAULT = "300"
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

		if rr.Code != http.StatusOK {
			// Trata o erro aqui, por exemplo, enviando uma resposta de erro
			http.Error(w, "O servidor de destino está inacessível", http.StatusServiceUnavailable)
			return
		}
		for k, vv := range rr.Header() {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(rr.Code)
		w.Write(rr.Body.Bytes())
	}
}

func config() []resource {
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)

	}
	resourcesList := config.Resources
	HOST := config.Server.Host
	PORT := config.Server.Listen_port
	TIMEOUT_STRING := config.Server.Timeout

	var TIMEOUT string
	if TIMEOUT_STRING == "" {

		TIMEOUT = TIMEOUT_DEFAULT
	} else {
		TIMEOUT = TIMEOUT_STRING
	}

	os.Setenv("HOST", HOST)
	os.Setenv("PORT", PORT)
	os.Setenv("TIMEOUT", TIMEOUT)
	return resourcesList
}

func GetTimeOut() time.Duration {
	timeout := os.Getenv("TIMEOUT")

	timeout_number, err := strconv.Atoi(timeout)
	if err != nil {

		log.Fatal(err)
	}
	duratoin := time.Duration(timeout_number) * time.Millisecond

	return duratoin
}

func main() {
	resourcesList := config()

	mux := http.NewServeMux()

	for _, resource := range resourcesList {
		url, _ := url.Parse(resource.Destination_URL)
		proxy := NewProxy(url)
		mux.HandleFunc(resource.Endpoint, ProxyRequestHandler(proxy, url, resource.Endpoint))
	}

	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))

	fmt.Printf("server started on port: %s\n", PORT)

	http.ListenAndServe(PORT, mux)

}
