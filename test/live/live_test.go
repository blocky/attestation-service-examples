package live_test

import (
	"path/filepath"
	"testing"

	"attestation-service-examples.test"
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
		CopyFile("tmp/x.wasm").
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
		CopyFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyFile("successFunc.json").
		CopyFile("errorFunc.json").
		CopyFile("panicFunc.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
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
		CopyFile("tmp/x.wasm").
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
		CopyFile("tmp/x.wasm").
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
		CopyFile("main.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
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
		CopyFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyFile("fn-call.json").
		CopyFile("fn-call-error.json").
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
		CopyFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyFile("fn-call.json").
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
		CopyFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		RenderTemplateFileFromEnvWithCleanup("fn-call.json", requiredEnvVars).
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}

func TestLiveTWAPAttestFnCall(t *testing.T) {
	projectName := "time_weighted_average_price_attest_fn_call"
	projectDir := filepath.Join(examplesDir, "time_weighted_average_price", "attest_fn_call")
	requiredEnvVars := []string{
		"LIVE_TEST_PLATFORM",
		"LIVE_TEST_CODE",
		"LIVE_TEST_AUTH_TOKEN",
		"LIVE_TEST_HOST",
		"YOUR_COINGECKO_API_KEY",
	}
	test.NewTestscriptTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("twap-call.json.template").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		RenderTemplateFileFromEnvWithCleanup("iteration-call.json.template", requiredEnvVars).
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
		CopyFile("tmp/x.wasm").
		RenderTemplateStringFromEnvWithCleanup(
			liveTestConfigTemplate,
			"config.toml",
			requiredEnvVars).
		CopyFile("fn-call.json").
		Run(filepath.Join(scriptDir, projectName+".txtar"))
}
