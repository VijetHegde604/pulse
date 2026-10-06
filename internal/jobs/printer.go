package jobs

import (
	"fmt"
	"io"
	"text/tabwriter"
)

func Write(w io.Writer, jobs []Job) {

	// Create a new tabwriter instance
	tw := tabwriter.NewWriter(w, 0, 0, 4, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "ID\tNAME\tCOMMAND\tSTATUS\tCREATED AT")

	for _, j := range jobs {
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%s\t%s\n",
			j.ID,
			j.Name,
			j.Command,
			j.Status,
			j.CreatedAt.Format("02 Jan 2006 15:04:05"),
		)
	}
}
