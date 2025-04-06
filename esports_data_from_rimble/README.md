# Getting Esports Data From Rimble

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the Rimble API and processes
it.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and [Error Handling](../error_handling/README.md) examples.

In this example, you'll learn how to:

- Fetch data from the Rimble API.

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/v0.1.0-beta.6/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the Rimble API](https://documenter.getpostman.com/view/16449503/Tzm8FvFw#authentication)
  and set it in the `api_key` field in `matchWinner.json` and
  `teamKillDiff.json`. For the purpose of this example you can use the demo key
  included in those files.

## Walkthrough

Let's say you're implementing an on chain fantasy application that needs to
access esports data, like Counter-Strike: Global Offensive (CS:GO) match
results, or more custom statistics like the kill count difference between two
match teams on a specific map.

### Step 1: Get match ID for the Rimble API

You can use
[Rimble API](https://documenter.getpostman.com/view/16449503/Tzm8FvFw#682e4cd5-97b3-455d-aa52-51b57a819473)
to get CS:GO matches. Let's say you're interested in match with the ID `2379357`
that took place on 2025-02-18.

### Step 2: Create a parameterized  function to attest match winner

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

The `matchWinnerFromRimble` function uses a helper function `getMatchDataFromRimble`
to fetch the match data from the Rimble API. The `matchWinnerFromRimble` function
then calls the `TeamWinner` function to get the match winner. 

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
	winningTeams := lo.Filter(match.Teams, func(team Team, _ int) bool {
		return team.WinResult == 1
	})
	switch {
	case len(winningTeams) == 0:
		return Team{}, fmt.Errorf("no winning team found")
	case len(winningTeams) > 1:
		return Team{}, fmt.Errorf("multiple winning teams found")
	}

	return winningTeams[0], nil
}
```

As you can see, in the `MatchData.Winner` function we use the imported 
[`samber/lo`](https://github.com/samber/lo) package to process deserialized
Rimble API data.

As in previous examples, we'll invoke `matchWinnerFromRimble` through the 
`bky-as` CLI, specifically using [`matchWinner.json`](./matchWinner.json), which
already contains the `match_id` and the `date` of the match, as well as the 
demo Rimble `api_key`.

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
    "MatchID": "2379357",
    "Date": "2025-02-18",
    "Winner": "MOUZ"
  }
}
```

which tells you that the team `MOUZ` won the match with ID `2379357` played on
2025-02-18.


### Step 2: Create a parameterized oracle function to attest team kill difference

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

	teamKillDiff, err := TeamKillDifferenceOnMap(match, input.Date, input.MapName)
	if err != nil {
		outErr := fmt.Errorf("getting team kill difference: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(teamKillDiff)
}
```

The `teamKillDifferenceFromRimble` function uses a helper function
`getMatchDataFromRimble` to fetch the match data from the Rimble API. The
`teamKillDifferenceFromRimble` function then calls the `TeamKillDifferenceOnMap`
function to get the team kill difference again using the `rimble` package.

We'll invoke `teamKillDifferenceFromRimble` through the
`bky-as` CLI, specifically using [`teamKillDiff.json`](./teamKillDiff.json),
which already contains the `match_id` and the `date` of the match and the
`map_name` of interest, as well as the demo Rimble `api_key`.

To invoke `matchWinnerFromRimble`, run:

```bash
make team-kill-diff
```

to get the output:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "MatchID": "2379357",
    "Date": "2025-02-18",
    "MapName": "Mirage",
    "Team1": "MOUZ",
    "Team2": "Virtus.pro",
    "KillDiff": 34
  }
}
```

which tells you that the team `MOUZ` scored 34 more kills than team `Virtus.pro`
on the map `Mirage` during the match with ID `2379357` played on 2025-02-18.


## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For example, you can add functions to 
[`rimble.go`](rimble/rimble.go) to compute additional game statistics and then
user them in [`main.go`](./main.go) to create additional Blocky AS oracle
functions. You can also expand this example with an on chain component, you may
explore the 
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example to learn you can bring the match results and statistics into a smart
contract.
