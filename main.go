package main

import (
	"log/slog"
	"os"
)

type application struct {
	db  *database
	log *slog.Logger
	cfg *config
}

func main() {
	// get the db path
	cfg := parseConfig()

	db, err := loadDB(cfg.dbPath)
	if err != nil {
		panic(err.Error())
	}

	// init g object
	app := &application{
		db:  db,
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		cfg: cfg,
	}

	// run server
	if err := app.serve(); err != nil {
		panic(err.Error())
	}
}
