package entity

const GAME_TYPE_DIVISION_A = 1
const GAME_TYPE_DIVISION_B = 2
const GAME_TYPE_PLAYOFF_STAGE_1 = 3
const GAME_TYPE_PLAYOFF_SEMIFINAL = 4
const GAME_TYPE_PLAYOFF_FINAL = 5

type Game struct {
	ID           int
	TournamentID int
	Team1ID      int
	Team2ID      int
	GameType     int
	WinnerId     *int
}
