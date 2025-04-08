package sportradar

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

// NBAGameSummary represents the summary of an NBA game.
// Note that we set up this struct to parse only certain fields game summary response
type NBAGameSummary struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Coverage  string    `json:"coverage"`
	Scheduled time.Time `json:"scheduled"`
	Home      NBATeam   `json:"home"`
	Away      NBATeam   `json:"away"`
}

type NBATeam struct {
	Name    string      `json:"name"`
	Market  string      `json:"market"`
	Players []NBAPlayer `json:"players"`
}

type NBAPlayer struct {
	FullName   string              `json:"full_name"`
	Statistics NBAPlayerStatistics `json:"statistics"`
}

type NBAPlayerStatistics struct {
	Minutes string `json:"minutes"`
	Points  int    `json:"points"`
}

func MakeNBAGameSummaryFromJSON(data []byte) (NBAGameSummary, error) {
	var gameSummary NBAGameSummary
	err := json.Unmarshal(data, &gameSummary)
	if err != nil {
		return NBAGameSummary{}, fmt.Errorf("unmarshalling game summary: %w", err)
	}
	return gameSummary, nil
}

func (game NBAGameSummary) Player(playerName string) (NBAPlayer, error) {
	playerHasPlayerName := func(player NBAPlayer, _ int) bool {
		return player.FullName == playerName
	}

	homePlayersWithPlayerName := lo.Filter(game.Home.Players, playerHasPlayerName)
	awayPlayersWithPlayerName := lo.Filter(game.Away.Players, playerHasPlayerName)
	players := append(homePlayersWithPlayerName, awayPlayersWithPlayerName...)
	switch {
	case len(players) == 0:
		err := fmt.Errorf("player '%s' not found", playerName)
		return NBAPlayer{}, err
	case len(players) > 1:
		err := fmt.Errorf("multiple players found with name '%s'", playerName)
		return NBAPlayer{}, err
	}

	return players[0], nil
}

func (player NBAPlayer) PointsPerMinute() (float64, error) {
	minutes, err := MinutesToFloat(player.Statistics.Minutes)
	if err != nil {
		err := fmt.Errorf("getting minutes for player '%s': %w", player.FullName, err)
		return 0, err
	}
	if minutes == 0 {
		err := fmt.Errorf("player '%s' has played 0 minutes", player.FullName)
		return 0, err
	}

	return float64(player.Statistics.Points) / minutes, nil
}

func MinutesToFloat(timeStr string) (float64, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes value: %s", parts[0])
	}
	if minutes < 0 {
		return 0, fmt.Errorf("negative minutes value: %s", parts[0])
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds value: %s", parts[1])
	}
	if seconds < 0 || seconds >= 60 {
		return 0, fmt.Errorf("invalid seconds value: %s", parts[1])
	}

	totalMinutes := float64(minutes) + float64(seconds)/60.0
	return totalMinutes, nil
}
