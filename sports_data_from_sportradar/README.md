# Getting Sports Data from SportRadar

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches NBA game results from the SportRadar API
and compares the points per minute of two players.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

In this example, you'll learn how to:

- Fetch and parse data from the SportRadar API

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/v0.1.0-beta.9/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the SportRadar API](https://developer.sportradar.com/getting-started/docs/authentication)
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
  "Error": "",
  "Value": {
    "points_per_minute": [
      {
        "player": "Jayson Tatum",
        "ppm": 0.689655172413793
      },
      {
        "player": "Luka Doncic",
        "ppm": 0.6454091432961967
      }
    ],
    "winner": "Jayson Tatum"
  }
}
```

## Walkthrough

Let's say you want to compare the scoring efficiency of two NBA players, Jayson
Tatum and Luka Doncic, in a specific game. You can use the SportRadar API to
fetch the game data and calculate the points per minute for each player.

SportRadar provides very detailed data about NBA games. In fact, their
[game summary](https://developer.sportradar.com/basketball/reference/nba-game-summary)
endpoint returns 213KB of JSON data, which could get expensive to bring into
a smart contract and parse it there.
In this example, we will use Blocky AS to parse a SportRadar API response and
extract the relevant data to compute the players' points per minute.

### Step 1: Create a parameterized oracle function

We'll implement the oracle as `getNBAPlayersPointsComparison` function in
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the
[`fn-call.json`](./fn-call.json) file contents:

```json
{
  "code_file": "tmp/x.wasm",
  "function": "getNBAPlayersPointsComparison",
  "input": {
    "game_id": "aaa3ddb3-dd1b-459e-a686-d2bfc4408881",
    "players": [
      "Jayson Tatum",
      "Luka Doncic"
    ]
  },
  "secret": {
    "api_key": "SportRadar API Key"
  }
}
```

Notice the `input` section, which contains the parameters for
`getNBAPlayersPointsComparison`, specifically the `game_id` of the game we're
interested in and an array of `players` we want to compare. The `secret` section
contains the `api_key` field, which you should set to your SportRadar API key.

Next, we define the `getNBAPlayersPointsComparison` function in
[`main.go`](./main.go):

```go
type Args struct {
	GameID  string   `json:"game_id"`
	Players []string `json:"players"`
}

type SecretArgs struct {
	SportRadarAPIKey string `json:"api_key"`
}

//export getNBAPlayersPointsComparison
func getNBAPlayersPointsComparison(inputPtr uint64, secretPtr uint64) uint64 {
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

	gameSummary, err := getNBAGameSummary(
		input.GameID,
		secret.SportRadarAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting points leader NBA: %w", err)
		return WriteError(outErr)
	}

	comparison, err := makeNBAPlayerPointsPerMinuteComparison(
		gameSummary,
		input.Players,
	)
	if err != nil {
		outErr := fmt.Errorf("getting points per minute comparison: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(comparison)
}
```

First, we get the input parameters and secrets. Next, we call
the `getNBAGameSummary` function to fetch the game results.
We then call the `makeNBAPlayerPOintsPerMinuteComparison` function to
calculate the points per minute for each player.  Finally, we send the
`comparison` struct to the Blocky AS server for attestation.

The `getNBAPlayersPointsComparison` function uses the 
`getNBAPlayersPointsComparison` function to construct the `comparison` struct.

```go
type NBAPlayerPointsPerMinuteComparison struct {
	PointsPerMinute []PointsPerMinute `json:"points_per_minute"`
	Winner          string            `json:"winner"`
}

type PointsPerMinute struct {
	Player string  `json:"player"`
	PPM    float64 `json:"ppm"`
}

func makeNBAPlayerPointsPerMinuteComparison(
	gameSummary sportradar.NBAGameSummary,
	players []string,
) (
	NBAPlayerPointsPerMinuteComparison,
	error,
) {
	if len(players) != 2 {
		outErr := fmt.Errorf("exactly two players are required")
		return NBAPlayerPointsPerMinuteComparison{}, outErr
	}

	var comparison NBAPlayerPointsPerMinuteComparison
	for _, playerName := range players {
		player, err := gameSummary.Player(playerName)
		if err != nil {
			outErr := fmt.Errorf("getting player: %w", err)
			return NBAPlayerPointsPerMinuteComparison{}, outErr
		}

		ppm, err := player.PointsPerMinute()
		if err != nil {
			outErr := fmt.Errorf("getting points per minute: %w", err)
			return NBAPlayerPointsPerMinuteComparison{}, outErr
		}

		comparison.PointsPerMinute = append(
			comparison.PointsPerMinute,
			PointsPerMinute{
				Player: player.FullName,
				PPM:    ppm,
			},
		)
	}

	comparison.Winner = comparison.PointsPerMinute[0].Player
	if comparison.PointsPerMinute[1].PPM >
		comparison.PointsPerMinute[0].PPM {
		comparison.Winner = comparison.PointsPerMinute[1].Player
	}

	return comparison, nil
}
```

In `makeNBAPlayerPOintsPerMinuteComparison`, we first check that exactly two
players are provided. We then iterate through the players and call the
`PointsPerMinute` method on each player to get their points per minute. Finally,
we compare the points per minute of the two players and set the `winner` field
in the `comparison` struct. The `makeNBAPlayerPointsPerMinuteComparison` relies
on the functions provided by the `sportradar` package in 
[`nba.go`](./sportradar/nba.go) to process the `NBAGameSummary` struct
deserialized from the SportRadar API response.

### Step 3: Run the oracle

To run `getNBAPlayersPointsComparison`, you can call:

```bash
make run
```

You'll see output:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "points_per_minute": [
      {
        "player": "Jayson Tatum",
        "ppm": 0.689655172413793
      },
      {
        "player": "Luka Doncic",
        "ppm": 0.6454091432961967
      }
    ],
    "winner": "Jayson Tatum"
  }
}
```

where `"Success": true,` tells you that the function call was successful and 
the `Value` field gives you a JSON-serialized 
`NBAPlayerPointsPerMinuteComparison` struct, which tells you that Jayson Tatum
is the winner with 0.689655172413793 points per minute.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For example, you can try passing in different players to
`getNBAPlayersPointsComparison`. You can also modify the 
[`nba.go`](./sportradar/nba.go) file to calculate different statistics.
You may also want to explore the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example to learn you can bring the `NBAPlayerPointsPerMinuteComparison` struct
into a smart contract.
