package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blocky/basm-go-sdk"
)

type SportRaderNBAPointsLeaderResponse struct {
	Season struct {
		Year int    `json:"year"`
		Type string `json:"type"`
	} `json:"season"`
	Name       string `json:"name"`
	Categories []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
}

type Price struct {
	Market    string    `json:"market"`
	CoinID    string    `json:"coin_id"`
	Currency  string    `json:"currency"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func getPointsLeaderNBAFromSportRadar(
	seasonYear string,
	seasonType string,
	apiKey string,
) (
	Price,
	error,
) {
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL: fmt.Sprintf(
			"https://api.sportradar.com/nba/trial/v8/en/seasons/%s/%s/leaders.json?api_key=%s",
			seasonYear,
			seasonType,
			apiKey,
		),
	}
	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return Price{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return Price{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	sportRadarResponse := SportRaderNBAPointsLeaderResponse{}
	err = json.Unmarshal(resp.Body, &sportRadarResponse)
	if err != nil {
		return Price{}, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}

	return Price{}, nil
}

type Args struct {
	SeasonYear string `json:"season_year"`
	SeasonType string `json:"season_type"`
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

	pointsLeaderNBA, err := getPointsLeaderNBAFromSportRadar(
		input.SeasonYear,
		input.SeasonType,
		secret.SportRadarAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting points leader NBA: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(pointsLeaderNBA)
}

func main() {}
