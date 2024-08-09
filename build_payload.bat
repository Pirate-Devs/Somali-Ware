@echo off
setlocal enabledelayedexpansion

cd /d %~dp0

pushd src

go build -ldflags "-s -w" -o ..\payload.exe

if errorlevel 1 (
    echo Build failed
    pause
    exit
)

popd

pause
exit