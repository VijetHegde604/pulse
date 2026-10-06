package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

func HandleFind(args []string, store *jobs.Store) {
	cmdFlags := flag.NewFlagSet("find", flag.ExitOnError)
	id := cmdFlags.Int("id", 0, "Job ID to find")
	cmdFlags.Parse(args)

	if *id == 0 {
		fmt.Println("Please provide a job ID using -id flag")
		return
	}

	foundJob, err := store.Find(*id)
	if err != nil {
		fmt.Println(err)
		return
	}
	if foundJob == nil {
		fmt.Println("No jobs found with that ID")
		return
	}

	jobs.Write(os.Stdout, []jobs.Job{*foundJob})
}
