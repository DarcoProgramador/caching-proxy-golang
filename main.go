package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/Darcoprogramador/caching-proxy-go/database"
	"github.com/Darcoprogramador/caching-proxy-go/utils"
	"github.com/redis/go-redis/v9"
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

	db := database.Init()
	defer db.Close()

	proxy := &proxyCache{
		transport: http.DefaultTransport,
		targetURL: *originServer,
		cache:     db,
	}

	http.ListenAndServe(":"+strconv.Itoa(*portNumber), proxy)
}

type proxyCache struct {
	transport http.RoundTripper
	targetURL string
	cache     *redis.Client
}

func (p *proxyCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path

	// Create a new HTTP request with the same method, URL, and body as the original request
	proxyReq, err := http.NewRequest(r.Method, p.targetURL+path, nil)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// Get Cache with redis
	//var bodyCache cacheResponse
	strResponseCache, err := p.cache.Get(ctx, path).Result()

	if err != nil || strResponseCache == "" {
		slog.Info("No cache", slog.String("path", path))
	} else {
		w.Header().Add("X-Cache", "HIT")
		rBody := bufio.NewReader(bytes.NewReader([]byte(strResponseCache)))
		responseBody, err := http.ReadResponse(rBody, nil)

		if err != nil {
			slog.Error("Error reading cache", "error", err.Error())
		}

		CopyHeaders(responseBody.Header, w.Header())

		io.Copy(w, responseBody.Body)

		return
	}

	// Copy the headers from the original request to the proxy request
	CopyHeaders(r.Header, proxyReq.Header)

	// Send the proxy request using the custom transport
	// resp, err := p.transport.RoundTrip(proxyReq)
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}
	// Save cache with redis
	if err := p.cache.Set(ctx, path, string(body), time.Minute*1).Err(); err != nil {
		slog.Error("Error setting cache", "error", err.Error())
	}

	// Copy the headers from the proxy response to the original response
	CopyHeaders(resp.Header, w.Header())

	// Set the status code of the original response to the status code of the proxy response
	w.Header().Add("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}

func CopyHeaders(from, to http.Header) {
	for name, values := range from {
		for _, value := range values {
			to.Add(name, value)
		}
	}
}
