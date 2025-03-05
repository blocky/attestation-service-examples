# Getting Esports Data From PandaScore

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the PandaScore API and processes
it.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
and [Error Handling](../error_handling/README.md) examples.

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
you the result of the StartCraft II PL Invitational 2025 tournament final match:

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
bet on the outcome of the StartCraft II PL Invitational 2025 tournament
final match.

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

### Step 2: Create a parameterized oracle function

We'll implement the oracle as `scoreFunc` in
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the [`fn-call.json`](./fn-call.json) 
file contents:

```json
[
  {
    "code_file": "tmp/x.wasm",
    "function": "scoreFunc",
    "input": {
      "match_id": "1121861"
    },
    "secret": {
      "api_key": "PandaScore API Key"
    }
  }
]
```

As you see, we already have the `match_id` value from the previous step in
[`fn-call.json`](./fn-call.json). If you want to look up the results for a 
different match you update the `match_id` value to another id. 
If you haven't already as part of the [Setup](#setup), go ahead and replace
the `api_key` value with your PandaScore API key.

Next, we define the `scoreFunc` function in [`main.go`](./main.go):

```go
type Args struct {
	MatchID string `json:"match_id"`
}

type SecretArgs struct {
	PandaScoreAPIKey string `json:"api_key"`
}

//export scoreFunc
func scoreFunc(inputPtr, secretPtr uint64) uint64 {
	var input Args
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return WriteError(outErr)
	}

	var secret SecretArgs
	secretData := as.Bytes(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	result, err := getMatchResult(input.MatchID, secret.PandaScoreAPIKey)
	if err != nil {
		outErr := fmt.Errorf("getting price: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(result)
}
```

The function takes two `uint64` arguments and returns a `uint64`. These are fat
pointers to shared memory managed by the Blocky AS server, where the first 32
bits are a memory address and the second 32 bits are the size of the data. The
memory space is sandboxed and shared between the TEE host program (Blocky AS
server) and the WASM runtime (your function). The `inputPtr` and `secretPtr`
arguments carry serialized `input` and `secret` sections of
[`fn-call.json`](./fn-call.json).

To parse the `input` data, we first fetch the data pointed to by `inputPtr`
using `as.Bytes` and then unmarshal it into the `Args` struct. We do the same
for the `secret` data. Next, we call the `getMatchResultFromPandaScore` function
to fetch the price of `input.MatchID` using the `secret.PandaScoreAPIKey` API
key. Finally, we return the `matchResult` to user by converting its data to fat
pointer using the `WriteOutput` function and returning the pointer from
`scoreFunc` to the Blocky AS server host runtime.

### Step 3: Make a request to the PandaScore API

The `getMatchResult` function in `scoreFunc` will make an HTTP request to the
PandaScore API to fetch match result data.

Let's start by setting up a struct to parse the relevant fields from the
PandaScore API response:

```go
type PandaScoreMatchResponse struct {
	EndAt    time.Time `json:"end_at"`
	Status   string    `json:"status"`
	WinnerID int       `json:"winner_id"`
	Id       int       `json:"id"`
	Slug     string    `json:"slug"`
	League   struct {
		Slug string `json:"slug"`
	} `json:"league"`
	Serie struct {
		Slug string `json:"slug"`
	} `json:"serie"`
	Tournament struct {
		Slug string `json:"slug"`
	} `json:"tournament"`
	Results []struct {
		PlayerId int `json:"player_id"`
		Score    int `json:"score"`
	} `json:"results"`
	Opponents []struct {
		Opponent struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"opponent"`
	} `json:"opponents"`

```

Next, we define the `getMatchResult` to fetch and parse the data from the
PandaScore API:

```go
type MatchResult struct {
	League     string `json:"league"`
	Serie      string `json:"serie"`
	Tournament string `json:"tournament"`
	Match      string `json:"match"`
	MatchID    int    `json:"match_id"`
	Winner     string `json:"winner"`
	Loser      string `json:"loser"`
	Score      string `json:"score"`
	EndAt      string `json:"end_at"`
}

func getMatchResultFromPandaScore(matchID string, apiKey string) (MatchResult, error) {
	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("https://api.pandascore.co/matches/%s", matchID),
		Headers: map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {"Bearer " + apiKey},
		},
	}
	resp, err := as.HostFuncHTTPRequest(req)
	switch {
	case err != nil:
		return MatchResult{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return MatchResult{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	match := PandaScoreMatchResponse{}
	err = json.Unmarshal(resp.Body, &match)
	if err != nil {
		return MatchResult{}, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}

	if match.Status != "finished" {
		return MatchResult{}, fmt.Errorf("match is not finished")
	}

	winnerOpponent, loserOpponent := match.Opponents[0].Opponent, match.Opponents[1].Opponent
	if winnerOpponent.Id != match.WinnerID {
		winnerOpponent, loserOpponent = loserOpponent, winnerOpponent
	}

	winnerResult, loserResult := match.Results[0], match.Results[1]
	if winnerResult.PlayerId != winnerOpponent.Id {
		winnerResult, loserResult = loserResult, winnerResult
	}

	return MatchResult{
		League:     match.League.Slug,
		Serie:      match.Serie.Slug,
		Tournament: match.Tournament.Slug,
		Match:      match.Slug,
		MatchID:    match.Id,
		Winner:     winnerOpponent.Name,
		Loser:      loserOpponent.Name,
		Score:      fmt.Sprintf("%d - %d", winnerResult.Score, loserResult.Score),
		EndAt:      match.EndAt.Format(time.RFC3339),
	}, nil
}
```

The `getMatchResult` function takes in the `matchID` and `apiKey` as arguments.
First, it constructs an HTTP request to the PandaScore API using the `matchID`
in the URL and the `apiKey` in the headers. It then sends the request to the
`as.HostFuncHTTPRequest` function, which makes the request through the Blocky AS
server networking stack. Next, it checks the response status code and
unmarshalls the JSON response into the `PandaScoreMatchResponse` struct.
Finally, it processes the response to populate the `MatchResult` struct and
returns it to the `scoreFunc` function. The `scoreFunc` function returns a
`Result` containing the `MatchResult` to the Blocky AS server to create an
attestation over the function call and the `Result` struct.

### Step 4: Run the oracle

To run `scoreFunc`, you need call:

```bash
make run
```

You'll see output similar to the following:

```json
{
  "Success": true,
  "Error": "",
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
fit your own needs. If you remember, the application we had in mind for this 
example was to settle esports bets on chain. If you want to take expand this 
example with an on chain component, you may explore the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example to learn you can bring the `MatchResult` struct into a smart contract,
which you could extend to accept and settle esports bets.
