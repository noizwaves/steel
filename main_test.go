package main

import (
	"bytes"
	"io"
	"os/exec"
	"testing"

	"github.com/creack/pty"
	"github.com/stretchr/testify/assert"
)

const steelFilename = "./steel_testable"

func TestProjects(t *testing.T) {
	buildSteel(t)

	tcs := []struct {
		name             string
		dir              string
		command          string
		expectedContains string
	}{
		{
			name:             "Ruby",
			dir:              "testdata/project_ruby",
			command:          "ruby version.rb",
			expectedContains: "RubyVersion==3.2.1",
		},
		{
			name:             "Go",
			dir:              "testdata/project_go",
			command:          "go run .",
			expectedContains: "GoVersion==go1.21.0",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			runSteel(t, "--workdir", tc.dir, "setup")

			output := runSteelWithInput(t, tc.command, "--workdir", tc.dir, "shell")

			assert.Contains(t, output, tc.expectedContains)
			assert.Contains(t, output, "🤘>")
		})
	}
}

func buildSteel(t *testing.T) {
	t.Helper()

	cmd := exec.Command("go", "build", "-o", steelFilename, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build steel: %s", err)
	}
}

func runSteel(t *testing.T, command ...string) string {
	t.Helper()

	cmd := exec.Command(steelFilename, command...)

	out := bytes.Buffer{}
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatal("command exited with an error")
	}

	return out.String()
}

func runSteelWithInput(t *testing.T, input string, command ...string) string {
	t.Helper()

	cmd := exec.Command(steelFilename, command...)

	tty, err := pty.Start(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tty.Close() }()

	go func() {
		tty.WriteString(input + "\n")
		tty.Write([]byte{4}) // EOT
	}()

	out := bytes.Buffer{}
	io.Copy(&out, tty)

	err = cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}

	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatal("command exited with an error")
	}

	return out.String()
}
