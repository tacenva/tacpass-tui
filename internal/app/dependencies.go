package app

import (
	"github.com/tacenva/database"
	"gorm.io/gorm"
)

type Dependencies struct {
	TacenvaDB *database.DB
	SqliteDB  *gorm.DB
}
