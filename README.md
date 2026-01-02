# VasGoTools

A command-line utility tool for managing Go projects, workspaces, applications, and libraries with built-in static analysis support.

## Features

- 🚀 **Quick Project Setup** - Create Go applications and libraries with a single command
- 📦 **Workspace Management** - Automatically generate Go workspaces with all submodules
- 🔍 **Static Analysis** - Integrated analyze scripts and golangci-lint configuration
- 🛠️ **Build Scripts** - Auto-generated build scripts for Windows and Linux/macOS
- 🔧 **VS Code Integration** - Automatic VS Code workspace opening
- 📝 **Template-based** - Consistent project structure with embedded templates
- 🔐 **Git Integration** - Automatic repository initialization and submodule management

## Installation

### Prerequisites

- Go 1.21 or higher
- Git (optional, for repository initialization)
- VS Code (optional, for IDE integration)

### Build from Source

```bash
git clone <repository-url>
cd VasGoTools
go build -o vasgotools.exe
```

## Usage

```bash
vasgotools.exe <command> [options]
```

### Available Commands

| Command | Description |
|---------|-------------|
| `work`  | Generate a Go workspace (go.work file) |
| `app`   | Create a new Go application |
| `lib`   | Create a new Go library |

### Global Options

| Option | Description |
|--------|-------------|
| `--path <path>` | Specify the folder path (defaults to current working directory) |
| `--module-prefix <prefix>` | Specify the module prefix (default: none) |
| `nogit` | Skip Git repository initialization |
| `nocode` | Skip creation and execution of the open_vscode file |
| `nomain` | Skip creation of the main.go file (app command only) |

### Module Prefix Shortcuts

- `vas` → `github.com/muellerbbm-vas/`
- `slb` → `github.com/mbbm-slb/`

## Examples

### Create a Go Workspace

Generate a workspace that includes all Go modules in subdirectories:

```bash
vasgotools.exe work --path "C:\projects\myworkspace"
```

This will:
- Scan for all `go.mod` files in subdirectories
- Create a `go.work` file with all found modules
- Initialize a Git repository (optional)
- Open VS Code (optional)

### Create a New Application

```bash
vasgotools.exe app myapp --path "C:\projects"
```

This creates a new application with:
- `go.mod` file
- `main.go` from template
- `build.bat` and `build.sh` scripts
- `analyze.bat` and `analyze.sh` scripts
- `golangci.yml` and `golangci_win.yml` configurations
- `open_vscode.bat` and `open_vscode.sh` scripts
- Git repository with initial commit

### Create a New Library

```bash
vasgotools.exe lib mylib --module-prefix vas
```

Creates a library with:
- Module name: `github.com/muellerbbm-vas/mylib`
- No `main.go` file (library only)
- All analysis and build scripts
- Git repository initialized

### Advanced Examples

Create an app without Git and VS Code integration:
```bash
vasgotools.exe app myapp nogit nocode
```

Create an app with custom module prefix:
```bash
vasgotools.exe app myapp --module-prefix "github.com/myorg/"
```

Create an app without main.go:
```bash
vasgotools.exe app myapp nomain
```

## Project Structure

### Recommended Workspace Structure

```
<workspace-root>/
├── go.work              # Go workspace file
├── app1/                # Application 1
│   ├── go.mod
│   ├── main.go
│   ├── build.bat
│   ├── build.sh
│   ├── analyze.bat
│   ├── analyze.sh
│   ├── golangci.yml
│   └── golangci_win.yml
├── app2/                # Application 2
│   └── ...
└── ext/                 # External libraries
    ├── lib1/
    │   ├── go.mod
    │   └── ...
    └── lib2/
        └── ...
```

### Generated Files

When creating an app or lib, the following files are automatically generated:

| File | Description | Platform |
|------|-------------|----------|
| `go.mod` | Go module file | All |
| `main.go` | Main application file (apps only) | All |
| `build.bat` | Build script | Windows |
| `build.sh` | Build script | Linux/macOS |
| `analyze.bat` | Static analysis script | Windows |
| `analyze.sh` | Static analysis script | Linux/macOS |
| `golangci.yml` | Linter configuration | Linux/macOS |
| `golangci_win.yml` | Linter configuration | Windows |
| `open_vscode.bat` | VS Code launcher | Windows |
| `open_vscode.sh` | VS Code launcher | Linux/macOS |

## Static Analysis

Each generated project includes comprehensive static analysis scripts that run:

1. **Go Modules** - `go mod tidy` and `go mod verify`
2. **Build Check** - `go build ./...`
3. **Code Formatting** - `gofmt` checks
4. **Import Check** - `goimports` validation
5. **Go Vet** - Standard Go analysis
6. **golangci-lint** - Extended linting with security checks
7. **Vulnerability Check** - `govulncheck` for known CVEs
8. **Tests** - `go test ./...`
9. **Test Coverage** - Coverage report generation
10. **Code Statistics** - File counts and metrics

### Running Analysis

**Windows:**
```bash
.\analyze.bat
```

**Linux/macOS:**
```bash
./analyze.sh
```

## Build Scripts

### Windows (build.bat)

```bash
.\build.bat
```

### Linux/macOS (build.sh)

```bash
./build.sh
```

Build scripts automatically compile your application for the current platform.

## Configuration Files

### golangci-lint Configuration

The tool generates two linter configurations:

- **golangci.yml** - For Linux/macOS
- **golangci_win.yml** - For Windows

Both include comprehensive security checks with gosec and best practice linters.

## Git Integration

When Git integration is enabled (default), the tool will:

1. Initialize a Git repository
2. Add all generated files
3. Create an initial commit with message "init"
4. For workspaces: detect and add existing Git repositories as submodules

To skip Git initialization:
```bash
vasgotools.exe app myapp nogit
```

## VS Code Integration

By default, projects automatically open in VS Code after creation. The tool creates platform-specific scripts:

- **Windows:** `open_vscode.bat`
- **Linux/macOS:** `open_vscode.sh`

To skip VS Code integration:
```bash
vasgotools.exe app myapp nocode
```

## License

Copyright © 2026 Müller-BBM VibroAkustik Systeme GmbH. All rights reserved.

This software is proprietary and confidential. Unauthorized copying, distribution, modification, or use of this software, via any medium, is strictly prohibited without the express written permission of Müller-BBM VibroAkustik Systeme GmbH.

## Contributing

Contributions to this project are currently limited to employees of Müller-BBM VibroAkustik Systeme GmbH.

For internal contributors:
- Follow the company's coding standards and guidelines
- Ensure all code passes static analysis (`analyze.bat`/`analyze.sh`)
- Run tests before submitting changes
- Document all significant changes in CHANGELOG.txt

## Support

For issues and questions, please contact slb.
