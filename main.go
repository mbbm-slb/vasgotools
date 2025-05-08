package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    // Define a command-line flag for the folder path
    folderPath := flag.String("path", "", "Path to the folder (defaults to current working directory)")
    flag.Parse()

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
}