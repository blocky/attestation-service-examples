# Fetch and process API call

This example shows how to use the Blocky Attestation Service (Blocky AS) to
fetch data from an API and process it. The example uses the Coingecko API to
fetch the current price of Bitcoin in USD on the Binance market.

To get started, install the Blocky AS CLI by following the
[setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
in the Blocky AS documentation. To run this example, make sure you also have
[Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
your system.

To run this example, all you need to do is run:

```bash
make
```

Your output should show the price of Bitcoin in USD on the Binance market:

```json
{
  "IsErr": false,
  "Value": {
    "market": "Binance",
    "coin_id": "BTC",
    "currency": "USD",
    "price": 98420,
    "timestamp": "2025-02-14T19:19:00Z"
  },
  "Error": ""
}
```

> Note: The `tmp/x.wasm` target in the `Makefiile` builds WASM for the
> `linux/amd64` platform. That is the platform used by the hosted Blocky AS
> server. If you're running Blocky AS locally, with `host = "local-server"` in
> `config.yaml`, you may need to update the `tmp/x.wasm` target to build for 
> your architecture.
