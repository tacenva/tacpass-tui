package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-core/user"
	"github.com/tacenva/tacpass-tui/internal/app"
	SoTService "github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth"
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

func New(deps *app.DatabaseDeps) Model {
	permissionRepository := permission.NewRepository(deps.SqliteDB)
	permissionService := permission.NewService(permissionRepository)
	userRepository := user.NewRepository(deps.SqliteDB)
	userService := user.NewService(userRepository)

	authService := auth.NewService(userService, permissionService)
	soTService := SoTService.NewService(deps.TacenvaDB, authService)

	return Model{
		Screen:        ScreenLogin,
		sotService:    soTService,
		SourceOfTruth: sourceoftruth.New(deps, soTService, authService, permissionService),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
