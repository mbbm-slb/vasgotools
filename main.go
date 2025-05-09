package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	openVSCodeBatchFile = "open_vscode.bat"
	openVSCodeShellFile = "open_vscode.sh"
)

func main() {
	// Ensure a subcommand is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <command> [options]")
		fmt.Println("Available commands:")
		fmt.Println("  generate-work    Generate a go workspace (i.e. a go.work file)")
		fmt.Println("  generate-app     Create a new Go application")
		os.Exit(1)
	}

	// Determine the subcommand
	switch os.Args[1] {
	case "generate-work":
		generateWorkCommand(os.Args[2:])
	case "generate-app":
		generateAppCommand(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Use 'go run main.go' for usage.")
		os.Exit(1)
	}
}

func generateWorkCommand(args []string) {
	// Define a flag set for the "generate-work" command
	fs := flag.NewFlagSet("generate-work", flag.ExitOnError)
	folderPath := fs.String("path", "", "Path to the folder (defaults to current working directory)")
	fs.Parse(args)

	// Check for optional flags
	noGit, noCode := parseOptionalFlags(fs.Args())

	// Use the current working directory if no path is provided
	err := setDefaultFolderPath(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Slice to store relative paths of subfolders containing go.mod
	var goModFolders []string

	// Walk through the directory
	err = filepath.Walk(*folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current item is a file named "go.mod"
		if info.Name() == "go.mod" {
			// Get the relative path of the folder containing go.mod
			relativePath, err := filepath.Rel(*folderPath, filepath.Dir(path))
			if err != nil {
				return err
			}
			goModFolders = append(goModFolders, relativePath)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
		return
	}

	// Print the collected relative paths
	fmt.Println("Subfolders containing go.mod:")
	for _, folder := range goModFolders {
		fmt.Println(folder)
	}

	// Run the "go work init" command with the relative paths
	if len(goModFolders) > 0 {
		args := append([]string{"work", "init"}, goModFolders...)
		cmd := exec.Command("go", args...)
		cmd.Dir = *folderPath // Set the working directory to the root folder
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("Running command:", cmd.String())
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running 'go work init':", err)
			return
		}

		fmt.Println("go.work file created successfully.")
	} else {
		fmt.Println("No subfolders with go.mod found. No go.work file created.")
	}

	// Initialize a Git repository (if not suppressed)
	if !noGit {
		err = initializeGitRepository(*folderPath)
		if err != nil {
			fmt.Println("Error initializing Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized successfully.")
	} else {
		fmt.Println("Git repository initialization skipped.")
	}

	// Create the open_vscode.bat file (if not suppressed)
	if !noCode {
		err = createOpenVSCodeFile(*folderPath)
		if err != nil {
			fmt.Println("Error creating open_vscode file:", err)
			return
		}

		// Execute the open_vscode file
		err = executeOpenVSCodeFile(*folderPath)
		if err != nil {
			fmt.Println("Error executing open_vscode file:", err)
			return
		}
		fmt.Println("Visual Studio Code opened successfully.")
	} else {
		fmt.Println("Creation and execution of open_vscode.bat skipped.")
	}
}

func generateAppCommand(args []string) {
	// Define a flag set for the "generate-app" command
	fs := flag.NewFlagSet("generate-app", flag.ExitOnError)
	folderPath := fs.String("path", "", "Path to create the application folder (defaults to current working directory)")
	fs.Parse(args)

	// Ensure the application name is provided as the first positional argument
	if fs.NArg() < 1 {
		fmt.Println("Error: Application name is required.")
		fmt.Println("Usage: vasgotools.exe generate-app <name> [--path <path>] [nogit] [nocode] [nomain]")
		os.Exit(1)
	}
	appName := fs.Arg(0)

	// Check for optional flags
	noGit, noCode := parseOptionalFlags(fs.Args()[1:])
	noMain := false
	for _, arg := range fs.Args()[1:] {
		if arg == "nomain" {
			noMain = true
		}
	}

	// Use the current working directory if no path is provided
	err := setDefaultFolderPath(folderPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create the application folder
	appFolder := filepath.Join(*folderPath, appName)
	err = os.MkdirAll(appFolder, 0755)
	if err != nil {
		fmt.Println("Error creating application folder:", err)
		return
	}

	// Run the "go mod init" command
	appFullName := "github.com/muellerbbm-vas/" + appName
	cmd := exec.Command("go", "mod", "init", appFullName)
	cmd.Dir = appFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running command:", cmd.String())
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running 'go mod init':", err)
		return
	}

	// Create a main.go file with a Hello World example (if not suppressed)
	if !noMain {
		mainGoContent := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`
		mainGoFilePath := filepath.Join(appFolder, "main.go")
		err = os.WriteFile(mainGoFilePath, []byte(mainGoContent), 0644)
		if err != nil {
			fmt.Println("Error creating main.go file:", err)
			return
		}
		fmt.Printf("A main.go file with a Hello World example has been created in '%s'.\n", mainGoFilePath)
	} else {
		fmt.Println("Creation of main.go file skipped.")
	}

	// Initialize a Git repository (if not suppressed)
	if !noGit {
		err = initializeGitRepository(appFolder)
		if err != nil {
			fmt.Println("Error initializing Git repository:", err)
			return
		}
		fmt.Println("Git repository initialized successfully.")
	} else {
		fmt.Println("Git repository initialization skipped.")
	}

	// Create the open_vscode.bat file (if not suppressed)
	if !noCode {
		err = createOpenVSCodeFile(appFolder)
		if err != nil {
			fmt.Println("Error creating open_vscode file:", err)
			return
		}

		// Execute the open_vscode file
		err = executeOpenVSCodeFile(appFolder)
		if err != nil {
			fmt.Println("Error executing open_vscode file:", err)
			return
		}
		fmt.Println("Visual Studio Code opened successfully.")
	} else {
		fmt.Println("Creation and execution of open_vscode.bat skipped.")
	}

	fmt.Printf("Application '%s' created successfully in folder '%s'.\n", appFullName, appFolder)
}

// parseOptionalFlags parses the optional "nogit" and "nocode" flags from the arguments.
func parseOptionalFlags(args []string) (bool, bool) {
	noGit := false
	noCode := false
	for _, arg := range args {
		if arg == "nogit" {
			noGit = true
		} else if arg == "nocode" {
			noCode = true
		}
	}
	return noGit, noCode
}

// setDefaultFolderPath sets the folder path to the current working directory if it is empty.
func setDefaultFolderPath(folderPath *string) error {
	if *folderPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current working directory: %w", err)
		}
		*folderPath = cwd
	}
	return nil
}

func initializeGitRepository(folderPath string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = folderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Initializing Git repository...")
	return cmd.Run()
}

func createOpenVSCodeBatchFile(folderPath string) error {
	batchFilePath := filepath.Join(folderPath, openVSCodeBatchFile)
	batchFileContent := "code . | exit 0\n"
	return os.WriteFile(batchFilePath, []byte(batchFileContent), 0644)
}

func executeOpenVSCodeBatchFile(folderPath string) error {
	batchFilePath := filepath.Join(folderPath, openVSCodeBatchFile)
	cmd := exec.Command("cmd", "/C", batchFilePath)
	cmd.Dir = folderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Opening Visual Studio Code...")
	return cmd.Run()
}

func createOpenVSCodeShellScript(folderPath string) error {
	scriptFilePath := filepath.Join(folderPath, openVSCodeShellFile)
	scriptContent := "#!/bin/bash\ncode . || exit 0\n"
	err := os.WriteFile(scriptFilePath, []byte(scriptContent), 0755) // Make the script executable
	if err != nil {
		return fmt.Errorf("error creating open_vscode.sh: %w", err)
	}
	return nil
}

func executeOpenVSCodeShellScript(folderPath string) error {
	scriptFilePath := filepath.Join(folderPath, openVSCodeShellFile)
	cmd := exec.Command("bash", scriptFilePath)
	cmd.Dir = folderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Opening Visual Studio Code...")
	return cmd.Run()
}

func createOpenVSCodeFile(folderPath string) error {
	if runtime.GOOS == "windows" {
		return createOpenVSCodeBatchFile(folderPath)
	} else {
		return createOpenVSCodeShellScript(folderPath)
	}
}

func executeOpenVSCodeFile(folderPath string) error {
	if runtime.GOOS == "windows" {
		return executeOpenVSCodeBatchFile(folderPath)
	} else {
		return executeOpenVSCodeShellScript(folderPath)
	}
}
