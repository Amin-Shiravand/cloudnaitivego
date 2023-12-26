@echo off

REM Navigate to the root directory of the Go project
echo Current directory: %CD%

REM Run the Docker build command
docker build -t cloudnative -f  docker\build.Dockerfile .