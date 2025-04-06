package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blocky/as-demo/sportradar"
	"github.com/blocky/basm-go-sdk"
)

func getNBAGameSummary(gameID string, apiKey string) (sportradar.NBAGameSummary, error) {
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL: fmt.Sprintf(
			"https://api.sportradar.com/nba/trial/v8/en/games/%s/summary.json?api_key=%s",
			gameID,
			apiKey,
		),
	}

	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return sportradar.NBAGameSummary{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return sportradar.NBAGameSummary{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	return sportradar.MakeNBAGameSummaryFromJSON(resp.Body)
}

type NBAPlayerPointsPerMinuteComparison struct {
	PointsPerMinute []PointsPerMinute `json:"points_per_minute"`
	Winner          string            `json:"winner"`
}

type PointsPerMinute struct {
	Player string  `json:"player"`
	PPM    float64 `json:"ppm"`
}

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

	if len(input.Players) != 2 {
		outErr := fmt.Errorf("exactly two players are required")
		return WriteError(outErr)
	}

	var comparison NBAPlayerPointsPerMinuteComparison
	for _, playerName := range input.Players {
		player, err := gameSummary.Player(playerName)
		if err != nil {
			outErr := fmt.Errorf("getting player: %w", err)
			return WriteError(outErr)
		}

		ppm, err := player.PointsPerMinute()
		if err != nil {
			outErr := fmt.Errorf("getting points per minute: %w", err)
			return WriteError(outErr)
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

	return WriteOutput(comparison)
}

func main() {
}
