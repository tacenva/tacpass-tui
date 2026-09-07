package app

import (
	"github.com/tacenva/database"
	coreEntity "github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/entity"
	"gorm.io/gorm"
)

type DatabaseDeps struct {
	TacenvaDB *database.DB
	SqliteDB  *gorm.DB
}

type Context struct {
	SelectedSoT *entity.SourceOfTruth
	AuthUser    *coreEntity.User
	// AuthService       *auth.Service
	// PermissionService *permission.Service
}
