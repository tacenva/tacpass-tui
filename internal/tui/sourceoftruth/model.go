package sourceoftruth

import (
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/permission"

	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/form"
)

type Model struct {
	Width  int
	Height int

	Detail detail.Model
	Form   form.Model

	SoTService *sourceoftruth.Service
	SoTList    []entity.SourceOfTruth

	Cursor int

	appDeps           *app.Deps
	authService       *auth.Service
	permissionService *permission.Service

	masterKey string
}

func New(
	appDeps *app.Deps,
	SoTService *sourceoftruth.Service,
	authService *auth.Service,
	permissionService *permission.Service,
) Model {
	return Model{
		appDeps:           appDeps,
		authService:       authService,
		permissionService: permissionService,
		SoTService:        SoTService,
		Cursor:            0,
		Form:              form.Model{},
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
