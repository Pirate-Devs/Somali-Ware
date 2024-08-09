@echo off
setlocal enabledelayedexpansion

cd /d %~dp0

python -m main

pause
exit