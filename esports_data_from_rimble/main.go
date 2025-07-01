package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blocky/attestation-service-examples/esports-data-from-rimble/rimble"
	"github.com/blocky/basm-go-sdk/basm"
)

func getRecentMatchesFromRimble(apiKey string) ([]rimble.MatchData, error) {
	rimbleURL := "https://rimbleanalytics.com/raw/csgo/completed-matches/"
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL:    rimbleURL,
		Headers: map[string][]string{
			"Accept":    {"application/json"},
			"x-api-key": {apiKey},
		},
	}
	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return nil, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	var matches []rimble.MatchData
	err = json.Unmarshal(resp.Body, &matches)
	if err != nil {
		return nil, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}
	return matches, nil
}

type MatchWinner struct {
	MatchID string
	Date    string
	Winner  string
}

func TeamWinner(match rimble.MatchData) (MatchWinner, error) {
	team, err := match.Winner()
	if err != nil {
		return MatchWinner{}, fmt.Errorf("getting team winner: %w", err)
	}

	return MatchWinner{
		MatchID: match.MatchID,
		Date:    match.Date,
		Winner:  team.Name,
	}, nil
}

type SecretArgs struct {
	RimbleAPIKey string `json:"api_key"`
}

//export matchWinnerFromRimble
func matchWinnerFromRimble(_ uint64, secretPtr uint64) uint64 {
	var secret SecretArgs
	secretData := basm.ReadFromHost(secretPtr)
	err := json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	matches, err := getRecentMatchesFromRimble(secret.RimbleAPIKey)
	switch {
	case err != nil:
		outErr := fmt.Errorf("getting recent matches: %w", err)
		return WriteError(outErr)
	case len(matches) == 0:
		outErr := fmt.Errorf("no recent matches found")
		return WriteError(outErr)
	}

	matchWinner, err := TeamWinner(matches[0])
	if err != nil {
		outErr := fmt.Errorf("getting match winner: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(matchWinner)
}

type TeamKillDiff struct {
	MatchID  string
	Date     string
	Team1    string
	Team2    string
	KillDiff int
}

func TeamKillDifference(match rimble.MatchData) (TeamKillDiff, error) {
	killDiff, err := match.TeamKillDifferenceInGames(match.Metadata.Games)
	if err != nil {
		return TeamKillDiff{}, fmt.Errorf("getting team kill difference: %w", err)
	}

	if killDiff < 1 {
		return TeamKillDiff{
			MatchID:  match.MatchID,
			Date:     match.Date,
			Team1:    match.Teams[1].Name,
			Team2:    match.Teams[0].Name,
			KillDiff: -killDiff,
		}, nil
	}

	return TeamKillDiff{
		MatchID:  match.MatchID,
		Date:     match.Date,
		Team1:    match.Teams[0].Name,
		Team2:    match.Teams[1].Name,
		KillDiff: killDiff,
	}, nil
}

//export teamKillDifferenceFromRimble
func teamKillDifferenceFromRimble(_ uint64, secretPtr uint64) uint64 {
	var secret SecretArgs
	secretData := basm.ReadFromHost(secretPtr)
	err := json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	matches, err := getRecentMatchesFromRimble(secret.RimbleAPIKey)
	switch {
	case err != nil:
		outErr := fmt.Errorf("getting recent matches: %w", err)
		return WriteError(outErr)
	case len(matches) == 0:
		outErr := fmt.Errorf("no recent matches found")
		return WriteError(outErr)
	}

	teamKillDiff, err := TeamKillDifference(matches[0])
	if err != nil {
		outErr := fmt.Errorf("getting team kill difference: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(teamKillDiff)
}

func main() {}
