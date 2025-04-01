package rimble

import (
	"errors"
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

type RimbleStat struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func GetMatchWinner(match MatchData) (
	RimbleStat,
	error,
) {
	for _, team := range match.Teams {
		if team.WinResult == 1 {
			return RimbleStat{
				Name:  "Winner",
				Value: team.Name,
			}, nil
		}
	}

	return RimbleStat{}, errors.New("no winner found in match data")
}

// func GetTeamKills(matchID, date, apiKey, teamName string) (int, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	totalKills := 0
// 	teamFound := false
//
// 	for _, team := range match.Teams {
// 		if team.Name == teamName {
// 			teamFound = true
// 			for _, result := range team.Results {
// 				totalKills += result.TotalKills
// 			}
// 			break
// 		}
// 	}
//
// 	if !teamFound {
// 		return 0, fmt.Errorf("team '%s' not found in match data", teamName)
// 	}
//
// 	return totalKills, nil
// }
//
// func GetPlayerKills(matchID, date, apiKey, playerUsername string) (int, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	totalKills := 0
// 	playerFound := false
//
// 	for _, team := range match.Teams {
// 		for _, player := range team.Players {
// 			if player.Username == playerUsername {
// 				playerFound = true
// 				for _, result := range player.Results {
// 					totalKills += result.Kills
// 				}
// 				break
// 			}
// 		}
// 		if playerFound {
// 			break
// 		}
// 	}
//
// 	if !playerFound {
// 		return 0, fmt.Errorf("player '%s' not found in match data", playerUsername)
// 	}
//
// 	return totalKills, nil
// }
//
// func GetMapWinner(matchID, date, apiKey, mapName string) (string, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if len(matches) == 0 {
// 		return "", fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	gameNumber := 0
// 	mapFound := false
//
// 	for _, game := range match.Metadata.Games {
// 		if game.MapName == mapName {
// 			gameNumber = game.GameNumber
// 			mapFound = true
// 			break
// 		}
// 	}
//
// 	if !mapFound {
// 		return "", fmt.Errorf("map '%s' not found in match data", mapName)
// 	}
//
// 	var team1, team2 Team
// 	for _, team := range match.Teams {
// 		if team.Designation == 1 {
// 			team1 = team
// 		} else if team.Designation == 2 {
// 			team2 = team
// 		}
// 	}
//
// 	var team1Rounds, team2Rounds int
//
// 	for _, result := range team1.Results {
// 		if result.GameNumber == gameNumber {
// 			team1Rounds = result.RoundsWon
// 			break
// 		}
// 	}
//
// 	for _, result := range team2.Results {
// 		if result.GameNumber == gameNumber {
// 			team2Rounds = result.RoundsWon
// 			break
// 		}
// 	}
//
// 	if team1Rounds > team2Rounds {
// 		return team1.Name, nil
// 	} else if team2Rounds > team1Rounds {
// 		return team2.Name, nil
// 	}
//
// 	return "Draw", nil
// }
//
// func GetPlayerAvgADR(matchID, date, apiKey, playerUsername string) (float64, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	totalADR := 0.0
// 	gamesPlayed := 0
// 	playerFound := false
//
// 	for _, team := range match.Teams {
// 		for _, player := range team.Players {
// 			if player.Username == playerUsername {
// 				playerFound = true
// 				for _, result := range player.Results {
// 					totalADR += result.ADR
// 					gamesPlayed++
// 				}
// 				break
// 			}
// 		}
// 		if playerFound {
// 			break
// 		}
// 	}
//
// 	if !playerFound {
// 		return 0, fmt.Errorf("player '%s' not found in match data", playerUsername)
// 	}
//
// 	if gamesPlayed == 0 {
// 		return 0, nil
// 	}
//
// 	return totalADR / float64(gamesPlayed), nil
// }
//
// func GetKillDifferentialBetweenTeams(matchID, date, apiKey string) (int, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	team1Kills := 0
// 	team2Kills := 0
//
// 	for _, team := range match.Teams {
// 		if team.Designation == 1 {
// 			for _, result := range team.Results {
// 				team1Kills += result.TotalKills
// 			}
// 		} else if team.Designation == 2 {
// 			for _, result := range team.Results {
// 				team2Kills += result.TotalKills
// 			}
// 		}
// 	}
//
// 	return team1Kills - team2Kills, nil
// }
//
// func GetTopFragger(matchID, date, apiKey string) (string, int, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return "", 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return "", 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	topKills := 0
// 	topFragger := ""
//
// 	for _, team := range match.Teams {
// 		for _, player := range team.Players {
// 			playerKills := 0
// 			for _, result := range player.Results {
// 				playerKills += result.Kills
// 			}
//
// 			if playerKills > topKills {
// 				topKills = playerKills
// 				topFragger = player.Username
// 			}
// 		}
// 	}
//
// 	if topFragger == "" {
// 		return "", 0, fmt.Errorf("no players found in match data")
// 	}
//
// 	return topFragger, topKills, nil
// }
//
// func GetRoundsWonByTeam(matchID, date, apiKey, teamName string) (int, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	totalRounds := 0
// 	teamFound := false
//
// 	for _, team := range match.Teams {
// 		if team.Name == teamName {
// 			teamFound = true
// 			for _, result := range team.Results {
// 				totalRounds += result.RoundsWon
// 			}
// 			break
// 		}
// 	}
//
// 	if !teamFound {
// 		return 0, fmt.Errorf("team '%s' not found in match data", teamName)
// 	}
//
// 	return totalRounds, nil
// }
//
// func GetMatchScore(matchID, date, apiKey string) (string, int, string, int, error) {
// 	matches, err := fetchRawMatchData(matchID, date, apiKey)
// 	if err != nil {
// 		return "", 0, "", 0, err
// 	}
//
// 	if len(matches) == 0 {
// 		return "", 0, "", 0, fmt.Errorf("no match data found")
// 	}
//
// 	match := matches[0]
//
// 	var team1Name string
// 	var team1Score int
// 	var team2Name string
// 	var team2Score int
//
// 	for _, team := range match.Teams {
// 		if team.Designation == 1 {
// 			team1Name = team.Name
// 			team1Score = team.GamesWon
// 		} else if team.Designation == 2 {
// 			team2Name = team.Name
// 			team2Score = team.GamesWon
// 		}
// 	}
//
// 	return team1Name, team1Score, team2Name, team2Score, nil
// }
//
// func main() {
// 	if len(os.Args) < 5 {
// 		fmt.Println("Usage: program [statistic] [matchID] [date] [apiKey] [additional params if needed]")
// 		fmt.Println("Available statistics:")
// 		fmt.Println("  winner - Get match winner")
// 		fmt.Println("  teamkills [teamName] - Get total kills for a team")
// 		fmt.Println("  playerkills [username] - Get total kills for a player")
// 		fmt.Println("  mapwinner [mapName] - Get winner for a specific map")
// 		fmt.Println("  playeradr [username] - Get average ADR for a player")
// 		fmt.Println("  killdiff - Get kill differential between teams")
// 		fmt.Println("  topfragger - Get the player with most kills")
// 		fmt.Println("  teamrounds [teamName] - Get total rounds won by a team")
// 		fmt.Println("  score - Get match score")
// 		return
// 	}
//
// 	statType := os.Args[1]
// 	matchID := os.Args[2]
// 	date := os.Args[3]
// 	apiKey := os.Args[4]
//
// 	switch statType {
// 	case "winner":
// 		winner, err := GetMatchWinner(matchID, date, apiKey)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Match winner: %s\n", winner)
//
// 	case "teamkills":
// 		if len(os.Args) < 6 {
// 			fmt.Println("Error: teamName parameter required")
// 			return
// 		}
// 		teamName := os.Args[5]
// 		kills, err := GetTeamKills(matchID, date, apiKey, teamName)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Total kills for %s: %d\n", teamName, kills)
//
// 	case "playerkills":
// 		if len(os.Args) < 6 {
// 			fmt.Println("Error: username parameter required")
// 			return
// 		}
// 		username := os.Args[5]
// 		kills, err := GetPlayerKills(matchID, date, apiKey, username)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Total kills for %s: %d\n", username, kills)
//
// 	case "mapwinner":
// 		if len(os.Args) < 6 {
// 			fmt.Println("Error: mapName parameter required")
// 			return
// 		}
// 		mapName := os.Args[5]
// 		winner, err := GetMapWinner(matchID, date, apiKey, mapName)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Winner on %s: %s\n", mapName, winner)
//
// 	case "playeradr":
// 		if len(os.Args) < 6 {
// 			fmt.Println("Error: username parameter required")
// 			return
// 		}
// 		username := os.Args[5]
// 		adr, err := GetPlayerAvgADR(matchID, date, apiKey, username)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Average ADR for %s: %.2f\n", username, adr)
//
// 	case "killdiff":
// 		diff, err := GetKillDifferentialBetweenTeams(matchID, date, apiKey)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Kill differential (team1 - team2): %d\n", diff)
//
// 	case "topfragger":
// 		username, kills, err := GetTopFragger(matchID, date, apiKey)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Top fragger: %s with %d kills\n", username, kills)
//
// 	case "teamrounds":
// 		if len(os.Args) < 6 {
// 			fmt.Println("Error: teamName parameter required")
// 			return
// 		}
// 		teamName := os.Args[5]
// 		rounds, err := GetRoundsWonByTeam(matchID, date, apiKey, teamName)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Total rounds won by %s: %d\n", teamName, rounds)
//
// 	case "score":
// 		team1, score1, team2, score2, err := GetMatchScore(matchID, date, apiKey)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 			return
// 		}
// 		fmt.Printf("Match score: %s %d - %s %d\n", team1, score1, team2, score2)
//
// 	default:
// 		fmt.Printf("Unknown statistic type: %s\n", statType)
// 	}
// }
