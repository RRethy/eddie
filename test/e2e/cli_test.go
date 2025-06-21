package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIBasics(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "help command",
			args:       []string{"--help"},
			wantOutput: "eddie",
			wantErr:    false,
		},
		{
			name:       "version flag",
			args:       []string{"--version"},
			wantOutput: "",
			wantErr:    true, // Version not implemented yet, should error
		},
		{
			name:       "no arguments",
			args:       []string{},
			wantOutput: "eddie",
			wantErr:    false,
		},
		{
			name:       "unknown command",
			args:       []string{"unknown"},
			wantOutput: "unknown command",
			wantErr:    true,
		},
		{
			name:       "view help",
			args:       []string{"view", "--help"},
			wantOutput: "Examine the contents",
			wantErr:    false,
		},
		{
			name:       "str_replace help",
			args:       []string{"str_replace", "--help"},
			wantOutput: "Replace all occurrences",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.True(t, err != nil || stderr != "", "Expected error but got none")
				if tt.wantOutput != "" {
					output := stdout + stderr
					assert.Contains(t, output, tt.wantOutput)
				}
				return
			}

			assert.NoError(t, err, "stderr: %s", stderr)
			if tt.wantOutput != "" {
				assert.Contains(t, stdout, tt.wantOutput)
			}
		})
	}
}

func TestCommandCompletion(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "view with partial args",
			args:    []string{"view"},
			wantErr: false, // Should show error message but not exit with error
		},
		{
			name:    "str_replace with partial args",
			args:    []string{"str_replace"},
			wantErr: false, // Should show error message but not exit with error
		},
		{
			name:    "str_replace with one arg",
			args:    []string{"str_replace", "file.txt"},
			wantErr: false, // Should show error message but not exit with error
		},
		{
			name:    "str_replace with two args",
			args:    []string{"str_replace", "file.txt", "old"},
			wantErr: false, // Should show error message but not exit with error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runEddie(t, tt.args...)

			if tt.wantErr {
				assert.True(t, err != nil || stderr != "", "Expected error but got none")
			} else {
				// Commands should handle missing args gracefully
				assert.Contains(t, stdout, "Error:")
			}
		})
	}
}