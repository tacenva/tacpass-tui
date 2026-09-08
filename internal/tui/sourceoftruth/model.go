package sourceoftruth

import (
	coreApp "github.com/tacenva/tacpass-core/app"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/form"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
)

type Model struct {
	Width  int
	Height int

	Detail detail.Model
	Form   form.Model

	SoTService *sourceoftruth.Service
	SoTList    []entity.SourceOfTruth

	Cursor int

	appDeps   *app.Deps
	masterKey string

	coreServices *coreApp.Services

	ScreenState *state.Async
	ActionState *state.Async
}

func New(
	appDeps *app.Deps,
	SoTService *sourceoftruth.Service,
	coreServices *coreApp.Services,
) Model {
	ScreenState := &state.Async{}
	ActionState := &state.Async{}

	return Model{
		appDeps:      appDeps,
		SoTService:   SoTService,
		coreServices: coreServices,
		Cursor:       0,
		Form:         form.Model{},
		ScreenState:  ScreenState,
		ActionState:  ActionState,
	}
}

func (m *Model) Load(masterKey string) error {
	soTList, err := m.SoTService.List()
	if err != nil {
		return err
	}

	m.SoTList = soTList

	if len(m.SoTList) == 0 {
		m.Cursor = 0
		return nil
	}

	if m.Cursor >= len(m.SoTList) {
		m.Cursor = len(m.SoTList) - 1
	}

	m.masterKey = masterKey

	return nil
}
