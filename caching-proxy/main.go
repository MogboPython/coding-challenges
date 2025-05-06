package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", 0, "Port number to listen on")
	origin := flag.String("origin", "", "Allowed origin")
	clearCache := flag.Bool("clear-cache", false, "Clear all items from the cache.")

	flag.Parse()

	if *origin == "" || *port == 0 {
		fmt.Println("No action specified.\nUse --clear-cache to clear the cache\nProvide --port and --origin to start the server")
		flag.Usage()
		os.Exit(1)
	}

	cacheTTL := 1 * time.Minute
	proxy, err := NewProxyHandler(*origin, cacheTTL)
	if err != nil {
		log.Fatal("Error creating proxy handler:", err)
	}

	if *clearCache {
		proxy.cache.Clear()
	}

	log.Printf("Server listening on port %d", *port)
	log.Printf("Forwarding requests to %s", *origin)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), proxy); err != nil {
		log.Fatal("Server error:", err)
	}
}
