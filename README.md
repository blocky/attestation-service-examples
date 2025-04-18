# Blocky Attestation Service Examples

This repository contains examples of how to use the Blocky Attestation
Service (Blocky AS) to create TEE attestations over API data and offchain
computation. The examples are designed to be easy to understand and modify, so
you can use them as a starting point for your own projects.

Here's a list of the examples in this repository.
If you're new to Blocky AS, we recommend working though these examples in the
following order:

- [Hello World - Attesting a Function Call](./hello_world_attest_fn_call) shows
  you how to use Blocky AS to attest simple function calls: You'll learn how to:
   - Create a function that returns a `"Hello, World!"` message
   - Log messages from the function
   - Invoke functions in the Blocky AS using its `bky-as` CLI
   - Extract function output from the Blocky AS attestation
- [Hello World - Bringing A Blocky AS Function Call Attestation On Chain](./hello_world_on_chain)
  shows you how to bring a Blocky AS function call attestation on chain. You'll
  learn how to:
   - Write a smart contract to verify and parse a function call attestation
   - Test the smart contract locally
   - Deploy the smart contract to Base Sepolia to verify a function call
     attestation on chain and use the attested function output in your smart
     contract
- [Error Handling - Attested Function Calls](./error_handling_attest_fn_call)
  shows you how to handle errors in your
  Blocky AS functions. You'll learn how to:
   - Use the result pattern in your function
   - Return structured data from your function
   - Report errors from your function
   - Log errors in your function
- [Random](./random) shows you how to generate random numbers in your Blocky AS
  functions. You'll learn how to:
  - Generate random numbers in the Blocky AS runtime
  - Create a function that returns a random number within a given range
- [Time](./time) shows you how to retrieve the wall clock time in your Blocky AS
  functions.
- [Getting Coin Prices From CoinGecko](./coin_prices_from_coingecko) shows you
  how to use Blocky AS to fetch and process coin price data from the CoinGecko 
  API. You'll learn how to:
   - Pass in parameters and secrets to your function
   - Make an HTTP request to an external API in your function
   - Parse a JSON response from an API
- [Getting Esports Data From PandaScore](./esports_data_from_pandascore) shows
  you how to use Blocky AS to attest and process esports data from the 
  PandaScore API.
- [Getting Esports Data From Rimble](./esports_data_from_rimble) shows you how
  to use Blocky AS to attest and process esports data from the Rimble API.
  PandaScore API. 
- [Getting Sports Data From SportRadar](./sports_data_from_sportradar) shows
  you how to use Blocky AS to attest and process sports data from the
  SportRadar API.
- [Tracking a Shipment with DHL](./shipment_tracking_with_dhl) shows
  you how to use Blocky AS to track a shipment using the DHL API.
- [Time Weighted Average Price](./time_weighted_average_price) is a more
  advanced example that shows you how to calculate the time weighted average
  price of an asset through iterative calls to Blocky AS.

To learn more about Blocky AS, check out our
[documentation](https://blocky-docs.redocly.app/).
