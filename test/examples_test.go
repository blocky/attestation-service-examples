package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestErrorHandlingAttestFnCall(t *testing.T) {
	projectName := "error_handling_attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		InputFile("successFunc.json").
		InputFile("errorFunc.json").
		InputFile("panicFunc.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestErrorHandlingOnChain(t *testing.T) {
	projectName := "error_handling_on_chain"
	projectDir := filepath.Join(examplesDir, projectName)
	NewHardhatTest(t, projectDir).
		NPMInstall().
		Run("--grep", "Local")
}

func TestErrorHandlingCombined(t *testing.T) {
	errorHandlingOnChainDir := filepath.Join(examplesDir, "error_handling_on_chain")
	onChainSuccessInputFile, err := filepath.Abs(filepath.Join(
		errorHandlingOnChainDir,
		"tmp/out-success.json",
	))
	require.NoError(t, err)
	onChainErrorInputFile, err := filepath.Abs(filepath.Join(
		errorHandlingOnChainDir,
		"tmp/out-error.json",
	))
	require.NoError(t, err)

	errorHandlingName := "error_handling_attest_fn_call"
	t.Run(errorHandlingName, func(t *testing.T) {
		errorHandlingDir := filepath.Join(examplesDir, errorHandlingName)
		NewTestscriptTest(t, errorHandlingDir).
			ExecuteMakeTarget("build").
			InputFile("tmp/x.wasm").
			InputFile("config.toml").
			InputFile("successFunc.json").
			InputFile("errorFunc.json").
			InputFile("panicFunc.json").
			OutputFile("out-success.json", onChainSuccessInputFile).
			OutputFile("out-error.json", onChainErrorInputFile).
			Run(filepath.Join(scriptDir, errorHandlingName+".txtar"))
	})

	require.FileExists(t, onChainSuccessInputFile)
	require.FileExists(t, onChainErrorInputFile)
	t.Setenv("TA_SUCCESS_FILE", onChainSuccessInputFile)
	t.Setenv("TA_ERROR_FILE", onChainErrorInputFile)
	NewHardhatTest(t, errorHandlingOnChainDir).
		NPMInstall().
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
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestESportsDataFromRimble(t *testing.T) {
	projectName := "esports_data_from_rimble"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"YOUR_RIMBLE_MATCH_DATE",
		"YOUR_RIMBLE_MATCH_ID",
		"YOUR_RIMBLE_API_KEY",
	}

	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup(
			"match-winner.json.template",
			requiredEnvVars,
		).
		RenderTemplateFileFromEnvWithCleanup(
			"team-kill-diff.json.template",
			requiredEnvVars,
		).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestAttestFnCall(t *testing.T) {
	projectName := "attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("main.wasm").
		InputFile("main.wasm").
		InputFile("config.toml").
		InputFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestHelloWorldOnChain(t *testing.T) {
	projectName := "hello_world_on_chain"
	projectDir := filepath.Join(examplesDir, projectName)
	NewHardhatTest(t, projectDir).
		NPMInstall().
		Run("--grep", "Local")
}

func TestAttestFnCallCombined(t *testing.T) {
	helloWorldOnChainDir := filepath.Join(examplesDir, "hello_world_on_chain")
	onChainInputFile, err := filepath.Abs(filepath.Join(
		helloWorldOnChainDir,
		"tmp/attest-fn-call-out.json",
	))
	require.NoError(t, err)

	attestFnCallName := "attest_fn_call"
	t.Run(attestFnCallName, func(t *testing.T) {
		attestFnCallDir := filepath.Join(examplesDir, attestFnCallName)
		NewTestscriptTest(t, attestFnCallDir).
			ExecuteMakeTarget("main.wasm").
			InputFile("main.wasm").
			InputFile("config.toml").
			InputFile("fn-call.json").
			OutputFile("out.json", onChainInputFile).
			Run(filepath.Join(scriptDir, attestFnCallName+".txtar"))
	})

	require.FileExists(t, onChainInputFile)
	t.Setenv("TA_FILE", onChainInputFile)
	NewHardhatTest(t, helloWorldOnChainDir).
		NPMInstall().
		Run("--grep", "Local")
}

func TestParamsAndSecrets(t *testing.T) {
	projectName := "params_and_secrets"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		InputFile("fn-call.json").
		InputFile("fn-call-error.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestRandom(t *testing.T) {
	projectName := "random"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		InputFile("fn-call.json").
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
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestTime(t *testing.T) {
	projectName := "time"
	projectDir := filepath.Join(examplesDir, projectName)
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		InputFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestTWAPAttestFnCall(t *testing.T) {
	projectName := "time_weighted_average_price_attest_fn_call"
	projectDir := filepath.Join(examplesDir, "time_weighted_average_price", "attest_fn_call")
	requiredEnvVars := []string{
		"YOUR_COINGECKO_API_KEY",
	}
	NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		InputFile("tmp/x.wasm").
		InputFile("config.toml").
		InputFile("twap-call.json.template").
		RenderTemplateFileFromEnvWithCleanup("iteration-call.json.template", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestTWAPOnChain(t *testing.T) {
	projectName := "time_weighted_average_price/on_chain"
	projectDir := filepath.Join(examplesDir, projectName)
	NewHardhatTest(t, projectDir).
		NPMInstall().
		Run("--grep", "Local")
}

func TestTWAPCombined(t *testing.T) {
	twapOnChainDir := filepath.Join(
		examplesDir,
		"time_weighted_average_price/on_chain",
	)
	onChainInputFile, err := filepath.Abs(filepath.Join(
		twapOnChainDir,
		"tmp/twap.json",
	))
	require.NoError(t, err)

	twapAttestFnCallName := "time_weighted_average_price/attest_fn_call"
	t.Run(twapAttestFnCallName, func(t *testing.T) {
		twapAttestFnCallDir := filepath.Join(examplesDir, twapAttestFnCallName)
		requiredEnvVars := []string{
			"YOUR_COINGECKO_API_KEY",
		}
		NewTestscriptTest(t, twapAttestFnCallDir).
			ExecuteMakeTarget("build").
			InputFile("tmp/x.wasm").
			InputFile("config.toml").
			InputFile("twap-call.json.template").
			RenderTemplateFileFromEnvWithCleanup(
				"iteration-call.json.template",
				requiredEnvVars,
			).
			OutputFile("tmp/twap.json", onChainInputFile).
			Run(filepath.Join(
				scriptDir,
				"time_weighted_average_price_attest_fn_call.txtar",
			))
	})

	require.FileExists(t, onChainInputFile)
	t.Setenv("TA_FILE", onChainInputFile)
	NewHardhatTest(t, twapOnChainDir).
		NPMInstall().
		Run("--grep", "Local")
}
