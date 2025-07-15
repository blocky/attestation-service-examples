package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/rogpeppe/go-internal/testscript"
)

type TestscriptTest struct {
	t          *testing.T
	params     testscript.Params
	setupFuncs []func(*testscript.Env) error
	projectDir string
}

func NewTestscriptTest(t *testing.T, projectDir string) *TestscriptTest {
	params := testscript.Params{
		Files:               []string{},
		Setup:               nil,
		RequireExplicitExec: true,
		RequireUniqueNames:  true,
	}
	return &TestscriptTest{
		t:          t,
		params:     params,
		projectDir: projectDir,
	}
}

func (e *TestscriptTest) CopyFile(srcRelPath string) *TestscriptTest {
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

func templateFromFile(srcAbsPath, leftDelim, rightDelim string) (*template.Template, error) {
	tmpl, err := template.
		New(filepath.Base(srcAbsPath)).
		Delims(leftDelim, rightDelim).
		ParseFiles(srcAbsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template file '%s': %v", srcAbsPath, err)
	}
	return tmpl, nil
}

func (e *TestscriptTest) RenderTemplateFileFromEnvWithCleanup(
	srcRelPath string,
	envKeys []string,
) *TestscriptTest {
	src := filepath.Join(e.projectDir, srcRelPath)
	tmpl, err := templateFromFile(src, "{{", "}}")
	if err != nil {
		e.t.Fatalf("failed to generate template file '%s': %v", src, err)
	}
	dstRelPath := srcRelPath
	return e.RenderTemplateFromEnvWithCleanup(tmpl, dstRelPath, envKeys)
}

func (e *TestscriptTest) RenderTemplateStringFromEnvWithCleanup(
	templateString string,
	dstRelPath string,
	envKeys []string,
) *TestscriptTest {
	tmpl, err := template.
		New(filepath.Base(dstRelPath)).
		Delims("{{", "}}").
		Parse(templateString)
	if err != nil {
		e.t.Fatalf("failed to parse template string: %v", err)
	}
	return e.RenderTemplateFromEnvWithCleanup(tmpl, dstRelPath, envKeys)
}

func renderTemplateFromEnv(
	tmpl *template.Template,
	dst string,
	envKeys []string,
	getEnvFunc func(string) string,
) error {
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		msg := "failed to create destination directory %s: %w"
		return fmt.Errorf(msg, dstDir, err)
	}

	envMap := make(map[string]string, len(envKeys))
	for _, key := range envKeys {
		val := getEnvFunc(key)
		if val == "" {
			return fmt.Errorf("environment variable '%s' is not set", key)
		}
		envMap[key] = val
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file '%s': %w", dst, err)
	}
	defer dstFile.Close()

	tmpl.Option("missingkey=error")
	if err := tmpl.Execute(os.Stdout, envMap); err != nil {
		return fmt.Errorf("failed to execute template %s: %v", tmpl.Name(), err)
	}
	return nil
}

func (e *TestscriptTest) RenderTemplateFromEnvWithCleanup(
	tmpl *template.Template,
	dstRelPath string,
	envKeys []string,
) *TestscriptTest {
	setupFunc := func(env *testscript.Env) error {
		dst := filepath.Join(env.WorkDir, dstRelPath)

		err := renderTemplateFromEnv(tmpl, dst, envKeys, os.Getenv)
		if err != nil {
			return fmt.Errorf("failed to render template '%s' to '%s': %v", tmpl.Name(), dst, err)
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

func (e *TestscriptTest) ExecuteMakeTarget(target string) *TestscriptTest {
	setupFunc := func(env *testscript.Env) error {
		cmd := exec.Command("make", target)
		cmd.Dir = e.projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		e.t.Logf("Running command: '%s'", cmd.String())
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run 'make %s': %v", target, err)
		}
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *TestscriptTest) Run(scriptFile string) {
	e.params.Setup = func(env *testscript.Env) error {
		for _, setupFunc := range e.setupFuncs {
			if err := setupFunc(env); err != nil {
				return err
			}
		}
		return nil
	}
	e.params.Files = []string{scriptFile}
	e.params.Cmds = map[string]func(*testscript.TestScript, bool, []string){
		"renderTemplateFileFromEnv": cmdRenderTemplateFileFromEnv,
	}
	testscript.Run(e.t, e.params)
}

// renderTemplateFileFromEnv [-cleanup] src dst envKeys
// If -cleanup is specified, the destination file will be removed after the testscript completes.
// src is a template file relative to the testscript work directory.
// dst is the destination file relative to the testscript work directory.
// envKeys is a comma-separated list of environment variable names that will be used to render the template.
func cmdRenderTemplateFileFromEnv(
	ts *testscript.TestScript,
	neg bool,
	args []string,
) {
	if neg {
		ts.Fatalf("unsupported: ! renderTemplateFileFromEnv")
	}
	var cleanup bool
	var srcRelPath, dstRelPath string
	var envKeys []string
	switch len(args) {
	case 3:
		srcRelPath = args[0]
		dstRelPath = args[1]
	case 4:
		if args[0] != "-cleanup" {
			ts.Fatalf("usage: renderTemplateFileFromEnv [-cleanup] src dst")
		}
		cleanup = true
		srcRelPath = args[1]
		dstRelPath = args[2]
		envKeys = strings.Split(args[3], ",")
	default:
		ts.Fatalf("usage: renderTemplateFileFromEnv [-cleanup] src dst")
	}
	workDir := ts.Getenv("WORK")
	if workDir == "" {
		ts.Fatalf("WORK environment variable is not set")
	}
	src := filepath.Join(workDir, srcRelPath)
	dst := filepath.Join(workDir, dstRelPath)
	tmpl, err := templateFromFile(src, "[[", "]]")
	if err != nil {
		ts.Fatalf("failed to generate template file '%s': %v", src, err)
	}

	err = renderTemplateFromEnv(tmpl, dst, envKeys, ts.Getenv)
	if err != nil {
		ts.Fatalf("failed to render template '%s' to '%s': %v", src, dst, err)
	}
	if cleanup {
		ts.Defer(func() {
			if err := os.Remove(dst); err != nil {
				ts.Logf("failed to clean-up file '%s': %v", dst, err)
			}
		})
	}
}
