package pgsql

import (
	"database/sql"
	"fmt"
	"tournament/internal/entity"
)

type TournamentRepository struct {
	DB        *sql.DB
	TableName string
}

func NewTournamentRepository(db *sql.DB) *TournamentRepository {
	return &TournamentRepository{
		DB:        db,
		TableName: "tournaments",
	}
}

func (t *TournamentRepository) Create(tournament entity.Tournament) (*entity.Tournament, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", t.TableName)
	err := t.DB.QueryRow(query, tournament.Name).Scan(&tournament.ID)
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (t *TournamentRepository) Delete(tournament entity.Tournament) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", t.TableName)
	_, err := t.DB.Exec(query, tournament.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *TournamentRepository) GetById(id int) (*entity.Tournament, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id = $1", t.TableName)
	tournament := entity.Tournament{}
	err := t.DB.QueryRow(query, id).Scan(&tournament.ID, &tournament.Name)
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (t *TournamentRepository) AddTeam(tournamentID int, team entity.Team) (*entity.Team, error) {
	query := "INSERT INTO teams (tournament_id, name) VALUES ($1, $2) RETURNING id"
	err := t.DB.QueryRow(query, tournamentID, team.Name).Scan(&team.ID)
	if err != nil {
		return nil, err
	}
	team.TournamentID = tournamentID
	return &team, nil
}

func (t *TournamentRepository) GetTeams(tournamentID int) ([]entity.Team, error) {
	query := "SELECT id, tournament_id, name FROM teams WHERE tournament_id = $1"
	rows, err := t.DB.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []entity.Team
	for rows.Next() {
		team := entity.Team{}
		err := rows.Scan(&team.ID, &team.TournamentID, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return teams, nil
}
