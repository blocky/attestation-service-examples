package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blocky/basm-go-sdk"
)

type GameSummary struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Coverage  string    `json:"coverage"`
	Scheduled time.Time `json:"scheduled"`
	Home      struct {
		Name    string `json:"name"`
		Market  string `json:"market"`
		Players []struct {
			FullName   string `json:"full_name"`
			Statistics struct {
				Minutes string `json:"minutes"`
				Points  int    `json:"points"`
			} `json:"statistics"`
		} `json:"players"`
	} `json:"home"`
	Away struct {
		Name    string `json:"name"`
		Market  string `json:"market"`
		Players []struct {
			FullName   string `json:"full_name"`
			Statistics struct {
				Minutes string `json:"minutes"`
				Points  int    `json:"points"`
			} `json:"statistics"`
		} `json:"players"`
	} `json:"away"`
}

func getGameSummary(gameID string, apiKey string) (GameSummary, error) {
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
		return GameSummary{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return GameSummary{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	gameSummary := GameSummary{}
	err = json.Unmarshal(resp.Body, &gameSummary)
	if err != nil {
		return GameSummary{}, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}

	return gameSummary, nil
}

type PointsPerMinute struct {
	PlayerName      string  `json:"player_name"`
	PointsPerMinute float64 `json:"points_per_minute"`
}

func pointsPerMinute(summary GameSummary, playerName string) (PointsPerMinute, error) {
	var playerFound bool
	var playerStats struct {
		Minutes string `json:"minutes"`
		Points  int    `json:"points"`
	}
	for _, player := range summary.Home.Players {
		if player.FullName == playerName {
			playerStats = player.Statistics
			playerFound = true
		}
	}
	for _, player := range summary.Away.Players {
		if player.FullName == playerName {
			playerStats = player.Statistics
			playerFound = true
		}
	}

	if !playerFound {
		return PointsPerMinute{}, fmt.Errorf(`player "%s" not found`, playerName)
	}

	if playerStats.Minutes == "00:00" {
		return PointsPerMinute{}, nil
	}

	minutes, err := minutesToFloat(playerStats.Minutes)
	if err != nil {
		return PointsPerMinute{}, err
	}

	return PointsPerMinute{
		PlayerName:      playerName,
		PointsPerMinute: float64(playerStats.Points) / minutes,
	}, nil
}

type Args struct {
	GameID  string   `json:"game_id"`
	Players []string `json:"players"`
}

type SecretArgs struct {
	SportRadarAPIKey string `json:"api_key"`
}

//export pointsLeaderNBA
func pointsLeaderNBA(inputPtr uint64, secretPtr uint64) uint64 {
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

	gameSummary, err := getGameSummary(
		input.GameID,
		secret.SportRadarAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting points leader NBA: %w", err)
		return WriteError(outErr)
	}

	var pointEfficiency []PointsPerMinute
	for _, player := range input.Players {
		ppm, err := pointsPerMinute(gameSummary, player)
		if err != nil {
			outErr := fmt.Errorf("computing points per minute: %w", err)
			return WriteError(outErr)
		}
		pointEfficiency = append(pointEfficiency, ppm)
	}

	return WriteOutput(pointEfficiency)
}

func main() {
}

func minutesToFloat(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes value: %s", parts[0])
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds value: %s", parts[1])
	}

	totalMinutes := float64(minutes) + float64(seconds)/60.0
	return totalMinutes, nil
}
