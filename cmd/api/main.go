package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/didil/nginx-lb-updater/server"
	"github.com/didil/nginx-lb-updater/server/handlers"
	"github.com/didil/nginx-lb-updater/services"
)

func main() {
	lbUpdater, err := services.NewLBUpdater()
	if err != nil {
		log.Fatalf("failed to initialize lb updater: %v", err)
	}

	logger, err := services.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	app := handlers.NewApp(lbUpdater, logger)

	r := server.NewRouter(app)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on %s\n", addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("ListenAndServer err: %v", err)
	}
}
