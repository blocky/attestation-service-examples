# Time-Weighted Average Price

This examples shows how to use the Blocky Attestation Service (Blocky AS) to
compute a time-weighted average price (TWAP) of a token. 

In the example we compute the TWAP based on a series of timestamp price samples.
To collect these samples, we use an iterative pattern, where the result of one 
function call (with a previous set of samples) is passed as input to the next
function call to expand the set of samples.

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the CoinGecko API](https://docs.coingecko.com/reference/setting-up-your-api-key)
  and set it in `fn-call.json` in the `api_key` field.

## Run

To run this example, collect the first price samples by calling:

```bash
make init
```

Then collect additional samples by calling:

```bash
make iteration
```

Finally, compute the TWAP by calling:

```bash
make twap
```

Your output should show the attested TWAP of WETH:

```json
{
  "Result": {
    "Success": true,
    "Error": ""
  },
  "TWAP": 2695.7840977738
}
```
