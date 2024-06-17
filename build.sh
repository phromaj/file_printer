#!/bin/bash

# Set the program name
PROGRAM_NAME="codeprinter"

# Create the bin directory if it doesn't exist
mkdir -p bin

# Build for M1 Mac (darwin/arm64)
echo "Building for M1 Mac (darwin/arm64)..."
GOOS=darwin GOARCH=arm64 go build -o "bin/${PROGRAM_NAME}-mac-arm64" main.go

# Build for x64 Windows (windows/amd64)
echo "Building for x64 Windows (windows/amd64)..."
GOOS=windows GOARCH=amd64 go build -o "bin/${PROGRAM_NAME}-windows-amd64.exe" main.go

# Build for x64 Linux (linux/amd64)
echo "Building for x64 Linux (linux/amd64)..."
GOOS=linux GOARCH=amd64 go build -o "bin/${PROGRAM_NAME}-linux-amd64" main.go

echo "Builds completed and placed in the bin folder."