package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

func HandleAdd(args []string, store *jobs.Store) {
	// Find -c / --command first.
	commandIndex := -1

	for i, arg := range args {
		if arg == "-c" || arg == "--command" {
			commandIndex = i
			break
		}
	}

	if commandIndex == -1 || commandIndex+1 >= len(args) {
		fmt.Println("Error: -c/--command is required.")
		os.Exit(1)
	}

	// Everything before -c belongs to Pulse.
	pulseArgs := args[:commandIndex]

	// Everything after -c belongs to the job command.
	commandArgs := args[commandIndex+1:]

	// Create a dedicated FlagSet for the add command.
	cmdFlags := flag.NewFlagSet("add", flag.ExitOnError)

	var name string

	cmdFlags.StringVar(&name, "name", "", "Name of the job")
	cmdFlags.StringVar(&name, "n", "", "Name of the job (shorthand)")

	// Parse only Pulse's arguments.
	cmdFlags.Parse(pulseArgs)

	if name == "" {
		fmt.Println("Error: -n/--name is required.")
		cmdFlags.Usage()
		os.Exit(1)
	}

	if len(commandArgs) == 0 {
		fmt.Println("Error: command cannot be empty.")
		os.Exit(1)
	}

	// First argument is the executable.
	command := commandArgs[0]

	// Everything after the executable is its arguments.
	jobArgs := commandArgs[1:]

	job := jobs.NewJob(name, command, jobArgs)

	if err := store.Add(job); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Job #%d ('%s') added successfully.\n", job.ID, job.Name)
}
