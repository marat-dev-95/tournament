package handler

import (
	"net/http"
	"strconv"
	"tournament/internal/entity"
	"tournament/internal/usecase"

	"github.com/gin-gonic/gin"
)

type TournamentHandler struct {
	TournamentUsecase *usecase.TournamentUseCase
}

func NewTournamentHandler(tournamentUsecase *usecase.TournamentUseCase) *TournamentHandler {
	return &TournamentHandler{
		TournamentUsecase: tournamentUsecase,
	}
}

func (t *TournamentHandler) CreateTournament(c *gin.Context) {
	var req usecase.CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: Converter(err), StatusCode: http.StatusBadRequest})
		return
	}

	res, err := t.TournamentUsecase.CreateTournament(req)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (t *TournamentHandler) RunTournament(c *gin.Context) {
	tournamentIDStr := c.Param("id")

	tournamentID, err := strconv.Atoi(tournamentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": "Invalid tournament_id"},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// Делим команды на 2 дивизиона и формируем расписание
	err = t.TournamentUsecase.GenerateDivisionSchedule(tournamentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	// генерация результатов матчей в дивизионах
	err = t.TournamentUsecase.GenerateDivisionResult(tournamentID)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	// формирование расписания первой стадии плейофф по результатам матчей в дивизионах
	err = t.TournamentUsecase.GeneratePlayoffStage1Schedule(tournamentID)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	// генерация результатов первой стадии плейофф
	err = t.TournamentUsecase.GenerateResultByGameType(tournamentID, entity.GAME_TYPE_PLAYOFF_STAGE_1)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// формирование расписания полуфинала по результатам первой стадии плейофф
	err = t.TournamentUsecase.GenerateSemininalSchedule(tournamentID)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// генерация результатов полуфинала
	err = t.TournamentUsecase.GenerateResultByGameType(tournamentID, entity.GAME_TYPE_PLAYOFF_SEMIFINAL)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// формирование	расписания финала
	err = t.TournamentUsecase.GenerateFinalSchedule(tournamentID)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// генерация результата финала

	res, err := t.TournamentUsecase.GenerateFinalResult(tournamentID)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, res)

}

func (t *TournamentHandler) AddTeam(c *gin.Context) {
	var req usecase.AddTeamRequest

	tournamentIDStr := c.Param("id")

	tournamentID, err := strconv.Atoi(tournamentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": "Invalid tournament_id"},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: Converter(err), StatusCode: http.StatusBadRequest})
		return
	}

	res, err := t.TournamentUsecase.AddTeam(tournamentID, req)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (t *TournamentHandler) DeleteTournament(c *gin.Context) {
	tournamentIDStr := c.Param("id")

	tournamentID, err := strconv.Atoi(tournamentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": "Invalid tournament_id"},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	res, err := t.TournamentUsecase.DeleteTournament(tournamentID)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors:     map[string]string{"message:": err.Error()},
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
