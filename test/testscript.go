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

func NewTestscriptTest(
	t *testing.T,
	projectDir string,
) *TestscriptTest {
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

func copyFile(src string, dst string) error {
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

// CopyProjectFile copies a file from the project directory to the test script's
// working directory, preserving the relative path.
func (e *TestscriptTest) CopyProjectFile(srcRelPath string) *TestscriptTest {
	setupFunc := func(env *testscript.Env) error {
		src := filepath.Join(e.projectDir, srcRelPath)
		dst := filepath.Join(env.WorkDir, srcRelPath)
		return copyFile(src, dst)
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

// CopyTestscriptFile copies a file from the test script's working directory
// after script execution to an absolute destination path.
func (e *TestscriptTest) CopyTestscriptFile(
	srcRelPath string,
	dstAbsPath string,
) *TestscriptTest {
	setupFunc := func(env *testscript.Env) error {
		env.Defer(func() {
			src := filepath.Join(env.WorkDir, srcRelPath)
			if !filepath.IsAbs(dstAbsPath) {
				e.t.Fatalf("destination path must be absolute: %s", dstAbsPath)
			}
			if err := copyFile(src, dstAbsPath); err != nil {
				e.t.Fatalf(
					"failed to copy file from '%s' to '%s': %v",
					src,
					dstAbsPath,
					err,
				)
			}
		})
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *TestscriptTest) RenderTemplateFileFromEnvWithCleanup(
	srcRelPath string,
	envKeys []string,
) *TestscriptTest {
	src := filepath.Join(e.projectDir, srcRelPath)
	tmpl := template.Must(template.ParseFiles(src))
	dstRelPath := srcRelPath
	return e.RenderTemplateFromEnvWithCleanup(tmpl, dstRelPath, envKeys)
}

func (e *TestscriptTest) RenderTemplateStringFromEnvWithCleanup(
	templateString string,
	dstRelPath string,
	envKeys []string,
) *TestscriptTest {
	tmpl := template.Must(
		template.New(filepath.Base(dstRelPath)).
			Parse(templateString),
	)
	return e.RenderTemplateFromEnvWithCleanup(tmpl, dstRelPath, envKeys)
}

func (e *TestscriptTest) RenderTemplateFromEnvWithCleanup(
	tmpl *template.Template,
	dstRelPath string,
	envKeys []string,
) *TestscriptTest {
	setupFunc := func(env *testscript.Env) error {
		dst := filepath.Join(env.WorkDir, dstRelPath)

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

		dstFile, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", dst, err)
		}
		defer dstFile.Close()

		tmpl.Option("missingkey=error")
		if err := tmpl.Execute(dstFile, envMap); err != nil {
			return fmt.Errorf("failed to execute template %s: %v", tmpl.Name(), err)
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
		"setEnvValueFromFile": setEnvValueFromFileCmd("setEnvValueFromFile"),
	}
	testscript.Run(e.t, e.params)
}

// setEnvValueFromFile reads the contents of a file (filePath) and sets an
// environment variable (envKey) to the rendered contents of the file
func setEnvValueFromFileCmd(
	cmdName string,
) func(*testscript.TestScript, bool, []string) {
	return func(ts *testscript.TestScript, neg bool, args []string) {
		workDir := ts.Getenv("WORK")
		switch {
		case neg:
			ts.Fatalf("unsupported: ! %s", cmdName)
		case len(args) != 2:
			ts.Fatalf("usage: %s envKey filePath", cmdName)
		case workDir == "":
			ts.Fatalf("WORK environment variable is not set")
		}

		src := filepath.Join(workDir, args[1])
		content, err := os.ReadFile(src)
		if err != nil {
			ts.Fatalf("failed to read file '%s': %v", src, err)
		}
		contentStr := strings.TrimSpace(string(content))
		if contentStr == "" {
			ts.Fatalf("file '%s' is empty", src)
		}
		ts.Setenv(args[0], contentStr)
	}
}
