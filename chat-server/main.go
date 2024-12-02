package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"layer7/chat-server/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize handlers
	wsHandler := handlers.NewWSHandler()
	restHandler := handlers.NewRESTHandler()

	// Start the WebSocket room manager
	go wsHandler.Run()

	// Create gin router
	router := gin.Default()

	// WebSocket endpoint
	router.GET("/ws", wsHandler.HandleConnections)

	// REST endpoints
	router.GET("/api/messages", func(c *gin.Context) {
		restHandler.HandleMessages(c.Writer, c.Request)
	})
	router.POST("/api/messages", func(c *gin.Context) {
		restHandler.HandleMessages(c.Writer, c.Request)
	})
	router.GET("/api/users", func(c *gin.Context) {
		restHandler.HandleUsers(c.Writer, c.Request)
	})
	router.POST("/api/users", func(c *gin.Context) {
		restHandler.HandleUsers(c.Writer, c.Request)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		log.Printf("Server starting on port %s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Printf("server error: %v", err)

	case sig := <-shutdown:
		log.Printf("shutdown started, signal: %v", sig)
		
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("forcing server to close: %v", err)
			}
		}
	}
}