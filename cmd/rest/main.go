package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"timble/internal/config"
)

func main() {
	restServer, err := config.NewRestServer()
	if err != nil {
		log.Fatalf("failed to create REST server %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func(r *config.RESTServer) {
		log.Printf("Server running on port %d", r.ServerConfig.ServerPort)
		if err := r.Server.ListenAndServe(); err != nil {
			log.Fatalf("Server error %s", err)
		}
	}(restServer)

	go func(r *config.RESTServer) {
		log.Printf("Prometheus running on port %d", r.PrometheusServerConfig.PrometheusPort)
		if err := r.PrometheusServer.ListenAndServe(); err != nil {
			log.Fatalf("Prometheus server error %s", err)
		}
	}(restServer)

	<-sigChan
	log.Println("Received terminate, graceful shutdown")

}
