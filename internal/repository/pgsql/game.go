package pgsql

import (
	"database/sql"
	"fmt"
	"tournament/internal/entity"
)

type GameRepository struct {
	DB        *sql.DB
	TableName string
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{
		DB:        db,
		TableName: "games",
	}
}

func (g *GameRepository) Create(tournamentID int, team1 entity.Team, team2 entity.Team, gameType int) (*entity.Game, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (tournament_id, team1_id, team2_id, game_type)
		VALUES ($1, $2, $3, $4) RETURNING id, tournament_id, team1_id, team2_id, game_type, winner_id
	`, g.TableName)

	game := entity.Game{}
	err := g.DB.QueryRow(query, tournamentID, team1.ID, team2.ID, gameType).Scan(&game.ID, &game.TournamentID, &game.Team1ID, &game.Team2ID, &game.GameType, &game.WinnerId)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (g *GameRepository) GetByTypeGames(tournamentID int, gameType int) ([]entity.Game, error) {
	query := fmt.Sprintf(`
		SELECT id, tournament_id, team1_id, team2_id, game_type, winner_id
		FROM %s WHERE tournament_id = $1 AND game_type = $2
	`, g.TableName)

	rows, err := g.DB.Query(query, tournamentID, gameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []entity.Game
	for rows.Next() {
		game := entity.Game{}
		err := rows.Scan(&game.ID, &game.TournamentID, &game.Team1ID, &game.Team2ID, &game.GameType, &game.WinnerId)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return games, nil
}

func (g *GameRepository) Update(game entity.Game) (*entity.Game, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET winner_id = $1
		WHERE id = $2
		RETURNING id, tournament_id, team1_id, team2_id, game_type, winner_id
	`, g.TableName)

	err := g.DB.QueryRow(query, game.WinnerId, game.ID).Scan(&game.ID, &game.TournamentID, &game.Team1ID, &game.Team2ID, &game.GameType, &game.WinnerId)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (g *GameRepository) GetTopTeams(tournamentID int, gameType int) ([]entity.Team, error) {
	query := `
		SELECT t.id, t.name
		FROM teams t
		JOIN games g ON (g.winner_id = t.id)
		WHERE g.tournament_id = $1 AND g.game_type = $2
		GROUP BY t.id
		ORDER BY COUNT(g.id) DESC
	`

	rows, err := g.DB.Query(query, tournamentID, gameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []entity.Team
	for rows.Next() {
		team := entity.Team{}
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func (g *GameRepository) GetTop4WinnersByType(tournamentID int, gameType int) ([]entity.Team, error) {
	query := `
		SELECT t.id, t.name
		FROM teams t
		JOIN games g ON g.winner_id = t.id
		WHERE g.tournament_id = $1 AND g.game_type = $2
		GROUP BY t.id
		ORDER BY COUNT(g.id) DESC
		LIMIT 4
	`

	rows, err := g.DB.Query(query, tournamentID, gameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []entity.Team
	for rows.Next() {
		team := entity.Team{}
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}

func (g *GameRepository) GetWinnersByType(tournamentID int, gameType int) ([]entity.Team, error) {
	query := `
		SELECT t.id, t.name
		FROM teams t
		JOIN games g ON (g.winner_id = t.id)
		WHERE g.tournament_id = $1 AND g.game_type = $2
	`

	rows, err := g.DB.Query(query, tournamentID, gameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var winners []entity.Team
	for rows.Next() {
		team := entity.Team{}
		err := rows.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, err
		}
		winners = append(winners, team)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return winners, nil
}
