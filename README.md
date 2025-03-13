# WindowsTempCleanify

**WindowsTempCleanify** is a command-line tool and script written in Go that automates the cleanup of temporary files on Windows. It cleans your user TEMP directory, Windows system temporary directory, and Prefetch directory. 

---
![cleanify.png](assets/cleanify.png)
---


## Features

- **Automated Cleanup:**  
  Removes files from:
  - The user's TEMP `%temp%` directory
  - `C:\Windows\Temp`
  - `C:\Windows\Prefetch`

- **Interactive TUI:**  
  Uses Bubble Tea to display:
  - A colorful banner
  - A progress spinner
  - Detailed logs for each deletion operation

- **Detailed Metrics:**  
  Provides per-directory and overall summaries:
  - Total files processed
  - Number of successful deletions
  - Number of failures

- **Exit Prompt:**  
  After cleanup, the TUI displays a prompt ("Press any key to exit...") and waits for your input before quitting.

  <!-- download the script -->


## Prerequisites

- **Go 1.18+** – [Download Go](https://golang.org/dl/)
- **Windows OS** – Designed to work on Windows (requires admin privileges for system directories)
- **Terminal with ANSI Color Support** – For the best experience, use Windows Terminal or PowerShell.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/WindowsTempCleanify.git
   cd WindowsTempCleanify

## Build application
```go build -o name.exe```


### Find the build
<a href="cleanify.exe" download="cleanify.exe">cleanify.exe</a>