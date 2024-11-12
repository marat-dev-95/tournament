package main

import (
	"log"
	"tournament/internal/entity"
	"tournament/internal/usecase"
)

func main() {
	teams := []entity.Team{
		{Name: "Team spirit"},
		{Name: "Team secret"},
		{Name: "Tundra"},
		{Name: "Virtus pro"},
		{Name: "Liquid"},
		{Name: "Heroic"},
		{Name: "NaVi"},
		{Name: "LGD"},
		{Name: "Falcons"},
		{Name: "Team Waska"},
		{Name: "Xtreme Gaming"},
		{Name: "Azure ray"},
		{Name: "Nouns"},
		{Name: "Avulus"},
		{Name: "Nigma"},
		{Name: "Gaming Gladiators"},
	}

	tournament, err := usecase.NewTournament(teams)

	if err != nil {
		log.Fatal(err.Error())
	}

	winner, err := tournament.Run()

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%s is the winner of the tournament", winner.Name)

}
