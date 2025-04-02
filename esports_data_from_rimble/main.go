package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blocky/attestation-service-examples/esports-data-from-rimble/rimble"
	"github.com/blocky/basm-go-sdk"
)

func getMatchDataFromRimble(
	date string,
	matchID string,
	apiKey string,
) (
	rimble.MatchData,
	error,
) {
	rimbleURL := "https://rimbleanalytics.com/raw/csgo/match-status/"
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("%s?matchid=%s&date=%s", rimbleURL, matchID, date),
		Headers: map[string][]string{
			"Accept":    {"application/json"},
			"x-api-key": {apiKey},
		},
	}
	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return rimble.MatchData{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return rimble.MatchData{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	var matches []rimble.MatchData
	err = json.Unmarshal(resp.Body, &matches)
	if err != nil {
		return rimble.MatchData{}, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}

	switch len(matches) {
	case 0:
		err = fmt.Errorf(
			`no match found for match ID: "%s" on date: "%s"`,
			matchID,
			date,
		)
		return rimble.MatchData{}, err
	case 1:
		break // only one match found, proceed to return it
	default:
		err = fmt.Errorf(
			`multiple matches found for match ID: "%s" on date: "%s"`,
			matchID,
			date,
		)
		return rimble.MatchData{}, err
	}

	return matches[0], nil
}

type MatchWinner struct {
	MatchID string
	Date    string
	Winner  string
}

func TeamWinner(match rimble.MatchData, date string) (MatchWinner, error) {
	team, err := match.TeamWinner()
	if err != nil {
		return MatchWinner{}, fmt.Errorf("getting team winner: %w", err)
	}

	return MatchWinner{
		MatchID: match.MatchID,
		Date:    date,
		Winner:  team,
	}, nil
}

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

type TeamKillDiff struct {
	MatchID  string
	Date     string
	MapName  string
	Team1    string
	Team2    string
	KillDiff int
}

func TeamKillDifferenceOnMap(
	match rimble.MatchData,
	date string,
	mapName string,
) (TeamKillDiff, error) {
	teams, err := match.TeamsOnMap(mapName)
	if err != nil {
		return TeamKillDiff{}, fmt.Errorf("getting teams on map: %w", err)
	}

	killDiff, err := match.TeamKillDifferenceOnMap(mapName, teams)
	if err != nil {
		return TeamKillDiff{}, fmt.Errorf("getting team kill difference: %w", err)
	}

	if killDiff < 1 {
		tmp := teams[0]
		teams[0] = teams[1]
		teams[1] = tmp
		killDiff = -killDiff
	}

	return TeamKillDiff{
		MatchID:  match.MatchID,
		Date:     date,
		MapName:  mapName,
		Team1:    teams[0],
		Team2:    teams[1],
		KillDiff: killDiff,
	}, nil
}

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

func main() {}
