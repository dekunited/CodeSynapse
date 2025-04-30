package main

import (
	"CodeSynapse/internal"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()

	if err := internal.Run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error (ctrlc): %v\n", err)
		os.Exit(1)
	}
}
