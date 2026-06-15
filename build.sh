#!/bin/bash

# 🔍 Statische Codeanalyse für Go Projekte
# Dieses Script führt verschiedene Go-Analysetools aus

echo "🔍 Statische Codeanalyse gestartet..."
echo "=========================================="

# Arbeitsverzeichnis wechseln
cd "$(dirname "$0")" || exit 1

echo ""
echo "📋 1. Go Modules überprüfen..."
echo "------------------------------"
go mod tidy
go mod verify

echo ""
echo "🔧 2. Build-Überprüfung..."
echo "-------------------------"

VERSION=$(git describe --tags)

if go build -ldflags "-X main.version=$VERSION" ./...; then
    echo "✅ Build erfolgreich"
else
    echo "❌ Build fehlgeschlagen"
    exit 1
fi

echo ""
echo "🎨 3. Code-Formatierung (gofmt)..."
echo "----------------------------------"
GOFMT_OUTPUT=$(gofmt -d .)
if [ -n "$GOFMT_OUTPUT" ]; then
    echo "⚠️  Code-Formatierung Probleme gefunden:"
    echo "$GOFMT_OUTPUT"
else
    echo "✅ Code-Formatierung OK"
fi

echo ""
echo "📦 4. Imports (goimports)..."
echo "---------------------------"
if command -v goimports >/dev/null 2>&1; then
    GOIMPORTS_OUTPUT=$(goimports -d .)
    if [ -n "$GOIMPORTS_OUTPUT" ]; then
        echo "⚠️  Import Probleme gefunden:"
        echo "$GOIMPORTS_OUTPUT"
    else
        echo "✅ Imports OK"
    fi
else
    echo "ℹ️  goimports nicht installiert (go install golang.org/x/tools/cmd/goimports@latest)"
fi

echo ""
echo "🔍 5. Go Vet (Standard-Analyse)..."
echo "---------------------------------"
if go vet ./...; then
    echo "✅ go vet: Keine Probleme gefunden"
else
    echo "❌ go vet: Probleme gefunden"
fi

echo ""
echo "🚨 6. Golangci-lint (Erweiterte Analyse + Security)..."
echo "---------------------------------------------------"
if command -v golangci-lint >/dev/null 2>&1; then
    golangci-lint run --config golangci.yml
    LINT_EXIT_CODE=$?
    if [ $LINT_EXIT_CODE -eq 0 ]; then
        echo "✅ golangci-lint: Keine Probleme gefunden"
    else
        echo "⚠️  golangci-lint: Issues gefunden (siehe oben)"
    fi
else
    echo "ℹ️  golangci-lint nicht installiert (brew install golangci-lint)"
fi

echo ""
echo "🛡️  7. Vulnerability Check (govulncheck)..."
echo "--------------------------------------------"
if command -v govulncheck >/dev/null 2>&1; then
    govulncheck ./...
    VULN_EXIT_CODE=$?
    if [ $VULN_EXIT_CODE -eq 0 ]; then
        echo "✅ govulncheck: Keine bekannten Vulnerabilities gefunden"
    else
        echo "⚠️  govulncheck: Vulnerabilities gefunden"
    fi
else
    echo "ℹ️  govulncheck nicht installiert"
    echo "📋 Installation: go install golang.org/x/vuln/cmd/govulncheck@latest"
fi

echo ""
echo "🧪 8. Tests..."
echo "-------------"
if go test ./...; then
    echo "✅ Tests erfolgreich"
else
    echo "❌ Tests fehlgeschlagen"
fi

echo ""
echo "📊 9. Test-Coverage..."
echo "---------------------"
go test -race -coverprofile=coverage.out ./... >/dev/null 2>&1
if [ -f coverage.out ]; then
    COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}')
    echo "📈 Test-Coverage: $COVERAGE"
    rm coverage.out
else
    echo "ℹ️  Keine Coverage-Daten verfügbar"
fi

echo ""
echo "📋 10. Code-Statistiken..."
echo "-------------------------"
echo "📁 Go-Dateien: $(find . -name "*.go" | wc -l)"
echo "📏 Zeilen Code: $(find . -name "*.go" -exec cat {} \; | wc -l)"
echo "📦 Packages: $(go list ./... | wc -l)"

echo ""
echo "=========================================="
echo "🏁 Statische Analyse abgeschlossen!"
echo ""
echo "💡 Tipps:"
echo "  - Behebe kritische Issues aus golangci-lint"
echo "  - Überprüfe Vulnerabilities mit govulncheck"
echo "  - Achte auf mindestens 80% Test-Coverage"
echo "  - Verwende 'go fmt' für einheitliche Formatierung"
echo "  - Version: $VERSION"