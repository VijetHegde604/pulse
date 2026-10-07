package jobs

import (
	"fmt"
	"os"
	"os/exec"
)

// Run the command with the arguements and update the status accordingly
func Run(j *Job) error {
	j.Status = StatusRunning
	fmt.Printf("Command: %q\n", j.Command)
	fmt.Printf("Args: %#v\n", j.Args)

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
