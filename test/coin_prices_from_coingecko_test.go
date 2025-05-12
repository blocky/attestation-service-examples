package test

import (
	"path/filepath"
	"testing"
)

const examplesDir = ".."

const configTemplate = `
acceptable_measurements = [
{ platform = "{{.LIVE_TEST_PLATFORM}}", code = "{{.LIVE_TEST_CODE}}" },
]
auth_token = "{{.LIVE_TEST_AUTH_TOKEN}}"
host = "{{.LIVE_TEST_HOST}}"
`

func TestCoinPricesFromCoingecko(t *testing.T) {
	t.Run("coingecko example local test", func(t *testing.T) {
		projectName := "coin_prices_from_coingecko"
		projectDir := filepath.Join(examplesDir, projectName)
		newProjectTest(t, projectDir).
			ExecuteMakeTarget("build").
			CopyFile("tmp/x.wasm").
			CopyFile("config.toml").
			RenderTemplateFileFromEnvWithCleanup(
				"fn-call.json",
				[]string{"YOUR_COINGECKO_API_KEY"}).
			RunScript(filepath.Join(".", "scripts", projectName+".txtar"))
	})

	t.Run("coingecko example live test", func(t *testing.T) {
		projectName := "coin_prices_from_coingecko"
		projectDir := filepath.Join(examplesDir, projectName)
		newProjectTest(t, projectDir).
			ExecuteMakeTarget("build").
			CopyFile("tmp/x.wasm").
			RenderTemplateFileFromEnvWithCleanup(
				"fn-call.json",
				[]string{"YOUR_COINGECKO_API_KEY"}).
			RenderTemplateStringFromEnvWithCleanup(
				configTemplate,
				"config.toml",
				[]string{
					"LIVE_TEST_PLATFORM",
					"LIVE_TEST_CODE",
					"LIVE_TEST_AUTH_TOKEN",
					"LIVE_TEST_HOST",
				}).
			RunScript(filepath.Join(".", "scripts", projectName+".txtar"))
	})
}
