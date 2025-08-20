# Getting Esports Data From Rimble

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the Rimble API and processes
it.

Before starting this example, make sure you are familiar with the
[Attesting a Function Call](../attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and [Error Handling](../error_handling_attest_fn_call/README.md) examples.

In this example, you'll learn how to:

- Fetch data from the Rimble API.

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://docs.blocky.rocks/{{{AS_VERSION}}}/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the Rimble API](https://documenter.getpostman.com/view/16449503/Tzm8FvFw#authentication)
  and set it in your environment. You may also use our demo key (below), but
  note that its scope and quota are limited. You can set the key in your
  environment by running:

  ```bash
   export RIMBLE_API_KEY=TU167z1Pwb9SAbUErPZN2aepia1MOsBN3nXbC1eE
   ```

  
## Walkthrough

Let's say you're implementing an on chain fantasy application that needs to
access esports data, like Counter-Strike: Global Offensive (CS:GO) match
results, or more custom statistics like the kill count difference between two
match teams on a specific map.

### Step 1: Get match ID for the Rimble API

This example uses the
[Rimble API](https://documenter.getpostman.com/view/16449503/Tzm8FvFw#682e4cd5-97b3-455d-aa52-51b57a819473)
to get CS:GO match data, but requires a match id and date as input. For the sake
of convenience, we have provided a makefile target `fetch-match-id-and-date` to
grab the id and date of the most recently completed CS:GO match. Of course,
you can look up ids and dates for specific matches if you would prefer.

```bash
# Remember to set your Rimble API key in your environment
export RIMBLE_API_KEY="your-rimble-api-key"

# Fetch the most recent match ID and date
make fetch-match-id-and-date
Fetching recent match ID and date...
Match ID: 2382613
Date: 2025-06-20

# Set the match ID and date as environment variables
export RIMBLE_MATCH_ID="2382613"
export RIMBLE_MATCH_DATE="2025-06-20"
```

### Step 2: Create a parameterized function to attest match winner

We'll implement the oracle in the `matchWinnerFromRimble` function in
[`main.go`](./main.go). 

```go
type SecretArgs struct {
	RimbleAPIKey string `json:"api_key"`
}

type MatchWinnerArgs struct {
	Date    string `json:"date"`
	MatchID string `json:"match_id"`
}

//export matchWinnerFromRimble
func matchWinnerFromRimble(inputPtr uint64, secretPtr uint64) uint64 {
	var input MatchWinnerArgs
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

	match, err := getMatchDataFromRimble(
		input.Date,
		input.MatchID,
		secret.RimbleAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting match data: %w", err)
		return WriteError(outErr)
	}

	matchWinner, err := TeamWinner(match, input.Date)
	if err != nil {
		outErr := fmt.Errorf("getting team kill difference: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(matchWinner)
}
```

The `matchWinnerFromRimble` function uses a helper function
`getMatchDataFromRimble` to fetch the match data from the Rimble API. The
`matchWinnerFromRimble` function then calls the `TeamWinner` function to get the
match winner.

```go
type MatchWinner struct {
	MatchID string
	Date    string
	Winner  string
}

func TeamWinner(match rimble.MatchData, date string) (MatchWinner, error) {
	team, err := match.Winner()
	if err != nil {
		return MatchWinner{}, fmt.Errorf("getting team winner: %w", err)
	}

	return MatchWinner{
		MatchID: match.MatchID,
		Date:    date,
		Winner:  team.Name,
	}, nil
}
```

The `TeamWinner` function uses the `MatchData.Winner` function defined in
[`rimble.go`](./rimble/rimble.go) to get the match winner.

```go
func (match MatchData) Winner() (Team, error) {
	// create a function that returns true if the team is a winner
	teamIsWinner := func(team Team, _ int) bool {
		return team.WinResult == 1
	}

	// filter match.Teams to find all teams that are winners
	winningTeams := lo.Filter(match.Teams, teamIsWinner)
	switch {
	case len(winningTeams) == 0:
		return Team{}, fmt.Errorf("no winning team found")
	case len(winningTeams) > 1:
		return Team{}, fmt.Errorf("multiple winning teams found")
	}

	return winningTeams[0], nil
}
```

In the `rimble` package we use the functional programming paradigm using the
[`samber/lo`](https://github.com/samber/lo) package to process deserialized
Rimble API data. In some ways, this is easier to read and understand than
looping through `MatchData` and checking conditions along the way. If you're
more comfortable looping and testing conditions, feel free to do so and
visit the 
[Getting Coin Prices From CoinGecko](../coin_prices_from_coingecko)
to see how we use that approach to process data from the CoinGecko API.

## Step 3: Attest match winner

As in previous examples, we'll invoke `matchWinnerFromRimble` through the 
`bky-as` CLI. We define [`matchWinner.json`](./match-winner.json.template), which
already contains the `match_id` and the `date` of the match, as well as the 
demo Rimble `api_key`, and then pass it to the `bky-as` CLI.

To invoke `matchWinnerFromRimble`, run:

```bash
make match_winner
```

to get the output:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "MatchID": "2382613",
    "Date": "2025-06-20",
    "Winner": "paiN"
  }
}
```

which tells you that the team `painN` won the match with ID `2382613`
played on 2025-06-20.


### Step 4: Create a parameterized oracle function to attest team kill difference

Now let's say that you want to compute a more custom statistic about the match
like the difference in kills between the two teams on a specific map. We'll
implement this oracle in the `teamKillDiffFromRimble` function in
[`main.go`](./main.go). 

```go
type TeamKillDiffArgs struct {
	Date    string `json:"date"`
	MatchID string `json:"match_id"`
	MapName string `json:"map_name"`
}

//export teamKillDifferenceFromRimble
func teamKillDifferenceFromRimble(inputPtr uint64, secretPtr uint64) uint64 {
	var input TeamKillDiffArgs
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

	match, err := getMatchDataFromRimble(
		input.Date,
		input.MatchID,
		secret.RimbleAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting match data: %w", err)
		return WriteError(outErr)
	}

	teamKillDiff, err := TeamKillDifference(match, input.Date)
	if err != nil {
		outErr := fmt.Errorf("getting team kill difference: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(teamKillDiff)
}
```

The `teamKillDifferenceFromRimble` function uses a helper function
`getMatchDataFromRimble` to fetch the match data from the Rimble API. The
`teamKillDifferenceFromRimble` function then calls the `TeamKillDifference`
function to get the team kill difference for the match.

```go
type TeamKillDiff struct {
	MatchID  string
	Date     string
	Team1    string
	Team2    string
	KillDiff int
}

func TeamKillDifference(
	match rimble.MatchData,
	date string,
) (TeamKillDiff, error) {
	killDiff, err := match.TeamKillDifference(match.Metadata.Games)
	if err != nil {
		return TeamKillDiff{}, fmt.Errorf("getting team kill difference: %w", err)
	}

	if killDiff < 1 {
		return TeamKillDiff{
			MatchID:  match.MatchID,
			Date:     date,
			Team1:    match.Teams[1].Name,
			Team2:    match.Teams[0].Name,
			KillDiff: -killDiff,
		}, nil
	}

	return TeamKillDiff{
		MatchID:  match.MatchID,
		Date:     date,
		Team1:    match.Teams[0].Name,
		Team2:    match.Teams[1].Name,
		KillDiff: killDiff,
	}, nil
}
```

In `TeamKillDifference`, we calculate the difference in kills between `Team1` 
and `Team2` for all combined games played during the match.

### Step 5: Attest team kill difference

As in previous examples, we'll invoke `teamKillDifferenceFromRimble` through the
`bky-as` CLI. We define [`teamKillDiff.json`](./team-kill-diff.json.template), which
already contains the `match_id` and the `date`of the match and the
`map_name` of interest, as well as the demo Rimble `api_key`.

To invoke `teamKillDifferenceFromRimble`, run:

```bash
make team-kill-diff
```

to get the output:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "MatchID": "2382613",
    "Date": "2025-06-20",
    "Team1": "paiN",
    "Team2": "FURIA",
    "KillDiff": 12
  }
}
```

which tells you that the team `Volt` scored 2 more kills than team `paiN`
on the map `Inferno` during the match with ID `2382613` played on 2025-06-20.


### Step 6: Work with Rimble data

The `rimble` package contains the data structures and functions to process
Rimble API data. You can extend [`rimble.go`](./rimble/rimble.go) to add 
functionality to process additional data from the Rimble API and then use it
in your oracle functions. The [`rimble_test.go`](./rimble/rimble_test.go) file
contains tests for the `rimble` package. You can run the tests using:

```bash
make test-rimble
```

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For example, you can add functions to 
[`rimble.go`](rimble/rimble.go) to compute additional game statistics and then
user them in [`main.go`](./main.go) to create additional Blocky AS oracle
functions. You can also expand this example with an on chain component.
Explore the 
[Bringing A Blocky AS Function Call Attestation On Chain](../on_chain/README.md)
example to learn how to bring the match results and statistics into a smart
contract.
