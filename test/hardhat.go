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
	prevEnv    map[string]string
}

func NewHardhatTest(t *testing.T, projectDir string) *HardhatTest {
	return &HardhatTest{
		t:          t,
		projectDir: projectDir,
		prevEnv:    make(map[string]string),
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

func (h *HardhatTest) SetEnv(key string, value string) *HardhatTest {
	h.prevEnv[key] = os.Getenv(key)
	err := os.Setenv(key, value)
	if err != nil {
		h.t.Fatalf("failed to set environment variable %s: %v", key, err)
	}
	return h
}

func (h *HardhatTest) Run(args ...string) {
	defer resetEnv(h)
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

func resetEnv(h *HardhatTest) {
	for key, value := range h.prevEnv {
		err := os.Setenv(key, value)
		if err != nil {
			h.t.Logf("failed to reset environment variable %s: %v", key, err)
		}
	}
}