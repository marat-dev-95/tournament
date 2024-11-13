package usecase

import "tournament/internal/entity"

type TournamentRepository interface {
	Create(tournament entity.Tournament) (*entity.Tournament, error)
	Delete(tournament entity.Tournament) error
	GetById(id int) (*entity.Tournament, error)
	AddTeam(tournamentID int, team entity.Team) (*entity.Team, error)
	GetTeams(tournamentId int) ([]entity.Team, error)
}

type GameRepository interface {
	Create(tournamentID int, team1 entity.Team, team2 entity.Team, gameType int) (*entity.Game, error)
	GetByTypeGames(tournamentID int, gameType int) ([]entity.Game, error)
	Update(game entity.Game) (*entity.Game, error)
	GetTopTeams(tournamentID int, gameType int) ([]entity.Team, error)
	GetWinnersByType(tournamentID int, gameType int) ([]entity.Team, error)
	GetTop4WinnersByType(tournamentID int, gameType int) ([]entity.Team, error)
}
