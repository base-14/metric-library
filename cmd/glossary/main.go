package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/base14/otel-glossary/internal/adapter/kubernetes/cadvisor"
	"github.com/base14/otel-glossary/internal/adapter/kubernetes/ksm"
	"github.com/base14/otel-glossary/internal/adapter/otelcontrib"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/kafka"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/mongodb"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/mysql"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/node"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/postgres"
	"github.com/base14/otel-glossary/internal/adapter/prometheus/redis"
	"github.com/base14/otel-glossary/internal/api"
	"github.com/base14/otel-glossary/internal/orchestrator"
	"github.com/base14/otel-glossary/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return runServe()
	}

	switch os.Args[1] {
	case "serve":
		return runServe()
	case "extract":
		return runExtract(os.Args[2:])
	default:
		return runServe()
	}
}

func runServe() error {
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

func runExtract(args []string) error {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	adapterName := fs.String("adapter", "otel-collector-contrib", "Adapter to use for extraction")
	cacheDir := fs.String("cache-dir", "", "Directory to cache git repositories")
	force := fs.Bool("force", false, "Force re-fetch even if cached")
	dbPath := fs.String("db", "", "Database path (default: $DATABASE_PATH or ./data/glossary.db)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *dbPath == "" {
		*dbPath = os.Getenv("DATABASE_PATH")
		if *dbPath == "" {
			*dbPath = "./data/glossary.db"
		}
	}

	if *cacheDir == "" {
		*cacheDir = os.Getenv("CACHE_DIR")
		if *cacheDir == "" {
			*cacheDir = "./.cache"
		}
	}

	if err := os.MkdirAll(filepath.Dir(*dbPath), 0750); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	log.Printf("Connecting to database at %s", *dbPath)
	s, err := store.NewSQLiteStoreWithMigrations(*dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer func() { _ = s.Close() }()

	var adp orchestrator.Adapter
	switch *adapterName {
	case "otel-collector-contrib":
		adp = otelcontrib.NewAdapter(*cacheDir)
	case "prometheus-postgres":
		adp = postgres.NewAdapter(*cacheDir)
	case "prometheus-node":
		adp = node.NewAdapter(*cacheDir)
	case "prometheus-redis":
		adp = redis.NewAdapter(*cacheDir)
	case "prometheus-mysql":
		adp = mysql.NewAdapter(*cacheDir)
	case "prometheus-mongodb":
		adp = mongodb.NewAdapter(*cacheDir)
	case "prometheus-kafka":
		adp = kafka.NewAdapter(*cacheDir)
	case "kubernetes-ksm":
		adp = ksm.NewAdapter(*cacheDir)
	case "kubernetes-cadvisor":
		adp = cadvisor.NewAdapter(*cacheDir)
	default:
		return fmt.Errorf("unknown adapter: %s", *adapterName)
	}

	log.Printf("Starting extraction with adapter: %s", adp.Name())
	log.Printf("Cache directory: %s", *cacheDir)

	ext := orchestrator.NewExtractor(adp, s)
	ctx := context.Background()

	result, err := ext.Run(ctx, orchestrator.Options{
		CacheDir: *cacheDir,
		Force:    *force,
	})
	if err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	log.Printf("Extraction completed successfully")
	log.Printf("  Adapter: %s", result.AdapterName)
	log.Printf("  Commit: %s", result.Commit)
	log.Printf("  Metrics extracted: %d", result.MetricsExtracted)
	log.Printf("  Metrics stored: %d", result.MetricsStored)
	log.Printf("  Duration: %s", result.Duration)

	return nil
}
