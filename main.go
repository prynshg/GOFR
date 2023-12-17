package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/prynshg/gooooo/datastore"
	"github.com/prynshg/gooooo/handler"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	s, err := datastore.New()
	if err != nil {
		// Handle the error, such as logging or exiting the program.
		fmt.Printf("Failed to initialize datastore: %v\n", err)
		os.Exit(1)
	}

	h := handler.New(s)

	// Handle panics
	defer func() {
		if r := recover(); r != nil {
			// Log or print the panic information
			fmt.Printf("Panic recovered: %v\n", r)
		}
		// Close the datastore connection on exit
		s.Close()
	}()

	// Set up graceful shutdown
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-gracefulShutdown
		// Perform cleanup before exiting
		fmt.Println("Received signal for graceful shutdown")
		s.Close()
		os.Exit(0)
	}()

	// Register routes
	app.GET("/students/{id}", h.GetByID)
	app.POST("/students", h.Create)
	app.PUT("/students/{id}", h.Update)
	app.DELETE("/students/{id}", h.Delete)

	// Start the server on a custom port
	app.Server.HTTP.Port = 9092
	app.Start()
}
