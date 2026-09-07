package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/config"
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

	appDBFilename := cfg.Path(config.AppDBFileName)
	if err != nil {
		return
	}
	sqliteDB, err := gorm.Open(
		sqlite.Open(appDBFilename),
		&gorm.Config{},
	)
	if err != nil {
		fmt.Println("failed to open sqlite:", err)
		os.Exit(1)
	}

	if err := sqliteDB.AutoMigrate(
		&entity.User{},
		&entity.Permission{},
		&entity.Vault{},
		&entity.VaultAccess{},
	); err != nil {
		fmt.Println("failed to migrate database:", err)
		os.Exit(1)
	}

	appDeps := app.Deps{
		Config:   cfg,
		AppDB:    tacenvaDB,
		SqliteDB: sqliteDB,
	}

	p := tea.NewProgram(
		tui.New(&appDeps),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
