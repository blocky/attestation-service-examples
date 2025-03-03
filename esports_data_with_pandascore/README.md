# Getting Esports Data with PandaScore

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the PandScore API and processes
it.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
[Error Handling](../error-handling/README.md),
and
[Fetch and Process A Call to CoinGecko API](../fetch_and_process_api_call/README.md)
examples.

In this example, you'll learn how to:

- Fetch data from the PandaScore API.

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the PandaScore API](https://app.pandascore.co/dashboard)
  and set it in `fn-call.json` in the `api_key` field.

## Quick Start

To run this example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response showing
you the price of Bitcoin in USD on the Binance market:

```json
{
  "Success": true,
  "Value": {
    "league": "starcraft-2-pl-invitational",
    "serie": "starcraft-2-pl-invitational-2025",
    "tournament": "starcraft-2-pl-invitational-2025-playoffs",
    "match": "solar-vs-cure-2025-02-09",
    "match_id": 1121861,
    "winner": "Cure",
    "loser": "Solar",
    "score": "3 - 1",
    "end_at": "2025-02-09T08:24:49Z"
  }
}
```

## Walkthrough

Let's say you're implementing an onchain betting application that allows players
to bet on the outcome of esports matches. In particular, let's say you set up a 
bet on the outcome of the StartCraft II PL Invitational 2025 final match.

### Step 1: Get match ID for the PandaScore API

You can use
[PandaScore API](https://developers.pandascore.co/docs/introduction)
to get the tournament results. PandaScore organizes its results into 
[leagues, series, tournaments, and matches](https://developers.pandascore.co/docs/fundamentals).
Our first step is to get the match ID for the StarCraft II PL Invitational 2025
final match.

If you look at [`scripts/get_match_id.sh`](./scripts/get_match_id.sh), 
you'll see a series of `curl` and `jq` commands that look up, using the
PandaScore API, the match ID for final match in the
`starcraft-2-pl-invitational` league in 2025. You can get the match ID by
running (after replacing `<PandaScore API Key>` with your PandaScore API key):

```bash
PANDASCORE_API_KEY=<PandaScore API Key> ./scripts/get_match_id.sh
```

which will print

```
1121861
```

If you're curious about PandaScore's APIs used throughout this example, check 
out its 
[documentation](https://developers.pandascore.co/reference/get_matches)
page.

### Step 1: Create a parameterized oracle function

We'll implement the oracle as `oracleFunc` in
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the [`fn-call.json`](./fn-call.json) 
file contents:

```json
[
  {
    "code_file": "tmp/x.wasm",
    "function": "oracleFunc",
    "input": {
      "match_id": "match ID"
    },
    "secret": {
      "api_key": "PandaScore API Key"
    }
  }
]
```

Replace `match ID` with the match ID you got from the previous step and
`PandaScore API Key` with your PandaScore API key.

Next, we define the `oracleFunc` function in ['main.go'](./main.go), which calls
the `getMatchResult` function to fetch and parse the match data from the 
PandaScore API. If you're curious how internals of such functions work, visit 
the
[Fetch and Process A Call to CoinGecko API](../fetch_and_process_api_call/README.md)
example. At a high level, he `getMatchResult` function takes the `matchID` and
`apiKey` as arguments, uses them to fetch and parse the match data from
`https://api.pandascore.co/matches` and returns the `MatchResult` struct 
populated with PandaScore data. The `oracleFunc` function returns a `Result`
containing the `MatchResult` to the Blocky AS server to create an
attestation over the function call and the `Result` struct.

### Step 3: Run the oracle

To run `oracleFunc`, you need call:

```bash
make run
```

You'll see output similar to the following:

```json
{
  "Success": true,
  "Value": {
    "league": "starcraft-2-pl-invitational",
    "serie": "starcraft-2-pl-invitational-2025",
    "tournament": "starcraft-2-pl-invitational-2025-playoffs",
    "match": "solar-vs-cure-2025-02-09",
    "match_id": 1121861,
    "winner": "Cure",
    "loser": "Solar",
    "score": "3 - 1",
    "end_at": "2025-02-09T08:24:49Z"
  }
}
```

Which lists the `league`, `series`, `tournament`, and `match` names. The output
also lists the `match_id` of 1121861. Finally, we get the `winner`, `loser`,
`score`, and `end_at` of when the game finished.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs.
You may also want to explore the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example to learn you can bring the `MatchResult` struct into a smart contract
to settle the bet there.
