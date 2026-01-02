@echo off
setlocal enabledelayedexpansion

REM Statische Codeanalyse fuer Go Projekte
REM Dieses Script fuehrt verschiedene Go-Analysetools aus

echo [*] Statische Codeanalyse gestartet...
echo ==========================================

REM Arbeitsverzeichnis wechseln
cd /d "%~dp0"

echo.
echo [1] Go Modules ueberpruefen...
echo ------------------------------
go mod tidy
go mod verify

echo.
echo [2] Build-Ueberpruefung...
echo -------------------------
go build ./...
if errorlevel 1 (
    echo [X] Build fehlgeschlagen
    exit /b 1
) else (
    echo [OK] Build erfolgreich
)

echo.
echo [3] Code-Formatierung ^(gofmt^)...
echo ----------------------------------
gofmt -d . > gofmt_temp.txt 2>&1
if exist gofmt_temp.txt (
    for %%A in (gofmt_temp.txt) do set size=%%~zA
    if !size! gtr 0 (
        echo [!] Code-Formatierung Probleme gefunden:
        type gofmt_temp.txt
    ) else (
        echo [OK] Code-Formatierung OK
    )
    del gofmt_temp.txt
)

echo.
echo [4] Imports ^(goimports^)...
echo ---------------------------
where goimports >nul 2>&1
if errorlevel 1 (
    echo [i] goimports nicht installiert ^(go install golang.org/x/tools/cmd/goimports@latest^)
) else (
    goimports -l . > goimports_temp.txt 2>&1
    if exist goimports_temp.txt (
        for %%A in (goimports_temp.txt) do set size=%%~zA
        if !size! gtr 0 (
            echo [!] Import Probleme in folgenden Dateien gefunden:
            type goimports_temp.txt
        ) else (
            echo [OK] Imports OK
        )
        del goimports_temp.txt
    )
)

echo.
echo [5] Go Vet ^(Standard-Analyse^)...
echo ---------------------------------
go vet ./...
if errorlevel 1 (
    echo [X] go vet: Probleme gefunden
) else (
    echo [OK] go vet: Keine Probleme gefunden
)

echo.
echo [6] Golangci-lint ^(Erweiterte Analyse + Security^)...
echo ---------------------------------------------------
where golangci-lint >nul 2>&1
if errorlevel 1 (
    echo [i] golangci-lint nicht installiert ^(https://golangci-lint.run/welcome/install/^)
) else (
    golangci-lint run --config golangci_win.yml
    if errorlevel 1 (
        echo [!] golangci-lint: Issues gefunden ^(siehe oben^)
    ) else (
        echo [OK] golangci-lint: Keine Probleme gefunden
    )
)

echo.
echo [7] Vulnerability Check ^(govulncheck^)...
echo --------------------------------------------
where govulncheck >nul 2>&1
if errorlevel 1 (
    echo [i] govulncheck nicht installiert
    echo [i] Installation: go install golang.org/x/vuln/cmd/govulncheck@latest
) else (
    govulncheck ./...
    if errorlevel 1 (
        echo [!] govulncheck: Vulnerabilities gefunden
    ) else (
        echo [OK] govulncheck: Keine bekannten Vulnerabilities gefunden
    )
)

echo.
echo [8] Tests...
echo -------------
go test ./...
if errorlevel 1 (
    echo [X] Tests fehlgeschlagen
) else (
    echo [OK] Tests erfolgreich
)

echo.
echo [9] Test-Coverage...
echo ---------------------
go test -coverprofile=coverage.out ./...
if errorlevel 1 (
    echo [!] Test-Coverage Befehl fehlgeschlagen
) else (
    if exist coverage.out (
        for /f "tokens=3" %%i in ('go tool cover -func^=coverage.out ^| findstr /R "total:"') do set COVERAGE=%%i
        if defined COVERAGE (
            echo [*] Test-Coverage: !COVERAGE!
        ) else (
            echo [*] Coverage-Daten generiert
        )
        del coverage.out
    ) else (
        echo [i] Keine Coverage-Daten verfuegbar
    )
)

echo.
echo [10] Code-Statistiken...
echo -------------------------
set GO_FILES=0
for /f %%a in ('dir /s /b *.go 2^>nul ^| find /c /v ""') do set GO_FILES=%%a
echo [*] Go-Dateien: !GO_FILES!

set GO_LINES=0
for /f %%a in ('type *.go 2^>nul ^| find /c /v ""') do set GO_LINES=%%a
echo [*] Zeilen Code: !GO_LINES!

for /f %%a in ('go list ./... 2^>nul ^| find /c /v ""') do set PACKAGES=%%a
echo [*] Packages: !PACKAGES!

echo.
echo ==========================================
echo [*] Statische Analyse abgeschlossen!
echo.
echo [*] Tipps:
echo   - Behebe kritische Issues aus golangci-lint
echo   - Ueberpruefe Vulnerabilities mit govulncheck
echo   - Achte auf mindestens 80%% Test-Coverage
echo   - Verwende 'go fmt' fuer einheitliche Formatierung
