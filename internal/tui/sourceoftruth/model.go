package sourceoftruth

import (
	"encoding/json"

	coreApp "github.com/tacenva/tacpass-core/app"
	"github.com/tacenva/tacpass-tui/internal/operations/app"
	"github.com/tacenva/tacpass-tui/internal/operations/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/operations/entity"
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

	for _, sot := range soTList {
		_, err := json.Marshal(sot)
		if err != nil {
			continue
		}
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
