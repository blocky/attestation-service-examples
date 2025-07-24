package live_test

import (
	"path/filepath"
	"testing"

	"attestation-service-examples.test"
	"github.com/stretchr/testify/require"
)

const liveTestConfigTemplate = `
acceptable_measurements = [
{ platform = "{{.LIVE_TEST_PLATFORM }}", code = "{{.LIVE_TEST_CODE}}" },
]
auth_token = "{{.LIVE_TEST_AUTH_TOKEN}}"
host = "{{.LIVE_TEST_HOST}}"
`

const examplesDir = "../.."
const scriptDir = "../scripts"

func TestLiveCoinPricesFromCoingecko(t *testing.T) {
	projectName := "coin_prices_from_coingecko"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
		"YOUR_COINGECKO_API_KEY",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveErrorHandlingAttestFnCall(t *testing.T) {
	projectName := "error_handling_attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyProjectFile("successFunc.json").
		CopyProjectFile("errorFunc.json").
		CopyProjectFile("panicFunc.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveErrorHandlingCombined(t *testing.T) {
	errorHandlingOnChainDir := filepath.Join(examplesDir, "error_handling_on_chain")
	onChainSuccessCopyProjectFile, err := filepath.Abs(filepath.Join(
		errorHandlingOnChainDir,
		"tmp/out-success.json",
	))
	require.NoError(t, err)
	onChainErrorCopyProjectFile, err := filepath.Abs(filepath.Join(
		errorHandlingOnChainDir,
		"tmp/out-error.json",
	))
	require.NoError(t, err)

	errorHandlingName := "error_handling_attest_fn_call"
	t.Run(errorHandlingName, func(t *testing.T) {
		errorHandlingDir := filepath.Join(examplesDir, errorHandlingName)
		requiredEnvVars := []string{
			"LIVE_TEST_PLATFORM",
			"LIVE_TEST_CODE",
			"LIVE_TEST_AUTH_TOKEN",
			"LIVE_TEST_HOST",
		}

		test.NewTestscriptTest(t, errorHandlingDir).
			ExecuteMakeTarget("build").
			CopyProjectFile("tmp/x.wasm").
			RenderTemplateStringFromEnvWithCleanup(
				liveTestConfigTemplate,
				"config.toml",
				requiredEnvVars).
			CopyProjectFile("successFunc.json").
			CopyProjectFile("errorFunc.json").
			CopyProjectFile("panicFunc.json").
			CopyTestscriptFile("out-success.json", onChainSuccessCopyProjectFile).
			CopyTestscriptFile("out-error.json", onChainErrorCopyProjectFile).
			Run(filepath.Join(scriptDir, errorHandlingName+".txtar"))
	})

	require.FileExists(t, onChainSuccessCopyProjectFile)
	require.FileExists(t, onChainErrorCopyProjectFile)
	t.Setenv("TA_SUCCESS_FILE", onChainSuccessCopyProjectFile)
	t.Setenv("TA_ERROR_FILE", onChainErrorCopyProjectFile)
	test.NewHardhatTest(t, errorHandlingOnChainDir).
		NPMInstall().
		Run("--grep", "Local")
}

func TestLiveESportsDataFromPandaScore(t *testing.T) {
	projectName := "esports_data_from_pandascore"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
		"YOUR_PANDASCORE_API_ENDPOINT",
		"YOUR_PANDASCORE_API_KEY",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveESportsDataFromRimble(t *testing.T) {
	projectName := "esports_data_from_rimble"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
		"YOUR_RIMBLE_MATCH_DATE",
		"YOUR_RIMBLE_MATCH_ID",
		"YOUR_RIMBLE_API_KEY",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
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

func TestLiveAttestFnCall(t *testing.T) {
	projectName := "attest_fn_call"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("main.wasm").
		CopyProjectFile("main.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyProjectFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveAttestFnCallCombined(t *testing.T) {
	helloWorldOnChainDir := filepath.Join(examplesDir, "hello_world_on_chain")
	onChainCopyProjectFile, err := filepath.Abs(filepath.Join(
		helloWorldOnChainDir,
		"tmp/attest-fn-call-out.json",
	))
	require.NoError(t, err)

	attestFnCallName := "attest_fn_call"
	t.Run(attestFnCallName, func(t *testing.T) {
		attestFnCallDir := filepath.Join(examplesDir, attestFnCallName)
		requiredEnvVars := []string{
			"LIVE_TEST_PLATFORM",
			"LIVE_TEST_CODE",
			"LIVE_TEST_AUTH_TOKEN",
			"LIVE_TEST_HOST",
		}

		test.NewTestscriptTest(t, attestFnCallDir).
			ExecuteMakeTarget("main.wasm").
			CopyProjectFile("main.wasm").
			RenderTemplateStringFromEnvWithCleanup(
				liveTestConfigTemplate,
				"config.toml",
				requiredEnvVars).
			CopyProjectFile("fn-call.json").
			CopyTestscriptFile("out.json", onChainCopyProjectFile).
			Run(filepath.Join(scriptDir, attestFnCallName+".txtar"))
	})

	require.FileExists(t, onChainCopyProjectFile)
	t.Setenv("TA_FILE", onChainCopyProjectFile)
	test.NewHardhatTest(t, helloWorldOnChainDir).
		NPMInstall().
		Run("--grep", "Local")
}

func TestLiveParamsAndSecrets(t *testing.T) {
	projectName := "params_and_secrets"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyProjectFile("fn-call.json").
		CopyProjectFile("fn-call-error.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveRandom(t *testing.T) {
	projectName := "random"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyProjectFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveShipmentTrackingWithDHL(t *testing.T) {
	projectName := "shipment_tracking_with_dhl"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
		"YOUR_DHL_API_KEY",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveTime(t *testing.T) {
	projectName := "time"
	projectDir := filepath.Join(examplesDir, projectName)
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
	}

	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyProjectFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyProjectFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}
