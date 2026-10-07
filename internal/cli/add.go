package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/VijetHegde604/pulse/internal/jobs"
)

// HandleAdd parses CLI arguments, isolates Pulse flags from the child command,
// and saves the new job to the store.
func HandleAdd(args []string, store *jobs.Store) {
	commandIndex := -1
	var pulseArgs []string
	var commandArgs []string

	// Locate the command delimiter (-c, --command, or --)
	for i, arg := range args {
		if arg == "-c" || arg == "--command" || arg == "--" {
			commandIndex = i
			pulseArgs = args[:i]
			commandArgs = args[i+1:]
			break
		}
		if strings.HasPrefix(arg, "-c=") {
			commandIndex = i
			pulseArgs = args[:i]
			val := strings.TrimPrefix(arg, "-c=")
			commandArgs = append([]string{val}, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "--command=") {
			commandIndex = i
			pulseArgs = args[:i]
			val := strings.TrimPrefix(arg, "--command=")
			commandArgs = append([]string{val}, args[i+1:]...)
			break
		}
	}

	if commandIndex == -1 {
		fmt.Println("Error: -c/--command (or --) is required.")
		fmt.Println("\nUsage: pulse add -n <name> -c <command> [arguments...]")
		os.Exit(1)
	}

	if len(commandArgs) == 0 {
		fmt.Println("Error: command cannot be empty.")
		os.Exit(1)
	}

	// Create a dedicated FlagSet for Pulse's add options
	cmdFlags := flag.NewFlagSet("add", flag.ContinueOnError)

	var name string
	cmdFlags.StringVar(&name, "name", "", "Name of the job")
	cmdFlags.StringVar(&name, "n", "", "Name of the job (shorthand)")

	if err := cmdFlags.Parse(pulseArgs); err != nil {
		os.Exit(1)
	}

	if name == "" {
		// Provide a helpful hint if the user placed -n/--name after the command flag
		for _, carg := range commandArgs {
			if carg == "-n" || carg == "--name" || strings.HasPrefix(carg, "-n=") || strings.HasPrefix(carg, "--name=") {
				fmt.Println("Error: -n/--name must be specified before -c/--command.")
				fmt.Println("Everything after -c/--command is treated as the command to execute.")
				os.Exit(1)
			}
		}

		fmt.Println("Error: -n/--name is required.")
		cmdFlags.Usage()
		os.Exit(1)
	}

	var command string
	var jobArgs []string

	if len(commandArgs) == 1 {
		tokens, err := jobs.SplitCommandLine(commandArgs[0])
		if err != nil {
			fmt.Printf("Error parsing command: %v\n", err)
			os.Exit(1)
		}
		if len(tokens) == 0 {
			fmt.Println("Error: command cannot be empty.")
			os.Exit(1)
		}
		command = tokens[0]
		jobArgs = tokens[1:]
	} else {
		if strings.Contains(commandArgs[0], " ") {
			firstTokens, err := jobs.SplitCommandLine(commandArgs[0])
			if err != nil {
				fmt.Printf("Error parsing command: %v\n", err)
				os.Exit(1)
			}
			if len(firstTokens) == 0 {
				fmt.Println("Error: command cannot be empty.")
				os.Exit(1)
			}
			command = firstTokens[0]
			jobArgs = append(firstTokens[1:], commandArgs[1:]...)
		} else {
			command = commandArgs[0]
			jobArgs = commandArgs[1:]
		}
	}

	job := jobs.NewJob(name, command, jobArgs)

	if err := store.Add(job); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Job #%d ('%s') added successfully.\n", job.ID, job.Name)
}
