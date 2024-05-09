package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	TIMEOUT_DEFAULT = "300"
)

func config() []resource {
	config, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	resourcesList := config.Resources
	PORT := config.Server.Listen_port
	TIMEOUT_STRING := config.Server.Timeout

	var TIMEOUT string
	if TIMEOUT_STRING == "" {
		TIMEOUT = TIMEOUT_DEFAULT
	} else {
		TIMEOUT = TIMEOUT_STRING
	}

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
