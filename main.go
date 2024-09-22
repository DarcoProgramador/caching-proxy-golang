package main

import (
	"flag"
	"fmt"
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

	//db := database.Init()

	http.ListenAndServe(":"+strconv.Itoa(*portNumber), http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//check the path
			path := r.URL.Path
			fmt.Println("Path: ", path)
			w.Write([]byte("Hello, World!"))
		},
	))
}
