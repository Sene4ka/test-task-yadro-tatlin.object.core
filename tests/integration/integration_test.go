package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	root := findProjectRoot(t)
	binPath := filepath.Join(t.TempDir(), "namecount")
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", binPath, "./cmd/namecount")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build ./cmd: %v\n%s", err, out)
	}
	return binPath
}

func TestIntegration(t *testing.T) {
	binary := buildBinary(t)

	tests := []struct {
		name     string
		args     []string
		input    string
		wantOut  string
		wantFile string
		wantErr  string
		wantCode int
	}{
		{
			name:    "default freq desc",
			args:    []string{},
			input:   "Алёна\nМиша\nАлёна\nДима\n",
			wantOut: "Алёна:2\nДима:1\nМиша:1\n",
		},
		{
			name:    "alph asc",
			args:    []string{"-sort-by", "alph"},
			input:   "Алёна\nМиша\nАлёна\nДима\n",
			wantOut: "Алёна:2\nДима:1\nМиша:1\n",
		},
		{
			name:    "freq asc",
			args:    []string{"-sort-by", "freq", "-order", "asc"},
			input:   "Алёна\nМиша\nАлёна\nДима\n",
			wantOut: "Дима:1\nМиша:1\nАлёна:2\n",
		},
		{
			name:    "ignore case",
			input:   "алёНа\nАлёна\n",
			wantOut: "Алёна:2\n",
		},
		{
			name:    "ignore case non-russian",
			input:   "алёНа\nАлёна\njoHn\n",
			wantOut: "Алёна:2\nJohn:1\n",
		},
		{
			name:    "preserve case",
			args:    []string{"-preserve-case"},
			input:   "алёна\nАлёна\n",
			wantOut: "Алёна:1\nалёна:1\n",
		},
		{
			name:     "missing file",
			args:     []string{},
			wantErr:  "Error while opening file",
			wantCode: 1,
		},
		{
			name:     "output to file",
			args:     []string{"-o", "result.txt"},
			input:    "Алёна\nМиша\nАлёна\nДима\n",
			wantFile: "Алёна:2\nДима:1\nМиша:1\n",
		},
		{
			name:     "invalid sort type",
			args:     []string{"-sort-by", "bad"},
			wantErr:  "Unknown sorting type",
			wantCode: 1,
		},
		{
			name:     "invalid order",
			args:     []string{"-order", "bad"},
			wantErr:  "Invalid order",
			wantCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			var inputFilePath string
			if tt.input != "" || (tt.args != nil && tt.wantOut != "" && tt.name != "missing file") {
				inputFilePath = filepath.Join(tmpDir, "names.txt")
				if tt.input != "" {
					if err := os.WriteFile(inputFilePath, []byte(tt.input), 0644); err != nil {
						t.Fatal(err)
					}
				} else {
					if err := os.WriteFile(inputFilePath, []byte{}, 0644); err != nil {
						t.Fatal(err)
					}
				}
			}

			var args []string
			if tt.wantFile != "" {
				outputPath := filepath.Join(tmpDir, "result.txt")
				args = append([]string{}, tt.args...)
				for i, a := range args {
					if a == "-o" && i+1 < len(args) {
						args[i+1] = outputPath
						break
					}
				}
				args = append(args, inputFilePath)
			} else {
				args = append(tt.args, inputFilePath)
				if tt.name == "missing file" {
					args = append(tt.args, "nonexistent_file.txt")
				}
			}

			cmd := exec.Command(binary, args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()

			var exitCode int
			if err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					t.Fatalf("failed to run: %v", err)
				}
			}
			if exitCode != tt.wantCode {
				t.Errorf("exit code = %d, want %d", exitCode, tt.wantCode)
			}

			if tt.wantOut != "" && !strings.Contains(stdout.String(), tt.wantOut) {
				t.Errorf("stdout = %q, want to contain %q", stdout.String(), tt.wantOut)
			}

			if tt.wantErr != "" && !strings.Contains(stderr.String(), tt.wantErr) {
				t.Errorf("stderr = %q, want to contain %q", stderr.String(), tt.wantErr)
			}

			if tt.wantFile != "" {
				outputPath := filepath.Join(tmpDir, "result.txt")
				data, err := os.ReadFile(outputPath)
				if err != nil {
					t.Fatalf("failed to read output file: %v", err)
				}
				if string(data) != tt.wantFile {
					t.Errorf("output file = %q, want %q", string(data), tt.wantFile)
				}
				if stdout.Len() != 0 {
					t.Errorf("stdout should be empty when -o is used, got: %q", stdout.String())
				}
			}
		})
	}
}
