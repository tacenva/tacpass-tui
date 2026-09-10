package app

import (
	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/config"
	"github.com/tacenva/tacpass-tui/internal/api"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"gorm.io/gorm"
)

type Deps struct {
	Config   *config.Config
	AppDB    *database.DB
	SqliteDB *gorm.DB
	Client   *api.Client
}

type Context struct {
	SelectedSoT *entity.SourceOfTruth
	NodeDB      *database.DB
	IsRemote    bool
}
