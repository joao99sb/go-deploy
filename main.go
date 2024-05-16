package main

import (
	"fmt"
	"log"
	"net/http"
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
	DEFAULT_PORT := config.Server.Default_port

	var TIMEOUT string
	if TIMEOUT_STRING == "" {
		TIMEOUT = TIMEOUT_DEFAULT
	} else {
		TIMEOUT = TIMEOUT_STRING
	}

	os.Setenv("PORT", PORT)
	os.Setenv("TIMEOUT", TIMEOUT)
	os.Setenv("DEFAULT_PORT", DEFAULT_PORT)

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

func HandleEvents(d *DockerHandler, p *Proxy) {

	for {
		message := <-d.MessageCh
		switch v := message.(type) {
		case NewProviderCommand:

			p.proxyChangeSignal <- v.Port
		}

	}

}

func main() {
	resourcesList := config()

	default_redirect_port := os.Getenv("DEFAULT_PORT")
	proxy := NewProxy()

	proxy.setProxyURL(default_redirect_port)

	dockerCli := NewDockerHandler(resourcesList)

	dockerCli.Start()
	go proxy.ChangeProxyURL()
	go HandleEvents(dockerCli, proxy)

	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))

	fmt.Printf("server started on port: %s\n", PORT)

	http.ListenAndServe(PORT, proxy)

}
