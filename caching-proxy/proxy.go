package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// ProxyHandler handles incoming requests and serves cached or remote responses
type ProxyHandler struct {
	cache     *InMemoryCache
	targetURL *url.URL
}

// NewProxyHandler creates a new ProxyHandler
func NewProxyHandler(target string, cacheTTL time.Duration) (*ProxyHandler, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	cache := NewInMemoryCache(cacheTTL)
	return &ProxyHandler{cache: cache, targetURL: targetURL}, nil
}

// ServeHTTP handles incoming HTTP requests
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.String()
	if entry, found := p.cache.Get(key); found {
		p.serveCachedResponse(w, entry)
		return
	}
	p.proxyAndCacheRequest(w, r)
}

// serveCachedResponse serves a cached response
func (p *ProxyHandler) serveCachedResponse(w http.ResponseWriter, entry CacheEntry) {
	w.Header().Set("Content-Type", entry.ContentType)
	w.Header().Set("X-Cache", "Hit")
	w.WriteHeader(entry.StatusCode)
	_, err := w.Write(entry.Response)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

// proxyAndCacheRequest proxies the request and caches the response
func (p *ProxyHandler) proxyAndCacheRequest(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	// r.RequestURI = ""
	resp, err := client.Get(p.targetURL.String() + r.URL.Path)
	// resp, err := client.Do(r)

	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Println("Error proxying request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Println("Error reading response body:", err)
		return
	}

	entry := CacheEntry{
		Response:    body,
		ContentType: resp.Header.Get("Content-Type"),
		StatusCode:  resp.StatusCode,
		Expiry:      time.Now().Add(p.cache.ttl),
	}
	p.cache.Set(r.URL.String(), entry)
	p.serveCachedResponse(w, entry)

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("X-Cache", "Miss")
	fmt.Println("set here")
	w.WriteHeader(resp.StatusCode)

	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response:", err)
	}
}
