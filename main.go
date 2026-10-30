package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	db, err := openDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	startRepl(newConfig(&plantStore{db: db}))
	return nil
}

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "app.db")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	goose.SetBaseFS(migrationFiles)
	if err := goose.SetDialect("sqlite3"); err != nil {
		db.Close()
		return nil, fmt.Errorf("configure migrations: %w", err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

func newConfig(store *plantStore) *config {
	return &config{
		store: store,
		commands: map[string]cliCommand{
			"help": {
				name:        "help",
				description: "Display a help message",
				callback:    commandHelp,
			},
			"search": {
				name:        "search",
				description: "Search for plants",
				callback:    commandSearch,
			},
			"add": {
				name:        "add",
				description: "adding plant to your list of plants",
				callback:    commandAddPlant,
			},
			"list": {
				name:        "list",
				description: "List plants saved in the SQL store",
				callback:    commandListPlants,
			},
			"stations": {
				name:        "stations",
				description: "List SMHI stations or show latest temperature and daily rainfall for a station",
				callback:    commandListWeatherStations,
			},
		},
	}
}
