package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/database"
	coreapp "github.com/tacenva/tacpass-core/app"
	"github.com/tacenva/tacpass-core/config"
	"github.com/tacenva/tacpass-tui/internal/operations/api"
	"github.com/tacenva/tacpass-tui/internal/operations/app"
	"github.com/tacenva/tacpass-tui/internal/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Run(dev bool) error {
	cfg, err := config.LoadOrCreate(dev)
	if err != nil {
		return err
	}

	tacenvaDB := database.New(
		cfg.BaseDir,
	)

	sqliteDB, err := OpenSQLite(
		cfg.Path(config.AppDBFileName),
	)
	if err != nil {
		return err
	}

	sqlDB, err := sqliteDB.DB()
	if err != nil {
		return fmt.Errorf(
			"get sqlite database: %w",
			err,
		)
	}
	defer sqlDB.Close()

	if err := coreapp.Migrate(sqliteDB); err != nil {
		return err
	}

	client := api.NewClient()

	deps := app.Deps{
		Config:   cfg,
		AppDB:    tacenvaDB,
		SqliteDB: sqliteDB,
		Client:   client,
	}

	services := coreapp.NewServices(
		sqliteDB,
		tacenvaDB,
	)

	p := tea.NewProgram(
		tui.New(&deps, services),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
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
