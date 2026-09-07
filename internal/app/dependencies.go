package app

import (
	"github.com/tacenva/database"
	coreEntity "github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/config"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"gorm.io/gorm"
)

type Deps struct {
	Config   *config.Config
	AppDB    *database.DB
	SqliteDB *gorm.DB
}

type Context struct {
	SelectedSoT *entity.SourceOfTruth
	NodeDB      *database.DB
	AuthUser    *coreEntity.User
}
