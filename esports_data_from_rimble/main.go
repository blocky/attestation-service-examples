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

type Args struct {
	Statistic string `json:"statistic"`
	Date      string `json:"date"`
	MatchID   string `json:"match_id"`
}

type SecretArgs struct {
	RimbleAPIKey string `json:"api_key"`
}

// todo: udpate function name
//
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

	match, err := getMatchDataFromRimble(
		input.Date,
		input.MatchID,
		secret.RimbleAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting match data: %w", err)
		return WriteError(outErr)
	}

	var stat any
	switch input.Statistic {
	case "match winner":
		stat, err = rimble.GetMatchWinner(match)
		if err != nil {
			outErr := fmt.Errorf("getting match winner: %w", err)
			return WriteError(outErr)
		}
	default:
		err = fmt.Errorf("unsupported statistic: %s", input.Statistic)
		return WriteError(err)
	}

	return WriteOutput(stat)
}

func main() {}
