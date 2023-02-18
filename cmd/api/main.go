package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/didil/nginx-lb-updater/server"
)

func main() {
	r := server.NewRouter()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on %s\n", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("ListenAndServer err: %v", err)
	}
}
