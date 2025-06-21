package e2e

import (
	"os"
	"os/exec"
	"testing"
)

const eddieBinary = "../../eddie"

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build", "-o", eddieBinary, "../../.")
	if err := cmd.Run(); err != nil {
		panic("Failed to build eddie binary: " + err.Error())
	}

	code := m.Run()

	os.Remove(eddieBinary)
	os.Exit(code)
}

func runEddie(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(eddieBinary, args...)

	stdout, err := cmd.Output()
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}

	return string(stdout), stderr, err
}
