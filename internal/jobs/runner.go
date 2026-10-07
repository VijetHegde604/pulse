package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Run executes the job's command with its arguments, streams stdout and stderr,
// and persists status transitions (running -> completed / failed) to the store.
func Run(j *Job, store *Store) error {
	// Auto-repair: if Command contains spaces and Args is empty,
	// parse it into command and args (handles legacy or single-string entries).
	if len(j.Args) == 0 && strings.Contains(j.Command, " ") {
		tokens, err := SplitCommandLine(j.Command)
		if err == nil && len(tokens) > 0 {
			j.Command = tokens[0]
			j.Args = tokens[1:]
			if store != nil {
				_ = store.Update(j)
			}
		}
	}

	j.Status = StatusRunning
	if store != nil {
		if err := store.UpdateStatus(j.ID, StatusRunning); err != nil {
			return fmt.Errorf("failed to update job status to running: %w", err)
		}
	}

	fmt.Printf("Command: %q\n", j.Command)
	fmt.Printf("Args: %#v\n\n", j.Args)

	cmd := exec.Command(j.Command, j.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		j.Status = StatusFailed
		if store != nil {
			_ = store.UpdateStatus(j.ID, StatusFailed)
		}
		return err
	}

	j.Status = StatusCompleted
	if store != nil {
		if err := store.UpdateStatus(j.ID, StatusCompleted); err != nil {
			return fmt.Errorf("failed to update job status to completed: %w", err)
		}
	}

	return nil
}
