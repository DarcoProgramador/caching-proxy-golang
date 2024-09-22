package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Darcoprogramador/caching-proxy-go/utils"
)

var portNumber = flag.Int("port", 8080, "Port to listen on")
var originServer = flag.String("origin", "", "Origin server to proxy requests to")

func main() {
	flag.Parse()
	fmt.Println("Starting proxy server...", *portNumber, *originServer)

	if *originServer == "" {
		panic("Origin server is required")
	}

	if !utils.IsUrl(*originServer) {
		panic("Invalid origin server")
	}

	proxy := &proxy{
		transport: http.DefaultTransport,
		targetURL: *originServer,
	}
	//db := database.Init()

	http.ListenAndServe(":"+strconv.Itoa(*portNumber), proxy)
}

type proxy struct {
	transport http.RoundTripper
	targetURL string
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP request with the same method, URL, and body as the original request
	path := r.URL.Path
	proxyReq, err := http.NewRequest(r.Method, p.targetURL+path, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Copy the headers from the original request to the proxy request
	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Send the proxy request using the custom transport
	resp, err := p.transport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}
