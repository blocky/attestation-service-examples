package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/rogpeppe/go-internal/testscript"
)

type projectTest struct {
	t          *testing.T
	params     testscript.Params
	setupFuncs []func(*testscript.Env) error
	projectDir string
}

func newProjectTest(t *testing.T, projectDir string) *projectTest {
	params := testscript.Params{
		Files:               []string{},
		Setup:               nil,
		RequireExplicitExec: true,
		RequireUniqueNames:  true,
	}
	return &projectTest{
		t:          t,
		params:     params,
		projectDir: projectDir,
	}
}

func (e *projectTest) CopyFile(srcRelPath string) *projectTest {
	setupFunc := func(env *testscript.Env) error {
		src := filepath.Join(e.projectDir, srcRelPath)
		dst := filepath.Join(env.WorkDir, srcRelPath)
		dstDir := filepath.Dir(dst)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			msg := "failed to create destination directory %s: %w"
			return fmt.Errorf(msg, dstDir, err)
		}

		srcContent, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("failed to read source file %s: %w", src, err)
		}

		if err := os.WriteFile(dst, srcContent, 0644); err != nil {
			msg := "failed to write destination file %s: %w"
			return fmt.Errorf(msg, dst, err)
		}
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *projectTest) RenderTemplateFromEnvWithCleanup(
	srcRelPath string,
	envKeys []string,
) *projectTest {
	setupFunc := func(env *testscript.Env) error {
		src := filepath.Join(e.projectDir, srcRelPath)
		dst := filepath.Join(env.WorkDir, srcRelPath)

		dstDir := filepath.Dir(dst)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			msg := "failed to create destination directory %s: %w"
			return fmt.Errorf(msg, dstDir, err)
		}

		envMap := make(map[string]string, len(envKeys))
		for _, key := range envKeys {
			val := os.Getenv(key)
			if val == "" {
				return fmt.Errorf("environment variable %s is not set", key)
			}
			envMap[key] = val
		}

		tmpl, err := template.
			New(filepath.Base(src)).
			Option("missingkey=error").
			ParseFiles(src)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", src, err)
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", dst, err)
		}
		defer dstFile.Close()

		if err := tmpl.Execute(dstFile, envMap); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", src, err)
		}

		// rendered files may contain secrets, register a cleanup func for removal
		env.Defer(func() {
			if err := os.Remove(dst); err != nil {
				e.t.Fatalf("failed to remove file %s: %v", dst, err)
			}
		})
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *projectTest) ExecuteMakeTarget(target string) *projectTest {
	setupFunc := func(env *testscript.Env) error {
		cmd := exec.Command("make", target)
		cmd.Dir = e.projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run make %s: %v", target, err)
		}
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *projectTest) RunScript(scriptFile string) {
	e.params.Setup = func(env *testscript.Env) error {
		for _, setupFunc := range e.setupFuncs {
			if err := setupFunc(env); err != nil {
				return err
			}
		}
		return nil
	}
	e.params.Files = []string{scriptFile}
	testscript.Run(e.t, e.params)
}
