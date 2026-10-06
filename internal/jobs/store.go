package jobs

import "errors"

type Store struct {
	jobs map[int]*Job
}

// Create a new Store for storing the job
func NewStore() *Store {
	return &Store{jobs: make(map[int]*Job)}
}

// Add new job
func (s *Store) Add(j Job) {
	jobCopy := j
	s.jobs[j.ID] = &jobCopy
}

// finding the job using id
// returns a pointer to the job if found
func (s *Store) Find(id int) (*Job, error) {
	j, exists := s.jobs[id]
	if !exists {
		return nil, errors.New("Job not found!")
	}
	return j, nil
}

// Return all jobs
func (s *Store) All() []Job {
	list := make([]Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		list = append(list, *j)
	}
	return list
}

// NextID returns the next available ID for a new job
func (s *Store) NextID() int {
	return len(s.jobs) + 1
}
