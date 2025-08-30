package util

import (
	"errors"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		substr   string
		expected bool
	}{
		{
			name:     "simple match",
			s:        "hello world",
			substr:   "world",
			expected: true,
		},
		{
			name:     "case insensitive match",
			s:        "Hello World",
			substr:   "world",
			expected: true,
		},
		{
			name:     "no match",
			s:        "hello world",
			substr:   "goodbye",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.s, tt.substr); got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Run("nil error should not panic", func(t *testing.T) {
		Must(nil)
	})

	t.Run("non-nil error should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Must() did not panic with error as expected")
			}
		}()

		Must(errors.New("test error"))
	})
}
