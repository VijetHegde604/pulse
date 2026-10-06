package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

func HandleAdd(args []string, store *jobs.Store) {
	//  Create a dedicated FlagSet for subcommands
	cmdFlags := flag.NewFlagSet("add", flag.ExitOnError)

	var name string
	var command string

	//  Bind long flags
	cmdFlags.StringVar(&name, "name", "", "Name of the job")
	cmdFlags.StringVar(&command, "command", "", "Command to run")

	//  Bind short flags to the SAME variables
	cmdFlags.StringVar(&name, "n", "", "Name of the job (shorthand)")
	cmdFlags.StringVar(&command, "c", "", "Command to run (shorthand)")

	// Parse ONLY the subcommand arguments (os.Args[2:])
	cmdFlags.Parse(args)

	if name == "" || command == "" {
		fmt.Println("Error: both -n/--name and -c/--command are required.")
		cmdFlags.Usage()
		os.Exit(1)
	}

	//  Capture remaining positional arguments for the job command
	cmdArgs := cmdFlags.Args()

	//  Pass job directly (NewJob returns Job value)
	job := jobs.NewJob(store.NextID(), name, command, cmdArgs)
	store.Add(*job)

	fmt.Printf("Job #%d ('%s') added successfully.\n", job.ID, job.Name)
}
