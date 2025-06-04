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

	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestErrorHandlingAttestFnCall(t *testing.T) {
	projectName := "error_handling_attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("successFunc.json").
		CopyFile("errorFunc.json").
		CopyFile("panicFunc.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestErrorHandlingOnChain(t *testing.T) {
	projectName := "error_handling_on_chain"
	projectDir := filepath.Join(examplesDir, projectName)

	expectedOutput1 := `Success: true
	Error: 
	Value: {"number":42}`
	expectedOutput2 := `Success: false
	Error: expected error
	Value: null`

	NewHardhatTest(t, projectDir).
		NPMInstall().
		OutputContains(expectedOutput1).
		OutputContains(expectedOutput2).
		OutputContains("2 passing").
		NoError().
		Run("--grep", "Local")
}

func TestESportsDataFromPandaScore(t *testing.T) {
	projectName := "esports_data_from_pandascore"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_PANDASCORE_API_ENDPOINT",
		"YOUR_PANDASCORE_API_KEY",
	}

	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

//func TestESportsDataFromRimble(t *testing.T) {
//	projectName := "esports_data_from_rimble"
//	projectDir := filepath.Join(examplesDir, projectName)
//	requiredEnvVars := []string{
//		"YOUR_RIMBLE_API_KEY",
//	}
//
//	NewTestscriptTest(t, projectDir).
//		ExecuteMakeTarget("build").
//		CopyFile("tmp/x.wasm").
//		CopyFile("config.toml").
//		RenderTemplateFileFromEnvWithCleanup(
//			"match-winner.json.template",
//			requiredEnvVars,
//		).
//		RenderTemplateFileFromEnvWithCleanup(
//			"team-kill-diff.json.template",
//			requiredEnvVars,
//		).
//		Run(filepath.Join(scriptDir, projectName+".txtar"))
//}

func TestHelloWorldAttestFnCall(t *testing.T) {
	projectName := "hello_world_attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestHelloWorldOnChain(t *testing.T) {
	projectName := "hello_world_on_chain"
	projectDir := filepath.Join(examplesDir, projectName)

	expectedOutput := `Verified attest-fn-call claims:
	Function: helloWorld
	Hash of code: 083a9a11fa7d1c1ffa224018aca0b1ee1e77c5c8aa007b413e5dae3d3b075a22151b1e1cf318eec8a73de8aea3066478324df942fc2cd1b76cf42e807240115c
	Hash of input: a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26
	Hash of secrets: 9375447cd5307bf7473b8200f039b60a3be491282f852df9f42ce31a8a43f6f8e916c4f8264e7d233add48746a40166eec588be8b7b9b16a5eb698d4c3b06e00
	Output,: Hello, World!`

	NewHardhatTest(t, projectDir).
		NPMInstall().
		OutputContains(expectedOutput).
		OutputContains("1 passing").
		NoError().
		Run("--grep", "Local")
}

func TestParamsAndSecrets(t *testing.T) {
	projectName := "params_and_secrets"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		CopyFile("fn-call-error.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestRandom(t *testing.T) {
	projectName := "random"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestShipmentTrackingWithDHL(t *testing.T) {
	projectName := "shipment_tracking_with_dhl"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_DHL_API_KEY",
	}

	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestSportsDataFromSportRadar(t *testing.T) {
	projectName := "sports_data_from_sportradar"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_SPORTRADAR_API_KEY",
	}

	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestTime(t *testing.T) {
	projectName := "time"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		CopyFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestTWAPOnChain(t *testing.T) {
	projectName := "time_weighted_average_price/on_chain"
	projectDir := filepath.Join(examplesDir, projectName)

	expectedOutput := `Verify attested TWAP in User contract`

	NewHardhatTest(t, projectDir).
		NPMInstall().
		OutputContains(expectedOutput).
		OutputContains("1 passing").
		NoError().
		Run("--grep", "Local")
}
