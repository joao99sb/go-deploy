package main

import (
	"fmt"
	"net/http"
	"os"
)

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
