package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"github.com/gorilla/mux"
)

func proxyTo(targetURL string, prefixToStrip string) http.HandlerFunc {
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefixToStrip)
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	r := mux.NewRouter()

	// Routing to microservices
	r.PathPrefix("/api/auth").HandlerFunc(proxyTo("http://localhost:3001", "/api"))
	r.PathPrefix("/api/tasks").HandlerFunc(proxyTo("http://localhost:3002", "/api"))
	r.PathPrefix("/api/ai").HandlerFunc(proxyTo("http://localhost:3003", "/api"))

	// Static files
	distPath := "./dist"
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(distPath, r.URL.Path)
		_, err := os.Stat(path)
		if os.IsNotExist(err) || r.URL.Path == "/" || !strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			return
		}
		http.FileServer(http.Dir(distPath)).ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" { port = "3000" }
	fmt.Printf("Gateway running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
