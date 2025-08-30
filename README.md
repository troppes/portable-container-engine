# Portable Container Engine

A lightweight tool to download and run Docker containers without Docker.

## Quick Start

```bash
# Download and extract an image
go run cmd/pce/main.go download alpine:latest -x

# Run a container
go run cmd/pce/main.go run ubuntu:latest /bin/bash
```

## Development

### Prerequisites
- VS Code with Remote-Containers extension
- Docker

### Setup Development Environment

1. Clone the repository
2. Open in VS Code
3. Click "Reopen in Container" when prompted

### Build

```bash
# Local build
make build

# Build all platforms
make build-all

# Clean and test
make clean
make test
```

### Run in Docker

```bash
# Build container
docker-compose build

# Run
docker-compose up
```

## Project Structure

```
├── bin/           # Compiled binaries
├── cmd/           # Application entrypoints
├── internal/      # Private application code
│   ├── image/     # Image management
│   ├── runtime/   # Container runtime
│   └── util/      # Utility functions
└── Makefile       # Build automation
```

## Not yet supported:

- not all namespaces
- cGroups
- does not run the docker intended command

## Special Thanks

Liz Rice for the idea and base for this project: [https://github.com/lizrice/containers-from-scratch](https://github.com/lizrice/containers-from-scratch)
