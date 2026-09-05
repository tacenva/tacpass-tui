package app

import (
	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/auth"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"gorm.io/gorm"
)

type DatabaseDeps struct {
	TacenvaDB *database.DB
	SqliteDB  *gorm.DB
}

type Context struct {
	SelectedSoT       *entity.SourceOfTruth
	AuthService       *auth.Service
	PermissionService *permission.Service
}
