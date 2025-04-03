package rimble

import (
	"fmt"

	"github.com/samber/lo"
)

type PlayerResult struct {
	Kills      int     `json:"kills"`
	GameNumber int     `json:"game_number"`
	KAST       float64 `json:"KAST"`
	Deaths     int     `json:"deaths"`
	ADR        float64 `json:"ADR"`
}

type Player struct {
	Name     string         `json:"name"`
	ID       string         `json:"id"`
	Results  []PlayerResult `json:"results"`
	Username string         `json:"username"`
}

type TeamResult struct {
	GameNumber int `json:"game_number"`
	TotalKills int `json:"totalKills"`
	RoundsWon  int `json:"rounds_won"`
}

type Team struct {
	GamesWon    int          `json:"games_won"`
	WinResult   int          `json:"win_result"`
	Players     []Player     `json:"players"`
	Name        string       `json:"name"`
	ID          string       `json:"id"`
	Designation int          `json:"designation"`
	Results     []TeamResult `json:"results"`
}

type Game struct {
	GameNumber int    `json:"game_number"`
	MapName    string `json:"map_name"`
}

type Metadata struct {
	Games []Game `json:"games"`
}

type MatchData struct {
	Metadata    Metadata `json:"metadata"`
	Team1Name   string   `json:"team_1_name"`
	Team2Name   string   `json:"team_2_name"`
	Teams       []Team   `json:"teams"`
	MatchID     string   `json:"matchid"`
	MatchStatus string   `json:"match_status"`
}

func (match MatchData) TeamWinner() (string, error) {
	winningTeams := lo.Filter(match.Teams, func(team Team, _ int) bool {
		return team.WinResult == 1
	})
	switch {
	case len(winningTeams) == 0:
		return "", fmt.Errorf("no winning team found")
	case len(winningTeams) > 1:
		return "", fmt.Errorf("multiple winning teams found")
	}

	return winningTeams[0].Name, nil
}

func (match MatchData) GameNumbersForMap(mapName string) ([]int, error) {
	gameNumbers := lo.FilterMap(match.Metadata.Games, func(game Game, _ int) (int, bool) {
		if game.MapName == mapName {
			return game.GameNumber, true
		}
		return 0, false
	})

	if len(gameNumbers) == 0 {
		return nil, fmt.Errorf("map '%s' not found in match data", mapName)
	}

	return gameNumbers, nil
}

func (match MatchData) PlayerKillsOnMap(mapName string, playerUsername string) (
	int,
	error,
) {
	gameNumbers, err := match.GameNumbersForMap(mapName)
	if err != nil {
		return 0, fmt.Errorf("getting game numbers for map '%s': %w", mapName, err)
	}

	teamsWithPlayer := lo.Filter(match.Teams, func(team Team, _ int) bool {
		return lo.ContainsBy(team.Players, func(player Player) bool {
			return player.Username == playerUsername
		})
	})

	switch {
	case len(teamsWithPlayer) == 0:
		return 0, fmt.Errorf("player '%s' not found in match data", playerUsername)
	case len(teamsWithPlayer) > 1:
		return 0, fmt.Errorf("player '%s' found in multiple teams", playerUsername)
	}

	player, found := lo.Find(teamsWithPlayer[0].Players, func(player Player) bool {
		return player.Username == playerUsername
	})

	if !found {
		return 0, fmt.Errorf("player '%s' not found in their team", playerUsername)
	}

	results := lo.Filter(player.Results, func(result PlayerResult, _ int) bool {
		return lo.Contains(gameNumbers, result.GameNumber)
	})

	if len(results) == 0 {
		return 0, fmt.Errorf("player '%s' has no results for map '%s'", playerUsername, mapName)
	}

	totalKills := lo.SumBy(results, func(result PlayerResult) int {
		return result.Kills
	})

	return totalKills, nil
}

func (match MatchData) TeamKillsOnMap(mapName string, teamName string) (int, error) {
	teams := lo.Filter(match.Teams, func(team Team, _ int) bool {
		return team.Name == teamName
	})
	switch {
	case len(teams) == 0:
		return 0, fmt.Errorf("team '%s' not found in match data", teamName)
	case len(teams) > 1:
		return 0, fmt.Errorf("team '%s' found multiple times in match data", teamName)
	}

	var playerKillsOnMapError error = nil
	totalKills := lo.SumBy(teams[0].Players, func(player Player) int {
		playerKills, err := match.PlayerKillsOnMap(mapName, player.Username)
		if err != nil {
			err = fmt.Errorf("getting kills for player '%s': %w", player.Username, err)
			playerKillsOnMapError = err
			return 0
		}
		return playerKills
	})

	return totalKills, playerKillsOnMapError
}

func (match MatchData) TeamKillDifferenceOnMap(mapName string) (int, error) {
	if len(match.Teams) != 2 {
		err := fmt.Errorf("expected 2 teams, got %d", len(match.Teams))
		return 0, err
	}

	var teamKillsOnMapError error = nil
	teamKills := lo.Map(match.Teams, func(team Team, _ int) int {
		kills, err := match.TeamKillsOnMap(mapName, team.Name)
		if err != nil {
			err = fmt.Errorf("getting kills for team '%s': %w", team.Name, err)
			teamKillsOnMapError = err
			return 0
		}
		return kills
	})

	if teamKillsOnMapError != nil {
		return 0, teamKillsOnMapError
	}

	return teamKills[0] - teamKills[1], nil
}
