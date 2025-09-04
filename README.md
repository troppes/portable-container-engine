# Portable Container Engine (PCE)

A lightweight, cross-platform tool to download and extract Docker containers without requiring Docker installation. PCE furthermore provides a simple and efficient way to run containers on Linux while maintaining a minimal footprint.

!Warning this is an repo to play around with containers, there are no proper safeguards used against container breaches. Do not run this with unsafe containers!

## Features

- 🐳 Run Docker containers without Docker installation
- 🔄 Partial Cross-platform support (Linux, Windows, macOS)
- 📦 Download and extract Docker images
- 🔒 Containerization using native Linux OS features
- 🚀 Simple and lightweight implementation

## Quick Start

### Installation

Download the latest binary for your platform from the releases page or build from source:

```bash
make build
```

The binary will be available in the `bin/` directory.

### Basic Usage

1. Download a Docker image:
```bash
pce download alpine:latest
```

Downloads are safed inside a folder called `pce-download`

2. Run a container:
```bash
pce run alpine:latest /bin/sh
```

If no command is specified the default command of the container is used:
```bash
go run cmd/pce/main.go run ghcr.io/patrickdappollonio/docker-http-server
```

## Development

### Prerequisites
- Go 1.25 or higher
- Make
- VS Code with Remote-Containers extension (for development) or Linux OS
- Docker (for development environment)

### Setup Development Environment

1. Clone the repository
2. Open in VS Code
3. Click "Reopen in Container" when prompted

### Building

```bash
# Local build
make build

# Build all platforms
make build-all

# Clean and test
make clean
make test
```

### Development in Docker

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
│   └── pce/       # Main application code
├── internal/      # Private application code
│   ├── image/     # Image management and Docker registry client
│   ├── runtime/   # Platform-specific container runtime implementations
│   └── util/      # Shared utility functions
└── Makefile      # Build automation
```

## Implementation Details

- **Image Management**: Handles downloading and extracting Docker images from registries
- **Container Runtime**: Implementation for container isolation under Linux
- **Partial Cross-Platform Support**: Download and Extract container images on all platforms

## Current Limitations

- Limited namespace support (not all Linux namespaces are implemented)
- No cgroups support yet
- No network namespace support

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'Add new feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Submit a pull request

## Special Thanks

This project was inspired by and builds upon the work of Liz Rice's [containers-from-scratch](https://github.com/lizrice/containers-from-scratch) project.

## License

[MIT License](LICENSE)
