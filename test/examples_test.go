package test

import (
	"path/filepath"
	"testing"
)

const examplesDir = ".."
const scriptDir = "scripts"

func TestCoinPricesFromCoingecko(t *testing.T) {
	projectName := "coin_prices_from_coingecko"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_COINGECKO_API_KEY",
	}

	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestHelloWorldAttestFnCall(t *testing.T) {
	projectName := "hello_world_attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}
