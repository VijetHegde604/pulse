# Pulse

> A lightweight, modular command-line job runner and task execution manager written in Go.

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](LICENSE)
[![Nix Flakes](https://img.shields.io/badge/Nix-Flakes%20Ready-5277C3?style=flat-square&logo=nixos)](flake.nix)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](https://github.com/VijetHegde604/pulse/pulls)

---

## Overview

**Pulse** is a minimalist CLI job runner built for developers who need a clean, structured way to define, inspect, and execute system commands directly from the terminal.

Built entirely with the **Go standard library** without heavy external dependencies, Pulse provides a focused workflow for task management, streaming command output in real time while tracking execution lifecycle states.

---

## Features

- **Zero External Runtime Dependencies**: Pure Go standard library (`os/exec`, `text/tabwriter`, `flag`).
- **Expressive CLI Interface**: Intuitive subcommands (`add`, `list`, `run`, `find`) with both short and long flag support.
- **Real-Time Output Streaming**: Directly streams task `stdout` and `stderr` to the terminal during execution.
- **Lifecycle Status Tracking**: Automatically monitors job status (`pending` -> `running` -> `completed` / `failed`).
- **Clean Tabular Formatting**: Aligned terminal tables showing Job IDs, names, commands, status, and creation timestamps.
- **Nix Flake & Direnv Support**: Ready-to-use isolated reproducible development environment out of the box.

---

## Architecture & Project Structure

Pulse is organized using standard Go modular architecture:

```
pulse/
├── flake.nix              # Nix development environment configuration
├── flake.lock             # Nix flake lockfile
├── go.mod                 # Go module definition
├── main.go                # Application entrypoint & subcommand router
└── internal/
    ├── cli/               # CLI handlers and flag parsing
    │   ├── add.go         # 'add' subcommand handler
    │   ├── find.go        # 'find' subcommand handler
    │   ├── list.go        # 'list' subcommand handler
    │   └── run.go         # 'run' subcommand handler
    └── jobs/              # Core domain models and business logic
        ├── jobs.go        # Job and JobStatus definitions
        ├── runner.go      # Process execution & status transition logic
        ├── store.go       # In-memory job repository
        └── printer.go     # Tabular formatting with text/tabwriter
```

### Core Components

| Package | Responsibility |
| :--- | :--- |
| `jobs.Job` | Represents a task unit containing an `ID`, `Name`, `Command`, `Args`, `Status`, and `CreatedAt` timestamp. |
| `jobs.Store` | Manages registration, indexing, sequential ID generation, and lookup of jobs. |
| `jobs.Run` | Spawns system processes via `os/exec`, streams stdout/stderr, and updates execution state. |
| `jobs.Write` | Renders job collections into human-readable aligned tables using `text/tabwriter`. |
| `cli.*` | Parses subcommand flags (`flag.FlagSet`) and delegates commands to the domain layer. |

---

## Getting Started

### Prerequisites

- **Go 1.26** or higher installed (or [Nix](https://nixos.org/) with Flakes enabled)
- A POSIX-compatible shell (Linux / macOS)

### Installation

#### 1. Clone the repository

```bash
git clone https://github.com/VijetHegde604/pulse.git
cd pulse
```

#### 2. Build the binary

```bash
# Build the binary
go build -o pulse .

# (Optional) Install globally to your $GOPATH/bin
go install .
```

#### 3. Using Nix & Direnv (Optional)

If you use Nix and direnv, enter the fully pre-configured dev shell with Go, `gopls`, `delve`, and `golangci-lint`:

```bash
direnv allow
# or manually:
nix develop
```

---

## Usage & Commands

Pulse operates via subcommands:

```bash
pulse <command> [options]
```

### 1. `pulse add` - Create a Job

Defines a new job with a descriptive name, target executable command, and any optional command arguments.

```bash
pulse add -n <name> -c <command> [arguments...]
```

**Flags:**
| Flag | Shorthand | Required | Description |
| :--- | :--- | :--- | :--- |
| `--name` | `-n` | Yes | Name / label for the job |
| `--command` | `-c` | Yes | Executable command to execute |

**Examples:**

```bash
# Add a ping test
pulse add -n "ping-dns" -c ping -c 4 8.8.8.8

# Add a Go build task
pulse add -n "build-app" -c go build ./...

# Add a directory listing
pulse add --name "list-files" --command ls -la
```

---

### 2. `pulse list` - Display All Jobs

Prints an aligned table of all registered jobs with their current status and creation timestamp.

```bash
pulse list
```

**Example Output:**

```text
ID    NAME          COMMAND    STATUS       CREATED AT
1     ping-dns      ping       pending      07 Oct 2026 01:30:15
2     build-app     go         completed    07 Oct 2026 01:31:02
3     list-files    ls         pending      07 Oct 2026 01:32:40
```

---

### 3. `pulse run` - Execute a Job

Runs the job matching the specified ID. Streams output in real time and updates the job's status to `completed` or `failed`.

```bash
pulse run -i <id>
```

**Flags:**
| Flag | Shorthand | Required | Description |
| :--- | :--- | :--- | :--- |
| `--id` | `-i` | Yes | Target Job ID |

**Example:**

```bash
pulse run -i 1
```

**Output:**

```text
Running Job #1 ('ping-dns')...

PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
64 bytes from 8.8.8.8: icmp_seq=1 ttl=116 time=14.2 ms
64 bytes from 8.8.8.8: icmp_seq=2 ttl=116 time=13.8 ms

--- 8.8.8.8 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms

Job #1 completed successfully.
```

---

### 4. `pulse find` - Inspect a Single Job

Looks up a specific job by its ID and displays its details.

```bash
pulse find -id <id>
```

**Flags:**
| Flag | Required | Description |
| :--- | :--- | :--- |
| `-id` | Yes | Target Job ID to search |

**Example:**

```bash
pulse find -id 1
```

**Output:**

```text
ID    NAME        COMMAND    STATUS       CREATED AT
1     ping-dns    ping       completed    07 Oct 2026 01:30:15
```

---

## Job Lifecycle

Each job transitions through distinct lifecycle states managed by the runner:

```
[ pending ] ---> ( pulse run ) ---> [ running ] ---> [ completed ] (exit 0)
                                            |
                                            +------> [ failed ]    (exit non-zero)
```

- **`pending`**: Job has been registered and is queued to run.
- **`running`**: Job process is actively executing with stdout/stderr attached.
- **`completed`**: Process finished successfully (exit code 0).
- **`failed`**: Process failed or exited with an error.
- **`cancelled`**: Job was cancelled before completion.

---

## Development & Testing

### Running locally

```bash
go run . <command>
```

### Formatting & Linting

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (if installed)
golangci-lint run
```

---

## Roadmap

- [ ] **Persistent Storage**: Save jobs to SQLite or JSON file across CLI sessions.
- [ ] **Background Execution**: Daemonized asynchronous job execution mode.
- [ ] **Execution History & Logs**: Log stdout and stderr to disk for historical replay.
- [ ] **Job Cancellation & Timeouts**: Support graceful cancellation and execution timeouts.
- [ ] **Cron / Scheduled Triggers**: Run recurring jobs on a schedule.

---

## Author

- **Vijet Hegde** - [@VijetHegde604](https://github.com/VijetHegde604)

---

## License

This project is licensed under the [MIT License](LICENSE) (or applicable project license).
