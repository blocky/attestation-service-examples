package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/blocky/as-demo/as"
)

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
}

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

	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("%s/%s", matchesAPIEndpoint, matchID),
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

func main() {}

func getMatchesAPIEndpoint() (string, error) {
	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("https://developers.pandascore.co/reference/get_matches"),
	}
	resp, err := as.HostFuncHTTPRequest(req)
	switch {
	case err != nil:
		return "", fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return "", fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	matchesURLRegex, err := regexp.Compile(`https.{21}matches`)
	if err != nil {
		return "", fmt.Errorf("compiling regex: %w", err)
	}

	return string(matchesURLRegex.Find(resp.Body)), nil
}
