package cli

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/VijetHegde604/pulse/internal/db"
	"github.com/VijetHegde604/pulse/internal/jobs"
	_ "modernc.org/sqlite"
)

func setupTestStore(t *testing.T) (*jobs.Store, *sql.DB) {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.Migrate(database); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	return jobs.NewStore(database), database
}

func TestAddAndRunExecution(t *testing.T) {
	store, database := setupTestStore(t)
	defer database.Close()

	// 1. Add unquoted arguments
	HandleAdd([]string{"-n", "job-unquoted", "-c", "echo", "hello", "world"}, store)

	allJobs, err := store.All()
	if err != nil {
		t.Fatalf("store.All() error: %v", err)
	}
	if len(allJobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(allJobs))
	}
	j1 := allJobs[0]
	if j1.Command != "echo" {
		t.Errorf("expected command 'echo', got %q", j1.Command)
	}
	if len(j1.Args) != 2 || j1.Args[0] != "hello" || j1.Args[1] != "world" {
		t.Errorf("expected args ['hello', 'world'], got %#v", j1.Args)
	}
	if j1.Status != jobs.StatusPending {
		t.Errorf("expected initial status pending, got %s", j1.Status)
	}

	// Run job 1
	err = jobs.Run(&j1, store)
	if err != nil {
		t.Fatalf("jobs.Run failed: %v", err)
	}

	refreshedJ1, err := store.Find(j1.ID)
	if err != nil {
		t.Fatalf("store.Find failed: %v", err)
	}
	if refreshedJ1.Status != jobs.StatusCompleted {
		t.Errorf("expected status 'completed' in DB, got %s", refreshedJ1.Status)
	}

	// 2. Add quoted command string
	HandleAdd([]string{"-n", "job-quoted", "-c", "echo quoted string"}, store)
	refreshedJobs, err := store.All()
	if err != nil {
		t.Fatalf("store.All() error: %v", err)
	}
	if len(refreshedJobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(refreshedJobs))
	}
	j2 := refreshedJobs[1]
	if j2.Command != "echo" {
		t.Errorf("expected command 'echo', got %q", j2.Command)
	}
	if len(j2.Args) != 2 || j2.Args[0] != "quoted" || j2.Args[1] != "string" {
		t.Errorf("expected args ['quoted', 'string'], got %#v", j2.Args)
	}

	// Run job 2
	err = jobs.Run(&j2, store)
	if err != nil {
		t.Fatalf("jobs.Run failed: %v", err)
	}

	refreshedJ2, err := store.Find(j2.ID)
	if err != nil {
		t.Fatalf("store.Find failed: %v", err)
	}
	if refreshedJ2.Status != jobs.StatusCompleted {
		t.Errorf("expected status 'completed' in DB, got %s", refreshedJ2.Status)
	}

	// 3. Add with POSIX -- delimiter
	HandleAdd([]string{"-n", "job-dashdash", "--", "echo", "posix"}, store)
	allJobs, _ = store.All()
	j3 := allJobs[2]
	if j3.Command != "echo" || len(j3.Args) != 1 || j3.Args[0] != "posix" {
		t.Errorf("expected command 'echo' and arg 'posix', got %q %#v", j3.Command, j3.Args)
	}

	// 4. Add with --command= syntax
	HandleAdd([]string{"-n", "job-equals", "--command=echo equals"}, store)
	allJobs, _ = store.All()
	j4 := allJobs[3]
	if j4.Command != "echo" || len(j4.Args) != 1 || j4.Args[0] != "equals" {
		t.Errorf("expected command 'echo' and arg 'equals', got %q %#v", j4.Command, j4.Args)
	}

	// 5. Test failing job updates status to failed in DB
	HandleAdd([]string{"-n", "job-fail", "-c", "false"}, store)
	allJobs, _ = store.All()
	j5 := allJobs[4]
	err = jobs.Run(&j5, store)
	if err == nil {
		t.Fatalf("expected command 'false' to fail, but it succeeded")
	}

	refreshedJ5, err := store.Find(j5.ID)
	if err != nil {
		t.Fatalf("store.Find failed: %v", err)
	}
	if refreshedJ5.Status != jobs.StatusFailed {
		t.Errorf("expected status 'failed' in DB, got %s", refreshedJ5.Status)
	}
}

func TestAutoRepairLegacyJobs(t *testing.T) {
	store, database := setupTestStore(t)
	defer database.Close()

	// Simulate legacy job that had entire command stored in Command and empty Args
	legacyJob := jobs.NewJob("legacy", "echo legacy repaired", []string{})
	if err := store.Add(legacyJob); err != nil {
		t.Fatalf("failed to add legacy job: %v", err)
	}

	// Run should auto-repair and complete successfully
	err := jobs.Run(legacyJob, store)
	if err != nil {
		t.Fatalf("legacy job failed to run: %v", err)
	}

	refreshed, err := store.Find(legacyJob.ID)
	if err != nil {
		t.Fatalf("failed to find job: %v", err)
	}
	if refreshed.Command != "echo" {
		t.Errorf("expected repaired command 'echo', got %q", refreshed.Command)
	}
	if len(refreshed.Args) != 2 || refreshed.Args[0] != "legacy" || refreshed.Args[1] != "repaired" {
		t.Errorf("expected repaired args ['legacy', 'repaired'], got %#v", refreshed.Args)
	}
	if refreshed.Status != jobs.StatusCompleted {
		t.Errorf("expected status completed, got %s", refreshed.Status)
	}
}
