@echo off

FOR /F "tokens=* USEBACKQ" %%F IN (`git describe --tags`) DO (
SET GIT_VERSION_INFO=%%F
)
ECHO %GIT_VERSION_INFO%

:: preset the value of the variable gitVersionInfo in main.go with Version tag from git
go build -ldflags "-X main.gitVersionInfo=%GIT_VERSION_INFO%"