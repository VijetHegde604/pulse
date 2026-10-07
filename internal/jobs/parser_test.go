package jobs

import (
	"reflect"
	"testing"
)

func TestSplitCommandLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace only",
			input: "   \t\n  ",
			want:  nil,
		},
		{
			name:  "simple command with arguments",
			input: "echo Hello World",
			want:  []string{"echo", "Hello", "World"},
		},
		{
			name:  "command with multiple spaces",
			input: "ls   -la   /tmp",
			want:  []string{"ls", "-la", "/tmp"},
		},
		{
			name:  "double quotes preserving spaces",
			input: `echo "Hello World"`,
			want:  []string{"echo", "Hello World"},
		},
		{
			name:  "single quotes preserving spaces",
			input: `echo 'Hello World'`,
			want:  []string{"echo", "Hello World"},
		},
		{
			name:  "escaped double quotes inside double quotes",
			input: `echo "Hello \"World\""`,
			want:  []string{"echo", `Hello "World"`},
		},
		{
			name:  "complex git commit command",
			input: `git commit -m "feat: initial commit"`,
			want:  []string{"git", "commit", "-m", "feat: initial commit"},
		},
		{
			name:  "sh -c with nested quotes",
			input: `sh -c 'echo "hello"'`,
			want:  []string{"sh", "-c", `echo "hello"`},
		},
		{
			name:  "empty quoted argument",
			input: `echo ""`,
			want:  []string{"echo", ""},
		},
		{
			name:  "escaped space outside quotes",
			input: `echo Hello\ World`,
			want:  []string{"echo", "Hello World"},
		},
		{
			name:    "unclosed double quote",
			input:   `echo "Hello World`,
			wantErr: true,
		},
		{
			name:    "unclosed single quote",
			input:   `echo 'Hello World`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitCommandLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SplitCommandLine(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitCommandLine(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}
