package rimble

import (
	"fmt"
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
	for _, team := range match.Teams {
		if team.WinResult == 1 {
			return team.Name, nil
		}
	}

	return "", fmt.Errorf("no match winner found")
}

func (match MatchData) GameNumbersForMap(mapName string) ([]int, error) {
	mapFound := false

	var gameNumbers []int
	for _, game := range match.Metadata.Games {
		if game.MapName == mapName {
			mapFound = true
			gameNumbers = append(gameNumbers, game.GameNumber)
		}
	}

	if !mapFound {
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

	totalKills := 0
	playerFound := false
	for _, gameNumber := range gameNumbers {
		for _, team := range match.Teams {
			for _, player := range team.Players {
				if player.Username == playerUsername {
					playerFound = true
					for _, result := range player.Results {
						if result.GameNumber == gameNumber {
							totalKills += result.Kills
						}
					}
				}
			}
		}
	}

	if !playerFound {
		return 0, fmt.Errorf("player '%s' not found in match data", playerUsername)
	}

	return totalKills, nil
}

func (match MatchData) TeamKillsOnMap(mapName string, teamName string) (int, error) {
	totalKills := 0

	for _, team := range match.Teams {
		if team.Name == teamName {
			for _, player := range team.Players {
				playerKills, err := match.PlayerKillsOnMap(mapName, player.Username)
				if err != nil {
					return 0, fmt.Errorf(
						"getting kills for player '%s' in team '%s': %w",
						player.Username,
						team.Name,
						err,
					)
				}
				totalKills += playerKills
			}
			return totalKills, nil
		}
	}

	return 0, fmt.Errorf("team '%s' not found in match data", teamName)
}

func (match MatchData) TeamsOnMap(mapName string) ([]string, error) {
	mapFound := false

	var teamNames []string
	for _, game := range match.Metadata.Games {
		if game.MapName == mapName {
			mapFound = true
			for _, team := range match.Teams {
				teamNames = append(teamNames, team.Name)
			}
		}
	}

	if !mapFound {
		return nil, fmt.Errorf("map '%s' not found in match data", mapName)
	}

	return teamNames, nil
}

func (match MatchData) TeamKillDifferenceOnMap(
	mapName string,
	teams []string,
) (
	int,
	error,
) {
	if len(teams) != 2 {
		err := fmt.Errorf("expected 2 teams, got %d", len(teams))
		return 0, err
	}

	teamKills := make(map[string]int)
	for _, team := range teams {
		kills, err := match.TeamKillsOnMap(mapName, team)
		if err != nil {
			err = fmt.Errorf("getting kills for team '%s': %w", team, err)
			return 0, err
		}
		teamKills[team] = kills
	}

	return teamKills[teams[0]] - teamKills[teams[1]], nil
}
