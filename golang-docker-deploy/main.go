package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-docker-deploy/docker"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HelloServer)
	// Create Server and Route Handlers
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8081",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Server!")
	id, err := docker.CreateNewContainer("artofimagination/worker-server", "0.0.0.0", "8082")
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	fmt.Fprintln(w, id)
	log.Println("Hello, Server...")
}
