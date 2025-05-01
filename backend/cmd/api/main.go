package main

import (
	"context"
	"fmt"
	//"log"
	"os"

	//"github.com/joho/godotenv"

	"CodeSynapse/internal"
)

func main() {
	ctx := context.Background()
  /* Comment out once we have actual keys
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/

	if err := internal.Run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error (ctrlc): %v\n", err)
		os.Exit(1)
	}
}
