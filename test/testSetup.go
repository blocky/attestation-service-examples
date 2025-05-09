package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/rogpeppe/go-internal/testscript"
)

type testSetup struct {
	t          *testing.T
	srcDir     string
	setupFuncs []func(t *testing.T, workDir string) error // Change to return errors
}

func newTestSetup(t *testing.T, srcDir string) *testSetup {
	return &testSetup{t: t, srcDir: srcDir}
}

func (ts *testSetup) SetupFunc() func(env *testscript.Env) error {
	return func(env *testscript.Env) error {
		for _, fn := range ts.setupFuncs {
			if err := fn(ts.t, env.WorkDir); err != nil {
				return err // Propagate the error if one occurs
			}
		}
		return nil
	}
}

func (ts *testSetup) RunMake(target string) *testSetup {
	ts.setupFuncs = append(ts.setupFuncs, func(t *testing.T, _ string) error {
		cmd := exec.Command("make", target)
		cmd.Dir = ts.srcDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run make %s: %v", target, err)
		}
		return nil
	})
	return ts
}

func (ts *testSetup) CopyFile(srcRelPath, dstRelPath string) *testSetup {
	ts.setupFuncs = append(ts.setupFuncs, func(t *testing.T, workDir string) error {
		src := filepath.Join(ts.srcDir, srcRelPath)
		dst := filepath.Join(workDir, dstRelPath)

		content, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("failed to read %s: %v", src, err)
		}

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return fmt.Errorf("failed to create destination directory for %s: %v", dst, err)
		}

		if err := os.WriteFile(dst, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %v", dst, err)
		}
		return nil
	})
	return ts
}

func (ts *testSetup) RenderFile(relPath string, envKeys []string, cleanup bool) *testSetup {
	ts.setupFuncs = append(ts.setupFuncs, func(t *testing.T, workDir string) error {
		src := filepath.Join(ts.srcDir, relPath)
		dst := filepath.Join(workDir, filepath.Base(relPath))

		// Prepare environment variables as a map
		envMap := make(map[string]string, len(envKeys))
		for _, key := range envKeys {
			val := os.Getenv(key)
			if val == "" {
				return fmt.Errorf("environment variable %s is not set", key)
			}
			envMap[key] = val
		}

		tmpl, err := template.New(filepath.Base(src)).Option("missingkey=error").ParseFiles(src)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", src, err)
		}

		// Create a buffer to store the rendered content
		var buf bytes.Buffer

		// Execute the template with the environment map
		err = tmpl.Execute(&buf, envMap)
		if err != nil {
			return fmt.Errorf("failed to render template %s: %v", src, err)
		}

		// Write the rendered content to the destination file
		if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write rendered file %s: %v", dst, err)
		}

		// Register cleanup if needed
		if cleanup {
			t.Cleanup(func() {
				if _, err := os.Stat(dst); err == nil {
					if err := os.Remove(dst); err != nil {
						t.Logf("cleanup: failed to remove %s: %v", dst, err)
					}
				}
			})
		}
		return nil
	})
	return ts
}
