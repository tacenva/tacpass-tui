package sourceoftruth

import (
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail"
)

type Model struct {
	Width  int
	Height int

	Detail     detail.Model
	SoTService *sourceoftruth.Service
	SoTList    []entity.SourceOfTruth

	Cursor int

	appDeps           *app.Deps
	authService       *auth.Service
	permissionService *permission.Service
}

func New(appDeps *app.Deps, SoTService *sourceoftruth.Service, authService *auth.Service, permissionService *permission.Service) Model {
	return Model{
		appDeps:           appDeps,
		authService:       authService,
		permissionService: permissionService,
		SoTService:        SoTService,
		Cursor:            0,
	}
}

func (m *Model) Load() error {
	soTList, err := m.SoTService.List()
	if err != nil {
		return err
	}

	m.SoTList = soTList
	return nil
}
