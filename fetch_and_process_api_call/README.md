# Fetch and process API call

This example shows how to use the Blocky Attestation Service (Blocky AS) to
fetch data from an API and process it. The example uses the Coingecko API to
fetch the current price of Bitcoin in USD on the Binance market.

To get started, install the Blocky AS CLI by following the
[setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
in the Blocky AS documentation.

Next, use `nix` to set up the correct environment for building WASM binaries
with TinyGo. Note that this process will take a while the first time you run it.

```bash
nix-shell
```

Build the `main.go` file to a WASM binary:

```bash
make build
```

Send the WASM function for execution to the Blocky AS:

```bash
make run
```



