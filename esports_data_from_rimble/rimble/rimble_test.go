package rimble_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blocky/attestation-service-examples/esports-data-from-rimble/rimble"
)

var RIMBLE_DEMO_API_KEY = "TU167z1Pwb9SAbUErPZN2aepia1MOsBN3nXbC1eE"

func fetchRawMatchData(date, matchID, apiKey string) (rimble.MatchData, error) {
	url := fmt.Sprintf("https://rimbleanalytics.com/raw/csgo/match-status/?matchid=%s&date=%s", matchID, date)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return rimble.MatchData{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return rimble.MatchData{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return rimble.MatchData{}, fmt.Errorf("API request failed with status code: %d, URL: %s", resp.StatusCode, req.URL.String())
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rimble.MatchData{}, fmt.Errorf("error reading response body: %w", err)
	}

	var matches []rimble.MatchData
	err = json.Unmarshal(body, &matches)
	if err != nil {
		return rimble.MatchData{}, fmt.Errorf("error unmarshaling JSON: %w", err)
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

func TestGetMatchWinner(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		matchID := "2379357"
		date := "2025-02-18"

		matches, err := fetchRawMatchData(date, matchID, RIMBLE_DEMO_API_KEY)
		require.NoError(t, err)

		// when
		winner, err := rimble.GetMatchWinner(matches)

		// then
		require.NoError(t, err)
		assert.Equal(t, "MOUZ", winner)
	})
}

func TestGetPlayerKills(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		matchID := "2379357"
		date := "2025-02-18"
		playerUsername := "Brollan"

		matches, err := fetchRawMatchData(date, matchID, RIMBLE_DEMO_API_KEY)
		require.NoError(t, err)

		// when
		kills, err := rimble.GetPlayerKills(matches, playerUsername)

		// then
		require.NoError(t, err)
		assert.Equal(t, 36, kills)
	})
}

func TestGetTeamKills(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		matchID := "2379357"
		date := "2025-02-18"
		team := "MOUZ"

		matches, err := fetchRawMatchData(date, matchID, RIMBLE_DEMO_API_KEY)
		require.NoError(t, err)

		// when
		kills, err := rimble.GetTeamKills(matches, team)

		// then
		require.NoError(t, err)
		assert.Equal(t, 229, kills)
	})
}
