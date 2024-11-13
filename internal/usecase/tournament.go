package usecase

import (
	"errors"
	"net/http"
	"time"
	"tournament/internal/entity"

	"math/rand"
)

type TournamentUseCase struct {
	TournamentRepository TournamentRepository
	GameRepository       GameRepository
}

func NewTournamentUsecase(tournamentRep TournamentRepository, gameRep GameRepository) *TournamentUseCase {
	return &TournamentUseCase{
		TournamentRepository: tournamentRep,
		GameRepository:       gameRep,
	}
}

type CreateTournamentRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateTournamentResponse struct {
	StatusCode int                `json:"status_code"`
	Tournament *entity.Tournament `json:"tournament"`
}

type DeleteTournamentResponse struct {
	StatusCode int                `json:"status_code"`
	Tournament *entity.Tournament `json:"tournament"`
}

type AddTeamRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddTeamResponse struct {
	StatusCode int          `json:"status_code"`
	Team       *entity.Team `json:"team"`
}

type TournamentResultResponse struct {
	StatusCode int         `json:"status_code"`
	Winner     entity.Team `json:"winner"`
}

func (t *TournamentUseCase) CreateTournament(req CreateTournamentRequest) (*CreateTournamentResponse, error) {

	tournament := entity.Tournament{
		Name: req.Name,
	}

	res, err := t.TournamentRepository.Create(tournament)

	if err != nil {
		return nil, err
	}

	return &CreateTournamentResponse{
		StatusCode: http.StatusOK,
		Tournament: res,
	}, nil
}

func (t *TournamentUseCase) DeleteTournament(tournamentID int) (*DeleteTournamentResponse, error) {
	tournament, err := t.TournamentRepository.GetById(tournamentID)

	if err != nil {
		return nil, err
	}

	err = t.TournamentRepository.Delete(*tournament)

	if err != nil {
		return nil, err
	}

	return &DeleteTournamentResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func (t *TournamentUseCase) AddTeam(tournamentID int, req AddTeamRequest) (*AddTeamResponse, error) {
	tournament, err := t.TournamentRepository.GetById(tournamentID)

	if err != nil {
		return nil, err
	}

	team := entity.Team{
		Name: req.Name,
	}

	res, err := t.TournamentRepository.AddTeam(tournament.ID, team)

	if err != nil {
		return nil, err
	}

	return &AddTeamResponse{
		StatusCode: http.StatusOK,
		Team:       res,
	}, nil
}

func (t *TournamentUseCase) GenerateDivisionSchedule(tournamentId int) error {
	teams, err := t.TournamentRepository.GetTeams(tournamentId)

	if err != nil {
		return err
	}

	//разделения на два дивизиона
	firstDivision, secondDivision, err := splitToDivisions(teams)

	if err != nil {
		return err
	}
	//генерация расписания для первого дивизиона
	for i := 0; i < len(firstDivision); i++ {
		for j := i + 1; j < len(firstDivision); j++ {
			_, err := t.GameRepository.Create(tournamentId, firstDivision[i], firstDivision[j], entity.GAME_TYPE_DIVISION_A)
			if err != nil {
				return err
			}
		}
	}

	//генерация расписания для второго дивизиона
	for i := 0; i < len(secondDivision); i++ {
		for j := i + 1; j < len(secondDivision); j++ {
			_, err := t.GameRepository.Create(tournamentId, secondDivision[i], secondDivision[j], entity.GAME_TYPE_DIVISION_B)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TournamentUseCase) GenerateDivisionResult(tournamentID int) error {
	err := t.GenerateResultByGameType(tournamentID, entity.GAME_TYPE_DIVISION_A)

	if err != nil {
		return err
	}

	err = t.GenerateResultByGameType(tournamentID, entity.GAME_TYPE_DIVISION_B)

	if err != nil {
		return err
	}

	return nil
}

func (t *TournamentUseCase) GeneratePlayoffStage1Schedule(tournamentId int) error {
	//генерация расписания для первой стадии плейофф

	//берем топ 4 команды с каждого дивизиона и формируем расписание для первой стадии плей офф
	firstDivisionWinnners, err := t.GameRepository.GetTop4WinnersByType(tournamentId, entity.GAME_TYPE_DIVISION_A)
	if err != nil {
		return err
	}
	secondDivisionWinners, err := t.GameRepository.GetTop4WinnersByType(tournamentId, entity.GAME_TYPE_DIVISION_B)
	if err != nil {
		return err
	}

	//лучшие играют с худшими с другого дивизиона
	_, err = t.GameRepository.Create(tournamentId, firstDivisionWinnners[0], secondDivisionWinners[3], entity.GAME_TYPE_PLAYOFF_STAGE_1)
	if err != nil {
		return err
	}
	_, err = t.GameRepository.Create(tournamentId, secondDivisionWinners[0], firstDivisionWinnners[3], entity.GAME_TYPE_PLAYOFF_STAGE_1)
	if err != nil {
		return err
	}
	_, err = t.GameRepository.Create(tournamentId, firstDivisionWinnners[1], secondDivisionWinners[2], entity.GAME_TYPE_PLAYOFF_STAGE_1)
	if err != nil {
		return err
	}
	_, err = t.GameRepository.Create(tournamentId, secondDivisionWinners[1], firstDivisionWinnners[2], entity.GAME_TYPE_PLAYOFF_STAGE_1)
	if err != nil {
		return err
	}
	return nil
}

func (t *TournamentUseCase) GenerateSemininalSchedule(tournamentID int) error {
	winners, err := t.GameRepository.GetWinnersByType(tournamentID, entity.GAME_TYPE_PLAYOFF_STAGE_1)
	if err != nil {
		return err
	}
	// Создаем пары для полуфинала
	_, err = t.GameRepository.Create(tournamentID, winners[0], winners[3], entity.GAME_TYPE_PLAYOFF_SEMIFINAL)
	if err != nil {
		return err
	}

	_, err = t.GameRepository.Create(tournamentID, winners[1], winners[2], entity.GAME_TYPE_PLAYOFF_SEMIFINAL)
	if err != nil {
		return err
	}

	return nil
}

func (t *TournamentUseCase) GenerateFinalSchedule(tournamentID int) error {
	winners, err := t.GameRepository.GetWinnersByType(tournamentID, entity.GAME_TYPE_PLAYOFF_SEMIFINAL)
	if err != nil {
		return err
	}
	_, err = t.GameRepository.Create(tournamentID, winners[0], winners[1], entity.GAME_TYPE_PLAYOFF_FINAL)
	if err != nil {
		return err
	}
	return nil
}

func (t *TournamentUseCase) GenerateFinalResult(tournamentID int) (*TournamentResultResponse, error) {
	err := t.GenerateResultByGameType(tournamentID, entity.GAME_TYPE_PLAYOFF_FINAL)
	if err != nil {
		return nil, err
	}
	winners, err := t.GameRepository.GetWinnersByType(tournamentID, entity.GAME_TYPE_PLAYOFF_FINAL)

	if err != nil {
		return nil, err
	}

	winner := winners[0]

	return &TournamentResultResponse{
		StatusCode: http.StatusOK,
		Winner:     winner,
	}, nil
}

func (t *TournamentUseCase) GenerateResultByGameType(tournamentID int, gameType int) error {
	games, err := t.GameRepository.GetByTypeGames(tournamentID, gameType)

	if err != nil {
		return err
	}

	for i := 0; i < len(games); i++ {
		winner := runGame(games[i].Team1ID, games[i].Team2ID)
		games[i].WinnerId = &winner
		//обновляем в базе инфу о победителе матча
		_, err := t.GameRepository.Update(games[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func splitToDivisions(teams []entity.Team) ([]entity.Team, []entity.Team, error) {
	if len(teams) != 16 {
		return nil, nil, errors.New("expected 16 teams")
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	rng.Shuffle(len(teams), func(i, j int) {
		teams[i], teams[j] = teams[j], teams[i]
	})

	firstDivision := teams[:8]
	secondDivision := teams[8:]

	return firstDivision, secondDivision, nil
}

func runGame(a int, b int) int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rng.Intn(2) == 0 {
		return a
	}
	return b
}
