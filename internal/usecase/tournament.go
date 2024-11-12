package usecase

import (
	"errors"
	"fmt"
	"sort"
	"time"
	"tournament/internal/entity"

	"math/rand"
)

const TEAM_COUNT = 16

type Tournament struct {
	Teams []entity.Team
}

func NewTournament(teams []entity.Team) (*Tournament, error) {
	if len(teams) != TEAM_COUNT {
		return nil, fmt.Errorf("team count must be %d", TEAM_COUNT)
	}

	return &Tournament{
		Teams: teams,
	}, nil
}

func (t *Tournament) Run() (*entity.Team, error) {
	firstDivision, secondDivision, err := splitToDivisions(t.Teams)

	if err != nil {
		return nil, err
	}

	winnersFirstDivisionTeams := runDivisionGames(firstDivision)
	winnersSecondDivisionTeams := runDivisionGames(secondDivision)

	fmt.Println(winnersFirstDivisionTeams, winnersSecondDivisionTeams)

	return nil, nil
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

func runGame(a entity.Team, b entity.Team) entity.Team {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if rng.Intn(2) == 0 {
		return a
	}
	return b
}

func runDivisionGames(teams []entity.Team) []entity.Team {
	results := make(map[entity.Team]int)
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ {
			winner := runGame(teams[i], teams[j])
			results[winner]++
		}
	}
	// сортируем teams по результатам матчей
	sort.Slice(teams, func(i, j int) bool {
		return results[teams[i]] > results[teams[j]]
	})
	// возвращаем топ 4 команды по победам
	return teams[:4]
}
