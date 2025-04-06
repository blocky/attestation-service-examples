package sportradar_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/blocky/as-demo/sportradar"
)

//go:embed testdata/nba_game_summary.json
var nbaGameSummaryJSON []byte

func TestMakeNBAGameSummaryFromJSON(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// when
		gameSummary, err := sportradar.MakeNBAGameSummaryFromJSON(nbaGameSummaryJSON)

		// then
		require.NoError(t, err)
		require.NotEmpty(t, gameSummary)
	})

	t.Run("unmarshalling error", func(t *testing.T) {
		// when
		_, err := sportradar.MakeNBAGameSummaryFromJSON([]byte("invalid json"))

		// then
		require.Error(t, err)
	})
}

func TestNBAGameSummary_Player(t *testing.T) {
	gameSummary, err := sportradar.MakeNBAGameSummaryFromJSON(nbaGameSummaryJSON)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		playerName := "Jayson Tatum"
		// when
		player, err := gameSummary.Player(playerName)

		// then
		require.NoError(t, err)
		require.Equal(t, playerName, player.FullName)
	})

	t.Run("player not found", func(t *testing.T) {
		// when
		_, err := gameSummary.Player("Unknown Player")

		// then
		require.ErrorContains(t, err, "not found")
	})

	t.Run("multiple players with same name", func(t *testing.T) {
		// given
		playerName := "Player 1"
		gameSummary := sportradar.NBAGameSummary{
			Home: sportradar.NBATeam{
				Players: []sportradar.NBAPlayer{{FullName: playerName}},
			},
			Away: sportradar.NBATeam{
				Players: []sportradar.NBAPlayer{{FullName: playerName}},
			},
		}

		// when
		_, err := gameSummary.Player(playerName)

		// then
		require.ErrorContains(t, err, "multiple players found")
	})
}

func TestNBAPlayer_PointsPerMinute(t *testing.T) {
	gameSummary, err := sportradar.MakeNBAGameSummaryFromJSON(nbaGameSummaryJSON)
	require.NoError(t, err)

	t.Run("happy path", func(t *testing.T) {
		// given
		playerName := "Jayson Tatum"
		player, err := gameSummary.Player(playerName)
		require.NoError(t, err)

		// when
		pointsPerMinute, err := player.PointsPerMinute()

		// then
		require.NoError(t, err)
		require.Equal(t, 0.689655172413793, pointsPerMinute)
	})

	t.Run("player played 0 minutes", func(t *testing.T) {
		// given
		player := sportradar.NBAPlayer{
			Statistics: sportradar.NBAPlayerStatistics{
				Minutes: "00:00",
			},
		}

		// when
		_, err = player.PointsPerMinute()

		// then
		require.ErrorContains(t, err, "played 0 minutes")
	})

	t.Run("invalid minutes format", func(t *testing.T) {
		// given
		player := sportradar.NBAPlayer{
			Statistics: sportradar.NBAPlayerStatistics{
				Minutes: "invalid",
			},
		}

		// when
		_, err = player.PointsPerMinute()

		// then
		require.ErrorContains(t, err, "getting minutes for player")
	})
}

func TestMinutesToFloat(t *testing.T) {
	happyPathTests := []struct {
		name          string
		minutesString string
		minutesFloat  float64
	}{
		{
			name:          "valid time format",
			minutesString: "12:34",
			minutesFloat:  12.566666666666666,
		},
		{
			name:          "valid time format with leading zero",
			minutesString: "01:02",
			minutesFloat:  1.0333333333333334,
		},
		{
			name:          "valid time format with seconds",
			minutesString: "00:45",
			minutesFloat:  0.75,
		},
		{
			name:          "valid time format with zero seconds",
			minutesString: "10:00",
			minutesFloat:  10.0,
		},
		{
			name:          "valid time format with zero minutes",
			minutesString: "00:00",
			minutesFloat:  0.0,
		},
		{
			name:          "valid time format with leading zero seconds",
			minutesString: "00:01",
			minutesFloat:  0.016666666666666666,
		},
	}

	for _, test := range happyPathTests {
		t.Run(test.name, func(t *testing.T) {
			// when
			minutes, err := sportradar.MinutesToFloat(test.minutesString)

			// then
			require.NoError(t, err)
			require.Equal(t, test.minutesFloat, minutes)
		})
	}

	errorTests := []struct {
		name          string
		minutesString string
	}{
		{
			name:          "invalid time format",
			minutesString: "invalid",
		},
		{
			name:          "invalid time format with extra colon",
			minutesString: "12:34:56",
		},
		{
			name:          "invalid time format with missing seconds",
			minutesString: "12:",
		},
		{
			name:          "invalid time format with missing minutes",
			minutesString: ":34",
		},
		{
			name:          "invalid time format with non-numeric minutes",
			minutesString: "abc:34",
		},
		{
			name:          "invalid time format with non-numeric seconds",
			minutesString: "12:abc",
		},
		{
			name:          "invalid time format with negative minutes",
			minutesString: "-12:34",
		},
		{
			name:          "invalid time format with negative seconds",
			minutesString: "12:-34",
		},
		{
			name:          "invalid time format with non-integer seconds",
			minutesString: "12:34.5",
		},
		{
			name:          "invalid time format with non-integer minutes",
			minutesString: "12.5:34",
		},
		{
			name:          "invalid time format with empty string",
			minutesString: "",
		},
		{
			name:          "invalid time format with whitespace",
			minutesString: "   ",
		},
		{
			name:          "too many seconds",
			minutesString: "12:60",
		},
	}
	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			// when
			_, err := sportradar.MinutesToFloat(test.minutesString)

			// then
			require.Error(t, err)
		})
	}
}
