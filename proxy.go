package main

import (
	"fmt"
	"net/url"
)

type Proxy struct {
	CurrentUrl        *url.URL
	proxyChangeSignal chan string
}

func NewProxy() *Proxy {
	return &Proxy{
		proxyChangeSignal: make(chan string),
	}
}

func (p *Proxy) setProxyURL(newURL string) {
	u, err := url.Parse(newURL)
	if err != nil {
		fmt.Println("Erro ao analisar a URL:", err)
		return
	}

	p.CurrentUrl = u
}

func (p *Proxy) ChangeProxyURL() {
	for {
		newURL := <-p.proxyChangeSignal
		p.setProxyURL(newURL)

		fmt.Println("URL do proxy atualizada para:", newURL)

	}
}
