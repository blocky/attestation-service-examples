package rimble_test

import (
	_ "embed"
	"regexp"
	"testing"

	lom "github.com/samber/lo/mutable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blocky/attestation-service-examples/esports-data-from-rimble/rimble"
)

//go:embed testdata/match_data.json
var matchDataJSON []byte

// todo: add a bit to the README about this package and how to trigger tests and update the test data

func TestMakeMatchDataFromMatchesJSON(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// when
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, match)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		// given
		invalidJSON := []byte("invalid JSON")

		// when
		_, err := rimble.MakeMatchDataFromMatchesJSON(invalidJSON)

		// then
		require.Error(t, err)
	})
}

func TestMatchData_Winner(t *testing.T) {
	match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// when
		winner, err := match.Winner()

		// then
		require.NoError(t, err)
		assert.Equal(t, "MOUZ", winner.Name)
	})

	t.Run("no winner", func(t *testing.T) {
		// given
		noWinnerMatch := rimble.MatchData{}

		// when
		_, err := noWinnerMatch.Winner()

		// then
		require.ErrorContains(t, err, "no winning team found")
	})

	t.Run("multiple winners", func(t *testing.T) {
		// given
		multipleWinnersMatch := rimble.MatchData{
			Teams: []rimble.Team{{WinResult: 1}, {WinResult: 1}},
		}

		// when
		_, err := multipleWinnersMatch.Winner()

		// then
		require.ErrorContains(t, err, "multiple winning teams found")
	})
}

func TestMatchData_GamesOnMap(t *testing.T) {
	match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		mapName := "Mirage"

		// when
		games, err := match.GamesOnMap(mapName)

		// then
		require.NoError(t, err)
		assert.Equal(t, []rimble.Game{{GameNumber: 3, MapName: "Mirage"}}, games)
	})

	t.Run("map not found", func(t *testing.T) {
		// given
		mapName := "Non existent map"

		// when
		_, err = match.GamesOnMap(mapName)

		// then
		require.ErrorContains(t, err, "no games")
	})
}

func TestMatchData_PlayerKillsInGames(t *testing.T) {
	match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
	require.NoError(t, err)
	games, err := match.GamesOnMap("Mirage")
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		playerUsername := "Brollan"

		// when
		kills, err := match.PlayerKillsInGames(games, playerUsername)

		// then
		require.NoError(t, err)
		assert.Equal(t, 10, kills)
	})

	t.Run("nil games", func(t *testing.T) {
		// given
		playerUsername := "Brollan"

		// when
		_, err := match.PlayerKillsInGames(nil, playerUsername)

		// then
		require.ErrorContains(t, err, "no games")
	})

	t.Run("empty games", func(t *testing.T) {
		// given
		playerUsername := "Brollan"

		// when
		_, err := match.PlayerKillsInGames([]rimble.Game{}, playerUsername)

		// then
		require.ErrorContains(t, err, "no games")
	})

	t.Run("player not found", func(t *testing.T) {
		// given
		playerUsername := "Non existent player"

		// when
		_, err = match.PlayerKillsInGames(games, playerUsername)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("player .* not found"), err.Error())
	})

	t.Run("player on multiple teams", func(t *testing.T) {
		// given
		playerUsername := "Brollan"
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		require.NoError(t, err)
		match.Teams = append(match.Teams, rimble.Team{
			Players: []rimble.Player{
				{Username: playerUsername},
			},
		})

		// when
		_, err = match.PlayerKillsInGames(games, playerUsername)

		// then
		require.ErrorContains(t, err, "on multiple teams")
	})

	t.Run("no results for games", func(t *testing.T) {
		// given
		playerUsername := "Brollan"
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		require.NoError(t, err)
		match.Teams[0].Players[0].Results = []rimble.PlayerResult{}

		// when
		_, err = match.PlayerKillsInGames(games, playerUsername)

		// then
		require.ErrorContains(t, err, "no result in resultsForGames")
	})
}

func TestMatchData_TeamKillsInGames(t *testing.T) {
	match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
	require.NoError(t, err)
	games, err := match.GamesOnMap("Mirage")
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		team := "MOUZ"

		// when
		kills, err := match.TeamKillsInGames(games, team)

		// then
		require.NoError(t, err)
		assert.Equal(t, 68, kills)
	})

	t.Run("nil games", func(t *testing.T) {
		// given
		team := "MOUZ"

		// when
		_, err := match.TeamKillsInGames(nil, team)

		// then
		require.ErrorContains(t, err, "no games")
	})

	t.Run("empty games", func(t *testing.T) {
		// given
		team := "MOUZ"

		// when
		_, err := match.TeamKillsInGames([]rimble.Game{}, team)

		// then
		require.ErrorContains(t, err, "no games")
	})

	t.Run("team not found", func(t *testing.T) {
		// given
		team := "Non existent team"

		// when
		_, err = match.TeamKillsInGames(games, team)

		// then
		require.Error(t, err)
		require.Regexp(t, regexp.MustCompile("team .* not found"), err.Error())
	})

	t.Run("team found multiple times", func(t *testing.T) {
		// given
		team := "MOUZ"
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		require.NoError(t, err)
		match.Teams = append(match.Teams, rimble.Team{
			Name: team,
		})

		// when
		_, err = match.TeamKillsInGames(games, team)

		// then
		require.ErrorContains(t, err, "found multiple times")
	})

	t.Run("getting kills for player", func(t *testing.T) {
		// given
		team := "MOUZ"
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		require.NoError(t, err)
		match.Teams[0].Players[0].Results = []rimble.PlayerResult{}

		// when
		_, err = match.TeamKillsInGames(games, team)

		// then
		require.ErrorContains(t, err, "getting kills for player")
	})
}

func TestMatchData_TeamKillDifferenceInGames(t *testing.T) {
	match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
	require.NoError(t, err)
	games, err := match.GamesOnMap("Mirage")
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// when
		teamKillDiff, err := match.TeamKillDifferenceInGames(games)

		// then
		require.NoError(t, err)
		assert.Equal(t, 34, teamKillDiff)
	})

	t.Run("happy path with swapped teams", func(t *testing.T) {
		// given
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		lom.Reverse(match.Teams)

		// when
		teamKillDiff, err := match.TeamKillDifferenceInGames(games)

		// then
		require.NoError(t, err)
		assert.Equal(t, -34, teamKillDiff)
	})

	t.Run("nil games", func(t *testing.T) {
		// when
		_, err := match.TeamKillDifferenceInGames(nil)

		// then
		require.ErrorContains(t, err, "no games")
	})

	t.Run("empty games", func(t *testing.T) {
		// when
		_, err := match.TeamKillDifferenceInGames([]rimble.Game{})

		// then
		require.ErrorContains(t, err, "no games")
	})

	t.Run("wrong number of teams", func(t *testing.T) {
		// given
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		match.Teams = match.Teams[:1]

		// when
		_, err = match.TeamKillDifferenceInGames(games)

		// then
		require.ErrorContains(t, err, "expected 2 teams, got 1")
	})

	t.Run("getting kills for team", func(t *testing.T) {
		// given
		match, err := rimble.MakeMatchDataFromMatchesJSON(matchDataJSON)
		match.Teams[0].Players[0].Results = []rimble.PlayerResult{}

		// when
		_, err = match.TeamKillDifferenceInGames(games)

		// then
		require.ErrorContains(t, err, "getting kills for team")
	})
}
