package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// logs holds all log messages. We'll print them after stopping the spinner.
var logs []string

func main() {
	// 1) Clear the screen and move the cursor to top-left:
	fmt.Print("\033[2J\033[H")

	// 2) Print a banner that says "Cleanify".
	banner := color.CyanString(`
   ____ _                    __ _         
  / ___| | ___  __ _ _ __  / _(_) __ _  
 | |   | |/ _ \/ _` + "`" + ` | '_ \| |_| |/ _` + "`" + ` | 
 | |___| |  __/ (_| | | | |  _| | (_| | 
  \____|_|\___|\__,_|_| |_|_| |_|\__, | 
                                  |___/ 
`)
	fmt.Println(banner)

	color.Yellow("Starting cleanup process...\n")

	// 3) Move the cursor to a “fixed” line for the spinner.
	spinnerLine := 8
	fmt.Printf("\033[%d;0H", spinnerLine)

	// 4) Start a spinner in that fixed position.
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Suffix = "  Cleaning in progress..."
	s.Start()

	// 5) Perform the cleanup (storing logs in memory).
	cleanup()

	// 6) Stop the spinner and move the cursor below it to print logs.
	s.Stop()
	logsStartLine := spinnerLine + 2
	fmt.Printf("\033[%d;0H", logsStartLine)

	// 7) Print all accumulated logs with colors intact.
	for _, line := range logs {
		fmt.Println(line)
	}

	color.Green("\nCleanup complete! 🚀\n")
}

// cleanup does the actual work and appends results to the global `logs` slice.
func cleanup() {
	// Dynamically retrieve the current user's TEMP directory.
	userTemp := os.Getenv("TEMP")
	if userTemp == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logs = append(logs, color.RedString("Error retrieving user home directory: %v", err))
			return
		}
		userTemp = filepath.Join(homeDir, "AppData", "Local", "Temp")
	}

	// List of directories to clean.
	tempDirs := []string{
		userTemp,              // Current user's TEMP directory
		`C:\Windows\Temp`,     // Windows system temporary directory
		`C:\Windows\Prefetch`, // Windows Prefetch directory
	}

	// Clean each directory in turn.
	for _, dir := range tempDirs {
		logs = append(logs, color.MagentaString("\nCleaning directory: %s", dir))

		successes, failures := cleanDir(dir)

		// Append success/failure messages to logs (with color).
		for _, msg := range successes {
			logs = append(logs, color.GreenString(msg))
		}
		for _, msg := range failures {
			logs = append(logs, color.RedString(msg))
		}

		if len(failures) > 0 {
			logs = append(logs, color.RedString("Some files in %s could not be removed.", dir))
		} else {
			logs = append(logs, color.GreenString("Finished cleaning directory: %s", dir))
		}
	}
}

// cleanDir attempts to remove all files/subfolders in dirPath, returning
// two lists of strings: one for successes and one for failures.
func cleanDir(dirPath string) (successes []string, failures []string) {
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		failures = append(failures, fmt.Sprintf("Reading directory %s failed: %v", dirPath, err))
		return
	}

	for _, item := range items {
		fullPath := filepath.Join(dirPath, item.Name())
		if err := os.RemoveAll(fullPath); err != nil {
			failures = append(failures, fmt.Sprintf("Failed to remove %s: %v", fullPath, err))
		} else {
			successes = append(successes, fmt.Sprintf("Successfully removed %s", fullPath))
		}
	}
	return
}
