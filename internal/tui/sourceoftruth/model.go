package sourceoftruth

import (
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
}

func New(SoTService *sourceoftruth.Service) Model {
	return Model{
		SoTService: SoTService,
		Cursor:     0,
	}
}

func (m *Model) Load() error {
	SoTList, err := m.SoTService.List()
	if err != nil {
		return err
	}

	m.SoTList = SoTList

	return nil
}
