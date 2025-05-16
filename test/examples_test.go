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

func TestErrorHandlingAttestFnCall(t *testing.T) {
	projectName := "error_handling_attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		CopyFile("fn-call-error.json").
		CopyFile("fn-call-panic.json").
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestESportsDataFromPandaScore(t *testing.T) {
	projectName := "esports_data_from_pandascore"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_PANDASCORE_API_KEY",
	}

	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestESportsDataFromRimble(t *testing.T) {
	projectName := "esports_data_from_rimble"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_RIMBLE_API_KEY",
	}

	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup(
			"fn-call-match-winner.json",
			requiredEnvVars,
		).
		RenderTemplateFileFromEnvWithCleanup(
			"fn-call-team-kill.json",
			requiredEnvVars,
		).
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

func TestParamsAndSecrets(t *testing.T) {
	projectName := "params_and_secrets"
	projectDir := filepath.Join(examplesDir, projectName)
	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		CopyFile("fn-call-error.json").
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestRandom(t *testing.T) {
	projectName := "random"
	projectDir := filepath.Join(examplesDir, projectName)
	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestShipmentTrackingWithDHL(t *testing.T) {
	projectName := "shipment_tracking_with_dhl"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_DHL_API_KEY",
	}

	NewProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		RunScript(filepath.Join(scriptDir, projectName+".txtar"))
}
