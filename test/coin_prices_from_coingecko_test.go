package test

import (
	"path/filepath"
	"testing"
)

const examplesDir = ".."

func TestCoinPricesFromCoingecko(t *testing.T) {
	projectName := "coin_prices_from_coingecko"
	projectDir := filepath.Join(examplesDir, projectName)
	newProjectTest(t, projectDir).
		ExecuteMakeTarget("build").
		CopyFile("tmp/x.wasm").
		CopyFile("config.toml").
		RenderTemplateFromEnvWithCleanup("fn-call.json", []string{"YOUR_COINGECKO_API_KEY"}).
		RunScript(filepath.Join(".", "scripts", projectName+".txtar"))
}
