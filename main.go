package main

import (
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/cli"
	"github.com/VijetHegde604/pulse/internal/jobs"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Initialize the store
	store := jobs.NewStore()

	subcommand := os.Args[1]
	args := os.Args[2:]

	switch subcommand {
	case "add":
		cli.HandleAdd(args, store)
	case "list":
		cli.HandleList(args, store)
	case "run":
		cli.HandleRun(args, store)
	// case "find":
	// 	handleFind(args, store)
	default:
		fmt.Printf("Unknown command: %s\n\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: pulse <command> [options]")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  add   - Create a new job")
	fmt.Println("  list  - Display all jobs")
	fmt.Println("  run   - Execute a job by ID")
	fmt.Println("  find  - Find and inspect a job by ID")
}
