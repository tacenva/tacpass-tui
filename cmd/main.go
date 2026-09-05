package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	devDir := "dev/"

	if err := os.MkdirAll(devDir, 0700); err != nil {
		fmt.Println("failed to create dev directory:", err)
		os.Exit(1)
	}

	tacenvaDB := database.New(devDir)

	sqliteDB, err := gorm.Open(
		sqlite.Open(devDir+"tacenva.db"),
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

	appDeps := app.Dependencies{
		TacenvaDB: tacenvaDB,
		SqliteDB:  sqliteDB,
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
