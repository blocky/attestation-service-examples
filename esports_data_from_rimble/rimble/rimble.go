package rimble

import (
	"encoding/json"
	"errors"
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

func MakeMatchDataFromMatchesJSON(matchesJSON []byte) (MatchData, error) {
	var matches []MatchData
	err := json.Unmarshal(matchesJSON, &matches)
	if err != nil {
		return MatchData{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	switch len(matches) {
	case 0:
		err = fmt.Errorf("found no matches")
		return MatchData{}, err
	case 1:
		break // only one match found, proceed to return it
	default:
		err = fmt.Errorf("found multiple matches")
		return MatchData{}, err
	}

	return matches[0], nil
}

func (match MatchData) Winner() (Team, error) {
	// create a function that returns true if the team is a winner
	teamIsWinner := func(team Team, _ int) bool {
		return team.WinResult == 1
	}

	// filter match.Teams to find all teams that are winners
	winningTeams := lo.Filter(match.Teams, teamIsWinner)
	switch {
	case len(winningTeams) == 0:
		return Team{}, fmt.Errorf("no winning team found")
	case len(winningTeams) > 1:
		return Team{}, fmt.Errorf("multiple winning teams found")
	}

	return winningTeams[0], nil
}

func (match MatchData) GamesOnMap(mapName string) ([]Game, error) {
	// create a function that returns true if the game was played on mapName
	gameOnMap := func(game Game, _ int) bool {
		return game.MapName == mapName
	}

	// filter match.Metadata.Games to find all games played on mapName
	games := lo.Filter(match.Metadata.Games, gameOnMap)
	if len(games) == 0 {
		return nil, fmt.Errorf("no games played on map '%s'", mapName)
	}

	return games, nil
}

func (match MatchData) PlayerKillsInGames(games []Game, playerUsername string) (
	int,
	error,
) {
	if len(games) == 0 {
		return 0, fmt.Errorf("no games")
	}

	// create a function that returns true if Player has the playerUsername
	playerHasUsername := func(player Player) bool {
		return player.Username == playerUsername
	}

	// create a function that returns true if the team has the player with playerUsername
	teamHasPlayer := func(team Team, _ int) bool {
		return lo.ContainsBy(team.Players, playerHasUsername)
	}

	// filter match.Teams to find all teams that have the player with playerUsername
	teamsWithPlayer := lo.Filter(match.Teams, teamHasPlayer)
	switch {
	case len(teamsWithPlayer) == 0:
		return 0, fmt.Errorf("player '%s' not found on any team", playerUsername)
	case len(teamsWithPlayer) > 1:
		return 0, fmt.Errorf("player '%s' found on multiple teams", playerUsername)
	}

	// find the player with playerUsername in the team
	player, found := lo.Find(teamsWithPlayer[0].Players, playerHasUsername)
	if !found {
		return 0, fmt.Errorf("player '%s' not found in their team", playerUsername)
	}

	// create a function that returns true if the result is for a particular game
	resultForGames := func(result PlayerResult, _ int) bool {
		gameNumbers := lo.Map(games, func(game Game, _ int) int {
			return game.GameNumber
		})
		return lo.Contains(gameNumbers, result.GameNumber)
	}

	// filter player.Results to find all resultsForGames for games
	resultsForGames := lo.Filter(player.Results, resultForGames)
	if len(resultsForGames) == 0 {
		return 0, fmt.Errorf("player '%s' has no result in resultsForGames", playerUsername)
	}

	// sum the kills in the resultsForGames
	totalKills := lo.SumBy(resultsForGames, func(result PlayerResult) int {
		return result.Kills
	})

	return totalKills, nil
}

func (match MatchData) TeamKillsInGames(games []Game, teamName string) (int, error) {
	if len(games) == 0 {
		return 0, fmt.Errorf("no games")
	}

	// create a function that returns true if the team has the teamName
	teamHasName := func(team Team, _ int) bool {
		return team.Name == teamName
	}

	// filter match.Teams to find all teams that have the teamName
	teams := lo.Filter(match.Teams, teamHasName)
	switch {
	case len(teams) == 0:
		return 0, fmt.Errorf("team '%s' not found in match data", teamName)
	case len(teams) > 1:
		return 0, fmt.Errorf("team '%s' found multiple times in match data", teamName)
	}

	// create a function that returns the kills for a player in the games
	// where internal errors are added to the playerKillsInGamesErr
	var playerKillsInGamesErr error = nil
	playerKillsInGames := func(player Player) int {
		playerKills, err := match.PlayerKillsInGames(games, player.Username)
		if err != nil {
			err = fmt.Errorf("getting kills for player '%s': %w", player.Username, err)
			playerKillsInGamesErr = errors.Join(playerKillsInGamesErr, err)
			return 0
		}
		return playerKills
	}

	// sum the kills for all players in the team
	totalKills := lo.SumBy(teams[0].Players, playerKillsInGames)

	return totalKills, playerKillsInGamesErr
}

func (match MatchData) TeamKillDifferenceInGames(games []Game) (int, error) {
	if len(games) == 0 {
		return 0, fmt.Errorf("no games")
	}

	if len(match.Teams) != 2 {
		err := fmt.Errorf("expected 2 teams, got %d", len(match.Teams))
		return 0, err
	}

	// create a function that returns team kills in games
	// where internal errors are added to the teamKillsOnMapError
	var teamKillsOnMapError error = nil
	teamKillsInGames := func(team Team, _ int) int {
		kills, err := match.TeamKillsInGames(games, team.Name)
		if err != nil {
			err = fmt.Errorf("getting kills for team '%s': %w", team.Name, err)
			teamKillsOnMapError = errors.Join(teamKillsOnMapError, err)
			return 0
		}
		return kills
	}

	// for each team, get the kills in the games
	teamKills := lo.Map(match.Teams, teamKillsInGames)

	if teamKillsOnMapError != nil {
		return 0, teamKillsOnMapError
	}

	return teamKills[0] - teamKills[1], nil
}
