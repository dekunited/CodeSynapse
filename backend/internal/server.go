package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"
)

/*
  Run() will start up the http server and listen for any errors
*/
func Run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	server := NewServer()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	errs := make(chan error, 1)
	go func() {
		fmt.Fprintf(w, "Starting server on :8080\n")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()

	// Block until interrupt or error
	select {
	case err := <-errs:
		fmt.Fprintf(w, "Error starting server: %v\n", err)
		return err
	case <-ctx.Done():
		fmt.Fprintf(w, "Shutdown signal received, gracefully shutting down...\n")

		// Create a timeout for shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer shutdownCancel()

		// Shutdown gracefully
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(w, "Error during shutdown: %v\n", err)
			return err
		}

		fmt.Fprintf(w, "Server gracefully stopped\n")
		return nil
	}
}

/*
  NewServer() sets up the routes, middleware, etc.
*/
func NewServer(
// add dependencies here
) http.Handler {
	mux := http.NewServeMux()

	// add routes here
	AddRoutes(mux)

	var handler http.Handler = mux

	// Setup CORS middleware (for local development)
	// Will need to be configured properly in production
	handler = cors.Default().Handler(handler)

	return handler
}
