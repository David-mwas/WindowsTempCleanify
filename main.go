package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/fatih/color"
)

// -------------------
// Model
// -------------------
type model struct {
	spinner spinner.Model
	logs    []string // Stores all log messages
	done    bool     // Indicates whether cleanup is finished
	err     error    // Stores an error if one occurs
}

// -------------------
// Message types
// -------------------
type logMsg []string // A batch of log lines
type errMsg error    // An error message

// -------------------
// Init: start spinner and cleanup concurrently
// -------------------
func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, doCleanupCmd())
}

// -------------------
// Update: process messages and key input
// -------------------
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		if !m.done {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd

	case logMsg:
		m.logs = msg
		m.done = true
		return m, nil

	case errMsg:
		m.err = msg
		m.done = true
		return m, nil

	case tea.KeyMsg:
		// If cleanup is complete, any key press quits the program.
		if m.done {
			return m, tea.Quit
		}
	}

	return m, nil
}

// -------------------
// View: display banner, spinner or logs, and exit prompt
// -------------------
func (m model) View() string {
	// Banner
	s := color.CyanString(`
   ____ _                    __ _         
  / ___| | ___  __ _ _ __  / _(_) __ _  
 | |   | |/ _ \/ _` + "`" + ` | '_ \| |_| |/ _` + "`" + ` | 
 | |___| |  __/ (_| | | | |  _| | (_| | 
  \____|_|\___|\__,_|_| |_|_| |_|\__, | 
                                  |___/ 
`) + "\n"

	// While not done, show spinner and progress text.
	if !m.done {
		s += m.spinner.View() + color.YellowString("  Cleaning in progress...\n\n")
		return s
	}

	// Once done, show logs or error.
	if m.err != nil {
		s += color.RedString("Error: %v\n\n", m.err)
	} else {
		for _, line := range m.logs {
			s += line + "\n"
		}
		s += color.GreenString("\nCleanup complete! 🚀\n")
		s += color.WhiteString("\nPress any key to exit...")
	}

	return s
}

// -------------------
// doCleanupCmd: perform cleanup in a goroutine
// -------------------
func doCleanupCmd() tea.Cmd {
	return func() tea.Msg {
		logs, err := runCleanup()
		if err != nil {
			return errMsg(err)
		}
		return logMsg(logs)
	}
}

// -------------------
// runCleanup: actual file deletion process
// -------------------
func runCleanup() ([]string, error) {
	var logs []string

	// Get the current user's TEMP directory.
	userTemp := os.Getenv("TEMP")
	if userTemp == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error retrieving home directory: %v", err)
		}
		userTemp = filepath.Join(homeDir, "AppData", "Local", "Temp")
	}

	// Directories to clean.
	dirs := []string{
		userTemp,              // Current user's TEMP directory
		`C:\Windows\Temp`,     // Windows system temporary directory
		`C:\Windows\Prefetch`, // Windows Prefetch directory
	}

	// Iterate through directories.
	for _, dir := range dirs {
		logs = append(logs, color.MagentaString("Cleaning directory: %s", dir))

		items, err := ioutil.ReadDir(dir)
		if err != nil {
			logs = append(logs, color.RedString("Error reading %s: %v", dir, err))
			continue
		}

		var failures int
		for _, item := range items {
			fullPath := filepath.Join(dir, item.Name())
			if err := os.RemoveAll(fullPath); err != nil {
				failures++
				logs = append(logs, color.RedString("Failed to remove %s: %v", fullPath, err))
			} else {
				logs = append(logs, color.GreenString("Successfully removed %s", fullPath))
			}
		}

		if failures > 0 {
			logs = append(logs, color.RedString("Some files in %s could not be removed.", dir))
		} else {
			logs = append(logs, color.GreenString("Finished cleaning directory: %s", dir))
		}
		logs = append(logs, "") // Blank line for spacing
	}

	return logs, nil
}

// -------------------
// Main entry point
// -------------------
func main() {
	s := spinner.New()
	s.Spinner = spinner.Dot

	p := tea.NewProgram(model{
		spinner: s,
	})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
