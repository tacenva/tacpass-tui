package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/database"
	coreapp "github.com/tacenva/tacpass-core/app"
	"github.com/tacenva/tacpass-core/config"
	"github.com/tacenva/tacpass-tui/internal/api"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadOrCreate()
	if err != nil {
		log.Fatal(err)
	}

	tacenvaDB := database.New(cfg.BaseDir)

	sqliteDB, err := OpenSQLite(
		cfg.Path(config.AppDBFileName),
	)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := sqliteDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := coreapp.Migrate(sqliteDB); err != nil {
		log.Fatal(err)
	}

	client := api.NewClient()

	appDeps := app.Deps{
		Config:   cfg,
		AppDB:    tacenvaDB,
		SqliteDB: sqliteDB,
		Client:   client,
	}

	services := coreapp.NewServices(sqliteDB, tacenvaDB)

	p := tea.NewProgram(
		tui.New(&appDeps, services),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func OpenSQLite(
	path string,
) (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open(path),
		&gorm.Config{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"open sqlite database: %w",
			err,
		)
	}

	return db, nil
}
