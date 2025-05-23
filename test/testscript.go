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

type ProjectTest struct {
	t          *testing.T
	params     testscript.Params
	setupFuncs []func(*testscript.Env) error
	projectDir string
}

func NewProjectTest(t *testing.T, projectDir string) *ProjectTest {
	params := testscript.Params{
		Files:               []string{},
		Setup:               nil,
		RequireExplicitExec: true,
		RequireUniqueNames:  true,
	}
	return &ProjectTest{
		t:          t,
		params:     params,
		projectDir: projectDir,
	}
}

func copyFile(srcPath string, dstPath string, mode os.FileMode) error {
	srcContent, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("reading source file %s: %w", srcPath, err)
	}

	if err := os.WriteFile(dstPath, srcContent, mode); err != nil {
		return fmt.Errorf("writing destination file %s: %w", dstPath, err)
	}
	return nil
}

func makeDir(dirPath string, mode os.FileMode) error {
	if err := os.MkdirAll(dirPath, mode); err != nil {
		return fmt.Errorf("creating directory %s: %w", dirPath, err)
	}
	return nil
}

func (e *ProjectTest) CopyFile(srcRelPath string) *ProjectTest {
	setupFunc := func(env *testscript.Env) error {
		src := filepath.Join(e.projectDir, srcRelPath)
		dst := filepath.Join(env.WorkDir, srcRelPath)
		dstDir := filepath.Dir(dst)

		if err := makeDir(dstDir, 0755); err != nil {
			return err
		}
		return copyFile(src, dst, 0644)
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *ProjectTest) CopyDir(srcRelPath string) *ProjectTest {
	setupFunc := func(env *testscript.Env) error {
		srcDir := filepath.Join(e.projectDir, srcRelPath)
		dstDir := filepath.Join(env.WorkDir, srcRelPath)

		walkFunc := func(srcPath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("walking '%s': %w", srcPath, err)
			}

			relPath, err := filepath.Rel(srcDir, srcPath)
			if err != nil {
				return fmt.Errorf("creating relative path: %w", err)
			}
			dstPath := filepath.Join(dstDir, relPath)

			if info.IsDir() {
				return makeDir(dstPath, info.Mode())
			} else {
				return copyFile(srcPath, dstPath, info.Mode())
			}
		}

		err := filepath.Walk(srcDir, walkFunc)
		if err != nil {
			return fmt.Errorf("copying directory %s: %w", srcDir, err)
		}
		return nil
	}

	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *ProjectTest) RenderTemplateFileFromEnvWithCleanup(
	srcRelPath string,
	envKeys []string,
) *ProjectTest {
	src := filepath.Join(e.projectDir, srcRelPath)
	tmpl := template.Must(template.ParseFiles(src))
	dstRelPath := srcRelPath
	return e.RenderTemplateFromEnvWithCleanup(tmpl, dstRelPath, envKeys)
}

func (e *ProjectTest) RenderTemplateStringFromEnvWithCleanup(
	templateString string,
	dstRelPath string,
	envKeys []string,
) *ProjectTest {
	tmpl := template.Must(
		template.New(filepath.Base(dstRelPath)).
			Parse(templateString),
	)
	return e.RenderTemplateFromEnvWithCleanup(tmpl, dstRelPath, envKeys)
}

func (e *ProjectTest) RenderTemplateFromEnvWithCleanup(
	tmpl *template.Template,
	dstRelPath string,
	envKeys []string,
) *ProjectTest {
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

func (e *ProjectTest) ExecuteMakeTarget(target string) *ProjectTest {
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

func (e *ProjectTest) NPMInstallDeps() *ProjectTest {
	setupFunc := func(env *testscript.Env) error {
		// Configure both npm and npx to only create and write files to the
		// current working directory instead of the default $HOME directory
		env.Setenv("NPM_CONFIG_CACHE", env.WorkDir+"/.npm")
		env.Setenv("NPM_CONFIG_PREFIX", env.WorkDir+"/.npm-global")
		env.Setenv("XDG_CACHE_HOME", env.WorkDir+"/.npx")
		env.Setenv("XDG_CONFIG_HOME", env.WorkDir+"/.npx")
		env.Setenv("XDG_DATA_HOME", env.WorkDir+"/.npx")

		cmd := exec.Command("npm", "install")
		cmd.Dir = env.WorkDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		e.t.Logf("Running command: '%s'", cmd.String())
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install npm deps: %v", err)
		}
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *ProjectTest) SetEnv(key string, value string) *ProjectTest {
	setupFunc := func(env *testscript.Env) error {
		env.Setenv(key, value)
		return nil
	}
	e.setupFuncs = append(e.setupFuncs, setupFunc)
	return e
}

func (e *ProjectTest) RunScript(scriptFile string) {
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
