# Getting Esports Data From PandaScore

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the PandaScore API and processes
it.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

In this example, you'll learn how to:

- Fetch data from the PandaScore API.

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/v0.1.0-beta.10/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the PandaScore API](https://app.pandascore.co/dashboard)
  and set it in `fn-call.json` in the `api_key` field.


## Walkthrough

Let's say you're implementing an on chain fantasy application that needs to
access esports data such as player statistics, or match results. In particular,
let's say want to bring on chain the outcome of the StartCraft II PL
Invitational 2025 tournament final match.

### Step 1: Get match ID for the PandaScore API

You can use
[PandaScore API](https://developers.pandascore.co/docs/introduction)
to get the tournament results. PandaScore organizes its results into 
[leagues, series, tournaments, and matches](https://developers.pandascore.co/docs/fundamentals).
If you want to look up the match ID for the StarCraft II PL Invitational 2025
final match, you can make a sequence of calls to PandaScore API endpoints to
find the
[league](https://developers.pandascore.co/reference/get_leagues),
[serie](https://developers.pandascore.co/reference/get_series),
[tournament](https://developers.pandascore.co/reference/get_tournaments),
and [match](https://developers.pandascore.co/reference/get_matches).
For this example, we went through this process and pulled the StarCraft II PL
Invitational 2025 final match ID `1121861`.

### Step 2: Create a parameterized oracle function

We'll implement the oracle as `scoreFunc` in
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the [`fn-call.json`](./fn-call.json) 
file contents:

```json
{
  "code_file": "tmp/x.wasm",
  "function": "scoreFunc",
  "input": {
    "matches_api_endpoint": "PandaScore API matches endpoint",
    "match_id": "1121861"
  },
  "secret": {
    "api_key": "PandaScore API Key"
  }
}
```

As you see, we already have the `match_id` value from the previous step in
[`fn-call.json`](./fn-call.json). If you want to look up the results for a 
different match you update the `match_id` value to another ID. 
If you haven't already as part of the [Setup](#setup), go ahead and replace
the `api_key` value with your PandaScore API key.

> PandaScore API [Terms and Conditions](https://pandascore.co/terms-and-condition)
> do not allow us to share their API endpoints publicly. To use this example,
> you need to replace `matches_api_endpoint` in [`fn-call.json`](./fn-call.json)
> with the PandaScore API matches endpoint, which you can find on their
> [List marches API call documentation page](https://developers.pandascore.co/reference/get_matches).

Next, we define the `scoreFunc` function in [`main.go`](./main.go):

```go
type Args struct {
	MatchID string `json:"match_id"`
}

type SecretArgs struct {
	PandaScoreAPIKey string `json:"api_key"`
}

//export scoreFunc
func scoreFunc(inputPtr uint64, secretPtr uint64) uint64 {
	var input Args
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return WriteError(outErr)
	}

	var secret SecretArgs
	secretData := basm.ReadFromHost(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	matchResult, err := getMatchResultFromPandaScore(
		input.MatchID,
		secret.PandaScoreAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting price: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(matchResult)
}
```

First, we get the input parameters and secrets. Next, we call
the `getMatchResultFromPandaScore` function to fetch the price of `input.MatchID`
using the`secret.PandaScoreAPIKey` API key.
Finally, we return the `matchResult` to user by converting its data to fat pointer 
using the `WriteOutput` function and returning the pointer from `scoreFunc`
to the Blocky AS server host runtime.

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
	matchesAPIEndpoint, err := getMatchesAPIEndpoint()
	if err != nil {
		return MatchResult{}, fmt.Errorf("getting matches api endpoint: %w", err)
	}

	req := basm.HTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("%s/%s", matchesAPIEndpoint, matchID),
		Headers: map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {"Bearer " + apiKey},
		},
	}
	resp, err := basm.HTTPRequest(req)
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
First, it looks up the PandaScore API endpoint using the `getMatchesAPIEndpoint`
and constructs an HTTP request to the PandaScore API using the `matchID`
in the URL and the `apiKey` in the headers. It then sends the request to the
`basm.HTTPRequest` function, which makes the request through the Blocky AS
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
fit your own needs. If you want to take expand this 
example with an on chain component, you may explore the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example to learn you can bring the `MatchResult` struct into a smart contract.
