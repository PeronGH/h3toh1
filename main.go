package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./h3toh1 <http3 host>")
		os.Exit(1)
	}

	parsedUrl, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing host: %v", err)
	}

	host := parsedUrl.Host
	hostUrl := url.URL{Scheme: "https", Host: host}

	proxy := httputil.NewSingleHostReverseProxy(&hostUrl)
	proxy.Transport = &http3.RoundTripper{}
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = host
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	fmt.Printf("visit %s at http://localhost:%d\n", host, port)
	http.Serve(listener, proxy)
}
