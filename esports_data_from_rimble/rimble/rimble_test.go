package rimble_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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

func TestMatchData_TeamWinner(t *testing.T) {
	match, err := fetchRawMatchData("2025-02-18", "2379357", RIMBLE_DEMO_API_KEY)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// when
		winner, err := match.TeamWinner()

		// then
		require.NoError(t, err)
		assert.Equal(t, "MOUZ", winner)
	})

	t.Run("no winner", func(t *testing.T) {
		// given
		noWinnerMatch := rimble.MatchData{}

		// when
		_, err := noWinnerMatch.TeamWinner()

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("no match winner"), err.Error())
	})
}

func TestMatchData_GameNumbersForMap(t *testing.T) {
	match, err := fetchRawMatchData("2025-02-18", "2379357", RIMBLE_DEMO_API_KEY)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		mapName := "Mirage"

		// when
		gameNumbers, err := match.GameNumbersForMap(mapName)

		// then
		require.NoError(t, err)
		assert.Equal(t, []int{3}, gameNumbers)
	})

	t.Run("map not found", func(t *testing.T) {
		// given
		mapName := "Non existent map"

		// when
		_, err = match.GameNumbersForMap(mapName)

		// then
		require.ErrorContains(t, err, "not found")
	})
}

func TestMatchData_PlayerKillsOnMap(t *testing.T) {
	match, err := fetchRawMatchData("2025-02-18", "2379357", RIMBLE_DEMO_API_KEY)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		mapName := "Mirage"
		playerUsername := "Brollan"

		// when
		kills, err := match.PlayerKillsOnMap(mapName, playerUsername)

		// then
		require.NoError(t, err)
		assert.Equal(t, 10, kills)
	})

	t.Run("map not found", func(t *testing.T) {
		// given
		mapName := "Non existent map"
		playerUsername := "Brollan"

		// when
		_, err = match.PlayerKillsOnMap(mapName, playerUsername)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("map .* not found"), err.Error())
	})

	t.Run("player not found", func(t *testing.T) {
		// given
		mapName := "Mirage"
		playerUsername := "Non existent player"

		// when
		_, err = match.PlayerKillsOnMap(mapName, playerUsername)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("player .* not found"), err.Error())
	})
}

func TestMatchData_TeamKillsOnMap(t *testing.T) {
	match, err := fetchRawMatchData("2025-02-18", "2379357", RIMBLE_DEMO_API_KEY)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		mapName := "Mirage"
		team := "MOUZ"

		// when
		kills, err := match.TeamKillsOnMap(mapName, team)

		// then
		require.NoError(t, err)
		assert.Equal(t, 68, kills)
	})

	t.Run("map not found", func(t *testing.T) {
		// given
		mapName := "Non existent map"
		team := "MOUZ"

		// when
		_, err = match.TeamKillsOnMap(mapName, team)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("map .* not found"), err.Error())
	})

	t.Run("team not found", func(t *testing.T) {
		// given
		mapName := "Mirage"
		team := "Non existent team"

		// when
		_, err = match.TeamKillsOnMap(mapName, team)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("team .* not found"), err.Error())
	})
}

func TestMatchData_TeamsOnMap(t *testing.T) {
	match, err := fetchRawMatchData("2025-02-18", "2379357", RIMBLE_DEMO_API_KEY)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		mapName := "Mirage"

		// when
		teams, err := match.TeamsOnMap(mapName)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"MOUZ", "Virtus.pro"}, teams)
	})

	t.Run("map not found", func(t *testing.T) {
		// given
		mapName := "Non existent map"

		// when
		_, err = match.TeamsOnMap(mapName)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("map .* not found"), err.Error())
	})
}

func TestMatchData_TeamKillDifferenceOnMap(t *testing.T) {
	match, err := fetchRawMatchData("2025-02-18", "2379357", RIMBLE_DEMO_API_KEY)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		mapName := "Mirage"
		teams, err := match.TeamsOnMap(mapName)
		require.NoError(t, err)

		// when
		teamKillDiff, err := match.TeamKillDifferenceOnMap(mapName, teams)

		// then
		require.NoError(t, err)
		assert.Equal(t, 34, teamKillDiff)
	})

	t.Run("happy path with swapped teams", func(t *testing.T) {
		// given
		mapName := "Mirage"
		teams := []string{"Virtus.pro", "MOUZ"}

		// when
		teamKillDiff, err := match.TeamKillDifferenceOnMap(mapName, teams)

		// then
		require.NoError(t, err)
		assert.Equal(t, -34, teamKillDiff)
	})

	t.Run("wrong number of teams", func(t *testing.T) {
		// given
		mapName := "Mirage"
		teams := []string{"MOUZ"}

		// when
		_, err = match.TeamKillDifferenceOnMap(mapName, teams)

		// then
		require.ErrorContains(t, err, "expected 2 teams, got 1")
	})

	t.Run("map not found", func(t *testing.T) {
		// given
		mapName := "Non existent map"
		teams := []string{"MOUZ", "Virtus.pro"}

		// when
		_, err = match.TeamKillDifferenceOnMap(mapName, teams)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("map .* not found"), err.Error())
	})
}
