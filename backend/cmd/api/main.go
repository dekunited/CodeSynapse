package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"CodeSynapse/internal"
)

func main() {
	ctx := context.Background()

	// Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Starting the server
	if err := internal.Run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error (ctrlc): %v\n", err)
		os.Exit(1)
	}
}
