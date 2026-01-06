@echo off
REM Cross-platform build script for Go projects
REM Builds executables for Windows, Linux, and macOS

echo Building Go project for multiple platforms...
echo.

REM Check if go.mod exists
if not exist "go.mod" (
    echo Error: go.mod not found in current directory
    exit /b 1
)

REM Get module name from go.mod and extract the last part as binary name
for /f "tokens=2" %%i in ('findstr /b "module " go.mod') do set MODULE_NAME=%%i
for %%i in ("%MODULE_NAME:/=" "%") do set BINARY_NAME=%%~nxi

echo Module: %MODULE_NAME%
echo Binary name: %BINARY_NAME%
echo.

REM Create output directory
if not exist "bin" mkdir bin

REM Build for Windows (amd64)
echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o bin\%BINARY_NAME%-windows-amd64.exe
if %errorlevel% neq 0 (
    echo Failed to build for Windows amd64
    exit /b 1
)

REM Build for Linux (amd64)
echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o bin\%BINARY_NAME%-linux-amd64
if %errorlevel% neq 0 (
    echo Failed to build for Linux amd64
    exit /b 1
)

REM Build for macOS (amd64 - Intel)
echo Building for macOS (amd64 - Intel)...
set GOOS=darwin
set GOARCH=amd64
go build -o bin\%BINARY_NAME%-darwin-amd64
if %errorlevel% neq 0 (
    echo Failed to build for macOS amd64
    exit /b 1
)

REM Build for macOS (arm64 - Apple Silicon)
echo Building for macOS (arm64 - Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build -o bin\%BINARY_NAME%-darwin-arm64
if %errorlevel% neq 0 (
    echo Failed to build for macOS arm64
    exit /b 1
)

echo.
echo Build completed successfully!
echo Binaries are located in the bin/ directory:
echo   - %BINARY_NAME%-windows-amd64.exe (Windows)
echo   - %BINARY_NAME%-linux-amd64 (Linux)
echo   - %BINARY_NAME%-darwin-amd64 (macOS Intel)
echo   - %BINARY_NAME%-darwin-arm64 (macOS Apple Silicon)
echo.
