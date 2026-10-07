package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

// HandleFind parses CLI flags and looks up a job by its ID.
func HandleFind(args []string, store *jobs.Store) {
	cmdFlags := flag.NewFlagSet("find", flag.ExitOnError)
	var id int64
	cmdFlags.Int64Var(&id, "id", 0, "Job ID to find")
	cmdFlags.Int64Var(&id, "i", 0, "Job ID to find (shorthand)")
	cmdFlags.Parse(args)

	if id == 0 {
		fmt.Println("Error: job ID is required (-i or --id)")
		cmdFlags.Usage()
		return
	}

	foundJob, err := store.Find(id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if foundJob == nil {
		fmt.Println("No jobs found with that ID")
		return
	}

	jobs.Write(os.Stdout, []jobs.Job{*foundJob})
}
