package test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

type HardhatTest struct {
	t          *testing.T
	setupFuncs []func() error
	projectDir string
}

func NewHardhatTest(t *testing.T, projectDir string) *HardhatTest {
	return &HardhatTest{
		t:          t,
		projectDir: projectDir,
	}
}

func (h *HardhatTest) NPMInstall() *HardhatTest {
	setupFunc := func() error {
		cmd := exec.Command("npm", "install")
		cmd.Dir = h.projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		h.t.Logf("Running command: '%s'", cmd.String())
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install npm deps: %v", err)
		}
		return nil
	}
	h.setupFuncs = append(h.setupFuncs, setupFunc)
	return h
}

func (h *HardhatTest) Run(args ...string) {
	for _, setupFunc := range h.setupFuncs {
		if err := setupFunc(); err != nil {
			h.t.Fatalf("failed to setup test: %v", err)
		}
	}

	baseArgs := []string{"hardhat", "test"}
	args = append(baseArgs, args...)
	cmd := exec.Command("npx", args...)

	cmd.Dir = h.projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	h.t.Logf("Running command: '%s'", cmd.String())
	if err := cmd.Run(); err != nil {
		h.t.Fatalf("failed to run hardhat test: %v", err)
	}
}
