package test

import (
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestCoinPricesFromCoingecko(t *testing.T) {
	exampleName := "coin_prices_from_coingecko"
	exampleDir := filepath.Join("..", exampleName)

	work := newTestSetup(t, exampleDir).
		RunMake("build").
		CopyFile("config.toml", "config.toml").
		RenderFile("fn-call.json", []string{"YOUR_COINGECKO_API_KEY"}).
		CopyFile("tmp/x.wasm", "tmp/x.wasm") // matches what fn-call.json expects

	testscript.Run(t, testscript.Params{
		Files: []string{
			filepath.Join(".", "scripts", exampleName+".txtar"),
		},
		RequireExplicitExec: true,
		Setup:               work.SetupFunc(),
	})
}
