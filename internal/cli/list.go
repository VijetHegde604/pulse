package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

// HandleList handles the 'list' subcommand flag parsing and invokes table printing.
func HandleList(args []string, store *jobs.Store) {
	cmdFlags := flag.NewFlagSet("list", flag.ExitOnError)
	cmdFlags.Parse(args)

	// Retrieve all jobs from store
	allJobs, err := store.All()
	if err != nil {
		fmt.Println(err)
	}

	// Handle empty state gracefully
	if len(allJobs) == 0 {
		fmt.Println("No jobs found. Add one using 'pulse add -n <name> -c <cmd>'")
		return
	}

	jobs.Write(os.Stdout, allJobs)

}
