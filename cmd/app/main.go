package main

import (
	"database/sql"
	"log"
	"os"
	"tournament/internal/handler"
	"tournament/internal/repository/pgsql"
	"tournament/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db := InitDB()
	defer db.Close()

	if err := InitMigrations(db); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	router := gin.Default()

	tournamentRepository := pgsql.NewTournamentRepository(db)
	gameRepository := pgsql.NewGameRepository(db)
	tournamentUsecase := usecase.NewTournamentUsecase(tournamentRepository, gameRepository)

	tournamentHandler := handler.NewTournamentHandler(tournamentUsecase)

	router.POST("/tournaments", tournamentHandler.CreateTournament)
	router.POST("/tournaments/:id", tournamentHandler.DeleteTournament)
	router.POST("/tournaments/:id/teams", tournamentHandler.AddTeam)
	router.GET("/tournaments/:id/run", tournamentHandler.RunTournament)

	router.Run()
}

func InitDB() *sql.DB {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	dsn := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	return db
}

func InitMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/pgsql",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Ошибка при создании экземпляра мигратора: %v", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}

	log.Println("Миграции успешно применены!")
	return nil
}
