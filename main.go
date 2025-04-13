package main

import (
	"log"
	"net/http"
)

const (
	port    = "8080"
	rootDir = "."
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(rootDir)))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving '%v' on port: %v\n", rootDir, port)
	log.Fatal(server.ListenAndServe())
}
