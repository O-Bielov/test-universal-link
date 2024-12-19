package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	server := startMainRouter()
	go startBackend(server)

	sig := <-shutdown // Directly receive from the channel
	log.Printf("main: %v : Start shutdown", sig)
}

func startBackend(server *Server) {
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
