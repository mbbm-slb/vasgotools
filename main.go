package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func main() {
    // Ensure a subcommand is provided
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <command> [options]")
        fmt.Println("Available commands:")
        fmt.Println("  generate-work  Generate a go.work file")
        // Add more commands here as needed
        os.Exit(1)
    }

    // Determine the subcommand
    switch os.Args[1] {
    case "generate-work":
        generateWorkCommand(os.Args[2:])
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

    // Use the current working directory if no path is provided
    if *folderPath == "" {
        cwd, err := os.Getwd()
        if err != nil {
            fmt.Println("Error getting current working directory:", err)
            return
        }
        *folderPath = cwd
    }

    // Slice to store relative paths of subfolders containing go.mod
    var goModFolders []string

    // Walk through the directory
    err := filepath.Walk(*folderPath, func(path string, info os.FileInfo, err error) error {
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
}