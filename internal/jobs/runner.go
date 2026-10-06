package jobs

import (
	"os"
	"os/exec"
)

// Run the command with the arguements and update the status accordingly
func Run(j *Job) error {
	j.Status = StatusRunning

	// Pass the running commands output and error to stdout
	// (Helps to print the output to the terminal)
	cmd := exec.Command(j.Command, j.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		j.Status = StatusFailed
		return err
	}
	j.Status = StatusCompleted
	return nil
}
