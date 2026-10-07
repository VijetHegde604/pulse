package jobs

import "time"

type Job struct {
	ID        int64
	Name      string
	Command   string
	Args      []string
	Status    JobStatus
	CreatedAt time.Time
}

// Custom string type to handle status
type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
	StatusCancelled JobStatus = "cancelled"
)

func NewJob(name, command string, args []string) *Job {
	return &Job{
		Name:      name,
		Command:   command,
		Args:      args,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}
