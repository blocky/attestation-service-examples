package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cbroglie/mustache"
	"github.com/rogpeppe/go-internal/testscript"
)

type testSetup struct {
	t          *testing.T
	srcDir     string
	setupFuncs []func(t *testing.T, workDir string)
}

func newTestSetup(t *testing.T, srcDir string) *testSetup {
	return &testSetup{t: t, srcDir: srcDir}
}

func (ts *testSetup) SetupFunc() func(env *testscript.Env) error {
	return func(env *testscript.Env) error {
		for _, fn := range ts.setupFuncs {
			fn(ts.t, env.WorkDir)
		}
		return nil
	}
}

func (ts *testSetup) RunMake(target string) *testSetup {
	ts.setupFuncs = append(ts.setupFuncs, func(t *testing.T, _ string) {
		cmd := exec.Command("make", target)
		cmd.Dir = ts.srcDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to run make %s: %v", target, err)
		}
	})
	return ts
}

func (ts *testSetup) CopyFile(srcRelPath, dstRelPath string) *testSetup {
	ts.setupFuncs = append(ts.setupFuncs, func(t *testing.T, workDir string) {
		src := filepath.Join(ts.srcDir, srcRelPath)
		dst := filepath.Join(workDir, dstRelPath)

		content, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("failed to read %s: %v", src, err)
		}

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatalf("failed to create destination directory for %s: %v", dst, err)
		}

		if err := os.WriteFile(dst, content, 0644); err != nil {
			t.Fatalf("failed to write %s: %v", dst, err)
		}
	})
	return ts
}

func (ts *testSetup) RenderFile(relPath string, envKeys []string, cleanup bool) *testSetup {
	ts.setupFuncs = append(ts.setupFuncs, func(t *testing.T, workDir string) {
		src := filepath.Join(ts.srcDir, relPath)
		dst := filepath.Join(workDir, filepath.Base(relPath))

		// Prepare environment variables as a map
		envMap := make(map[string]string, len(envKeys))
		for _, key := range envKeys {
			val := os.Getenv(key)
			if val == "" {
				t.Fatalf("environment variable %s is not set", key)
			}
			envMap[key] = val
		}

		// Render the file using mustache
		rendered, err := mustache.RenderFile(src, envMap)
		if err != nil {
			t.Fatalf("failed to render template %s: %v", src, err)
		}

		// Write the rendered result to destination
		if err := os.WriteFile(dst, []byte(rendered), 0644); err != nil {
			t.Fatalf("failed to write rendered file %s: %v", dst, err)
		}

		if cleanup {
			t.Cleanup(func() {
				if err := os.Remove(dst); err != nil {
					t.Logf("cleanup: failed to remove %s: %v", dst, err)
				}
			})
		}
	})
	return ts
}
