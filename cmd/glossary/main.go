package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/base14/otel-glossary/internal/api"
	"github.com/base14/otel-glossary/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/glossary.db"
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0750); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	log.Printf("Connecting to database at %s", dbPath)
	s, err := store.NewSQLiteStoreWithMigrations(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer func() { _ = s.Close() }()

	handler := api.NewHandler(s)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		close(done)
	}()

	log.Printf("Starting server on :%s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	<-done
	log.Println("Server stopped")
	return nil
}
