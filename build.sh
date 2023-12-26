#!/bin/bash

# Navigate to the root directory of the Go project
echo "Current directory after navigation: $(pwd)"


# Run the Docker build command
docker build -t cloudnative -f docker/build.Dockerfile .
