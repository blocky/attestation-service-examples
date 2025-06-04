package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HardhatTest struct {
	t          *testing.T
	checkFuncs []func(string, error)
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

func (h *HardhatTest) NoError() *HardhatTest {
	checkFunc := func(_ string, err error) {
		assert.NoError(h.t, err)
	}
	h.checkFuncs = append(h.checkFuncs, checkFunc)
	return h
}

func (h *HardhatTest) OutputContains(expected string) *HardhatTest {
	checkFunc := func(output string, _ error) {
		assert.Contains(h.t, output, expected)
	}
	h.checkFuncs = append(h.checkFuncs, checkFunc)
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

	var outBuff, errBuff bytes.Buffer
	outTee := io.MultiWriter(&outBuff, os.Stdout)
	errTee := io.MultiWriter(&errBuff, os.Stderr)
	cmd.Stdout = outTee
	cmd.Stderr = errTee

	h.t.Logf("Running command: '%s'", cmd.String())
	err := cmd.Run()

	for _, checkFunc := range h.checkFuncs {
		checkFunc(outBuff.String(), err)
	}
}
