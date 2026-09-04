package sourceoftruth

import (
	"github.com/tacenva/tacpass-tui/internal/app/sourceoftruth"
	"github.com/tacenva/tacpass-tui/internal/entity"
)

type Model struct {
	Width  int
	Height int

	SoTService *sourceoftruth.Service
	SoTList    []entity.SourceOfTruth

	Cursor int
}

func New(SoTService *sourceoftruth.Service) Model {
	return Model{
		SoTService: SoTService,
	}
}
