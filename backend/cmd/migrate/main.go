package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Desire-in-tech/morrow-and-muse/backend/db"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is unset")
		os.Exit(1)
	}

	ctx := context.Background()
	var err error
	switch os.Args[1] {
	case "up":
		err = db.Up(ctx, databaseURL)
	case "down-to":
		if len(os.Args) != 3 {
			usage()
			os.Exit(2)
		}
		version, parseErr := strconv.ParseInt(os.Args[2], 10, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "version: %v\n", parseErr)
			os.Exit(2)
		}
		err = db.DownTo(ctx, databaseURL, version)
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate up")
	fmt.Fprintln(os.Stderr, "       go run ./cmd/migrate down-to 0")
}
