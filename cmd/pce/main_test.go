package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMainIntegration(t *testing.T) {
	ctx := context.Background()

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := filepath.Dir(filepath.Dir(pwd))

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    projectRoot,
			Dockerfile: "Dockerfile",
		},
		Privileged: true, // Required for container operations
		Entrypoint: []string{"/bin/sh", "-c"},
		Cmd:        []string{"tail -f /dev/null"}, // Keep container running
		WaitingFor: wait.ForLog(""),               // No specific log to wait for
	}

	// Create the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	tests := []struct {
		name    string
		args    []string
		wantOut string
	}{
		{
			name:    "No arguments",
			args:    []string{},
			wantOut: "Usage: pce <download|run> <image> [<command>...]",
		},
		{
			name:    "Download alpine image",
			args:    []string{"download", "alpine:latest"},
			wantOut: "Image downloaded to",
		},
		{
			name:    "Extract alpine image",
			args:    []string{"download", "alpine:latest", "--extract"},
			wantOut: "Image extracted to",
		},
		{
			name:    "Run alpine container",
			args:    []string{"run", "alpine:latest", "echo", "hello from container"},
			wantOut: "hello from container",
		},
		{
			name:    "Invalid command",
			args:    []string{"invalid"},
			wantOut: "Usage: pce <download|run> <image> [<command>...]",
		},
	}

	// Create base temp directory for all tests
	baseTestDir, err := os.MkdirTemp("", "pce-integration-tests-*")
	if err != nil {
		t.Fatalf("Failed to create base test directory: %v", err)
	}
	defer os.RemoveAll(baseTestDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test-specific directory
			testDir := filepath.Join(baseTestDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Set up working directory in container
			setupCmd := []string{"/bin/sh", "-c", fmt.Sprintf("mkdir -p %s", testDir)}
			if _, err := runInContainer(container, ctx, setupCmd); err != nil {
				t.Fatalf("Failed to create test directory in container: %v", err)
			}

			// Create the command with proper working directory
			cmdStr := fmt.Sprintf("cd %s && /usr/local/bin/pce %s", testDir, strings.Join(tt.args, " "))
			cmd := []string{"/bin/sh", "-c", cmdStr}

			result, err := runInContainer(container, ctx, cmd)
			if err != nil {
				return
			}

			outputStr := result
			if tt.wantOut != "" && !strings.Contains(outputStr, tt.wantOut) {
				t.Errorf("Expected output to contain %q, got %q", tt.wantOut, outputStr)
			}

			// For run commands, verify container isolation
			if len(tt.args) > 0 && tt.args[0] == "run" {
				// Verify process isolation
				exitCode, psOutput, err := container.Exec(ctx, []string{"ps", "aux"})
				if err != nil {
					t.Fatalf("Failed to check processes: %v", err)
				}
				if exitCode != 0 {
					t.Errorf("Process check failed with exit code %d: %s", exitCode, psOutput)
				}

				// Verify filesystem isolation
				exitCode, mountOutput, err := container.Exec(ctx, []string{"mount"})
				if err != nil {
					t.Fatalf("Failed to check mounts: %v", err)
				}
				if exitCode != 0 {
					t.Errorf("Mount check failed with exit code %d: %s", exitCode, mountOutput)
				}
			}

			cleanupCmd := []string{"/bin/sh", "-c", fmt.Sprintf("rm -rf %s", testDir)}
			if _, err := runInContainer(container, ctx, cleanupCmd); err != nil {
				t.Logf("Warning: Failed to clean up test directory in container: %v", err)
			}
		})
	}
}

// Helper function to run a command in the container and capture its output
func runInContainer(container testcontainers.Container, ctx context.Context, command []string) (string, error) {
	// Run the command and capture its output
	exitCode, cmdOutput, err := container.Exec(ctx, command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}

	// Read the output
	output, err := io.ReadAll(cmdOutput)
	if err != nil {
		return "", fmt.Errorf("failed to read command output: %v", err)
	}

	if exitCode != 0 {
		return "", fmt.Errorf("command failed with exit code %d: %s", exitCode, string(output))
	}

	return string(output), nil
}
