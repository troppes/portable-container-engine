package runtime

import (
	"fmt"
	"testing"
)

// mockRuntime implements ContainerRuntime interface for testing
type mockRuntime struct {
	runCalled           bool
	createProcessCalled bool
	lastImage           string
	lastCommand         []string
	lastPath            string
	shouldError         bool
}

func (m *mockRuntime) Run(image string, command []string) error {
	m.runCalled = true
	m.lastImage = image
	m.lastCommand = command
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func (m *mockRuntime) CreateChildProcess(path string, command []string) error {
	m.createProcessCalled = true
	m.lastPath = path
	m.lastCommand = command
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	return nil
}

func TestGetRuntime(t *testing.T) {
	runtime := GetRuntime()
	if runtime == nil {
		t.Error("GetRuntime() returned nil")
	}
}

func TestRuntimeInterface(t *testing.T) {
	tests := []struct {
		name        string
		image       string
		path        string
		command     []string
		shouldError bool
	}{
		{
			name:        "basic run command",
			image:       "ubuntu:latest",
			command:     []string{"echo", "hello"},
			shouldError: false,
		},
		{
			name:        "create process",
			path:        "/bin/bash",
			command:     []string{"-c", "echo hello"},
			shouldError: false,
		},
		{
			name:        "error case",
			image:       "invalid:image",
			command:     []string{},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRuntime{shouldError: tt.shouldError}

			// Test Run method
			if tt.image != "" {
				err := mock.Run(tt.image, tt.command)
				if tt.shouldError && err == nil {
					t.Error("expected error but got none")
				}
				if !tt.shouldError && err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !mock.runCalled {
					t.Error("Run was not called")
				}
				if mock.lastImage != tt.image {
					t.Errorf("wrong image used. got %s, want %s", mock.lastImage, tt.image)
				}
			}

			// Test CreateChildProcess method
			if tt.path != "" {
				err := mock.CreateChildProcess(tt.path, tt.command)
				if tt.shouldError && err == nil {
					t.Error("expected error but got none")
				}
				if !tt.shouldError && err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !mock.createProcessCalled {
					t.Error("CreateChildProcess was not called")
				}
				if mock.lastPath != tt.path {
					t.Errorf("wrong path used. got %s, want %s", mock.lastPath, tt.path)
				}
			}

			// Check command args for both methods
			if len(mock.lastCommand) != len(tt.command) {
				t.Errorf("wrong command length. got %d, want %d", len(mock.lastCommand), len(tt.command))
			}
			for i, arg := range mock.lastCommand {
				if arg != tt.command[i] {
					t.Errorf("wrong command arg at position %d. got %s, want %s", i, arg, tt.command[i])
				}
			}
		})
	}
}
