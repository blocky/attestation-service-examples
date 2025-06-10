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
