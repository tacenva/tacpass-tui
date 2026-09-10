package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	coreApp "github.com/tacenva/tacpass-core/app"
	"github.com/tacenva/tacpass-tui/internal/operations/app"
	SoTService "github.com/tacenva/tacpass-tui/internal/operations/app/sourceoftruth"
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
	deps          *app.Deps
}

func New(
	deps *app.Deps,
	coreServices *coreApp.Services,
) Model {
	soTService := SoTService.NewService(deps, coreServices.Auth, coreServices.Permission)

	return Model{
		Screen:        ScreenLogin,
		sotService:    soTService,
		SourceOfTruth: sourceoftruth.New(deps, soTService, coreServices),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
