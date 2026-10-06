package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

// HandleRun parses CLI flags, retrieves the target job from the store, and executes it.
func HandleRun(args []string, store *jobs.Store) {
	cmdFlags := flag.NewFlagSet("run", flag.ExitOnError)

	var id int
	cmdFlags.IntVar(&id, "id", 0, "ID of the job to run")
	cmdFlags.IntVar(&id, "i", 0, "ID of the job to run (shorthand)")

	cmdFlags.Parse(args)

	if id == 0 {
		fmt.Println("Error: job ID is required (-i or --id)")
		cmdFlags.Usage()
		os.Exit(1)
	}

	// 1. Look up pointer to job in store
	job, err := store.Find(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running Job #%d ('%s')...\n\n", job.ID, job.Name)

	// 2. Pass *Job directly to jobs.Run
	// Status changes (StatusRunning -> StatusCompleted / StatusFailed)
	// apply directly to the job in the store since 'job' is a pointer.
	if err := jobs.Run(job); err != nil {
		fmt.Printf("\nJob #%d failed: %v\n", job.ID, err)
		os.Exit(1)
	}

	fmt.Printf("\nJob #%d completed successfully.\n", job.ID)
}
