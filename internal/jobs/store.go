package jobs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type Store struct {
	db *sql.DB
}

// Create a new Store for storing the job
func NewStore(database *sql.DB) *Store {
	return &Store{db: database}
}

// Add new job
func (s *Store) Add(j *Job) error {

	// Serialize the args into JSON
	// sqlite doesnt have []string support
	argsJSON, err := json.Marshal(j.Args)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		"INSERT INTO jobs (name, command, args, status, created_at) VALUES (?, ?, ?, ?, ?)",
		j.Name, j.Command, argsJSON, j.Status, j.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	j.ID = id

	return nil
}

// finding the job using id
// returns a pointer to the job if found
func (s *Store) Find(id int64) (*Job, error) {
	row := s.db.QueryRow(
		"SELECT id, name, command, args, status, created_at FROM jobs WHERE id = ?",
		id,
	)

	var job Job
	var argsJSON string

	err := row.Scan(
		&job.ID,
		&job.Name,
		&job.Command,
		&argsJSON,
		&job.Status,
		&job.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("job #%d not found", id)
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(argsJSON), &job.Args)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// UpdateStatus updates the status of an existing job in the database.
func (s *Store) UpdateStatus(id int64, status JobStatus) error {
	result, err := s.db.Exec("UPDATE jobs SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("job #%d not found", id)
	}
	return nil
}

// Update updates an existing job's name, command, args, and status in the database.
func (s *Store) Update(j *Job) error {
	argsJSON, err := json.Marshal(j.Args)
	if err != nil {
		return err
	}

	result, err := s.db.Exec(
		"UPDATE jobs SET name = ?, command = ?, args = ?, status = ? WHERE id = ?",
		j.Name, j.Command, argsJSON, j.Status, j.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("job #%d not found", j.ID)
	}
	return nil
}

// Return all jobs
func (s *Store) All() ([]Job, error) {
	list := []Job{}

	rows, err := s.db.Query(
		"SELECT id, name, command, args, status, created_at FROM jobs",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var job Job
		var argsJSON string

		err := rows.Scan(
			&job.ID,
			&job.Name,
			&job.Command,
			&argsJSON,
			&job.Status,
			&job.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(argsJSON), &job.Args)
		if err != nil {
			return nil, err
		}

		list = append(list, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
