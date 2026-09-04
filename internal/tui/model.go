package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-core/user"
	SoTService "github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth"
	"gorm.io/gorm"
)

type Screen int

const (
	ScreenLogin Screen = iota
	ScreenSourceOfTruth
)

type Model struct {
	Width  int
	Height int

	Screen Screen

	sotService   *SoTService.Service
	Input        string
	ErrorMessage string

	SourceOfTruth sourceoftruth.Model
}

func New(sqliteDB *gorm.DB, tacenvaDB *database.DB) Model {
	permissionRepository := permission.NewRepository(sqliteDB)
	permissionService := permission.NewService(permissionRepository)
	userRepository := user.NewRepository(sqliteDB)
	userService := user.NewService(userRepository)

	authService := auth.NewService(userService, permissionService)
	soTService := SoTService.NewService(tacenvaDB, authService)

	return Model{
		Screen:        ScreenLogin,
		sotService:    soTService,
		SourceOfTruth: sourceoftruth.New(soTService),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
