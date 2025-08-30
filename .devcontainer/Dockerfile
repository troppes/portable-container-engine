FROM mcr.microsoft.com/devcontainers/go:1-1.21-bullseye

# Install additional packages
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*
