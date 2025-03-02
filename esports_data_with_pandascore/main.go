package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blocky/as-demo/as"
)

type Result struct {
	Success bool
	Value   any
}

func (r Result) jsonMarshalWithError(err error) []byte {
	resultStr := fmt.Sprintf(`{ "Success": false, "Value": "%v" }`, err)
	data := []byte(resultStr)
	return data
}

func writeOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		as.Log(fmt.Sprintf("Error marshalling result: %v", err))
		return writeError(err)
	}
	return as.WriteToHost(data)
}

func writeError(err error) uint64 {
	data := Result{}.jsonMarshalWithError(err)
	return as.WriteToHost(data)
}

type PandaScoreMatchResponse struct {
	EndAt    time.Time `json:"end_at"`
	Status   string    `json:"status"`
	WinnerId int       `json:"winner_id"`
	Id       int       `json:"id"`
	Results  []struct {
		TeamId int `json:"team_id"`
		Score  int `json:"score"`
	} `json:"results"`
	Opponents []struct {
		Opponent struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"opponent"`
	} `json:"opponents"`
}

type MatchResult struct {
	Winner string `json:"winner"`
	Loser  string `json:"loser"`
	Score  string `json:"score"`
	EndAt  string `json:"end_at"`
}

func getMatchResult(matchID string, apiKey string) (MatchResult, error) {
	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("https://api.pandascore.co/matches/%s", matchID),
		Headers: map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {"Bearer " + apiKey},
		},
	}
	resp, err := as.HostFuncHTTPRequest(req)
	if err != nil {
		return MatchResult{}, fmt.Errorf("making http request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
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

	var winner string
	var loser string
	for _, opponent := range match.Opponents {
		if opponent.Opponent.Id == match.WinnerId {
			winner = opponent.Opponent.Name
		} else {
			loser = opponent.Opponent.Name
		}
	}

	// Get the score
	var winnerScore int
	var loserScore int
	for _, result := range match.Results {
		if result.TeamId == match.WinnerId {
			winnerScore = result.Score
		} else {
			loserScore = result.Score
		}
	}

	return MatchResult{
		Winner: winner,
		Loser:  loser,
		Score:  fmt.Sprintf("%d - %d", winnerScore, loserScore),
		EndAt:  match.EndAt.Format(time.RFC3339),
	}, nil
}

type Args struct {
	MatchID string `json:"match_id"`
}

type SecretArgs struct {
	PandaScoreAPIKey string `json:"api_key"`
}

//export oracleFunc
func oracleFunc(inputPtr, secretPtr uint64) uint64 {
	var input Args
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return writeError(outErr)
	}

	var secret SecretArgs
	secretData := as.Bytes(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return writeError(outErr)
	}

	result, err := getMatchResult(input.MatchID, secret.PandaScoreAPIKey)
	if err != nil {
		outErr := fmt.Errorf("getting price: %w", err)
		return writeError(outErr)
	}

	return writeOutput(result)
}

func main() {}
