// main.go — Minimal Go starter with HTTP server and graceful shutdown
// Run:
//   go run main.go
// Build:
//   go build -o app ./

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on")
	shutdownTimeout := flag.Duration("shutdown-timeout", 5*time.Second, "graceful shutdown timeout")
	flag.Parse()

	logger := log.New(os.Stdout, "app: ", log.LstdFlags|log.Lmsgprefix)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, map[string]string{"message": "Hello from Go starter!"})
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Printf("starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown: %v", err)
	}
	logger.Println("server stopped gracefully")
}

func respondJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		// If encoding fails, fall back to plain error
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
