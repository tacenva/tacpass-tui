package accesscontrol

import (
	"github.com/tacenva/tacpass-core/accesscontrol"
	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/permission"
	"github.com/tacenva/tacpass-core/user"

	"github.com/tacenva/tacpass-tui/internal/app"
	accessControlApp "github.com/tacenva/tacpass-tui/internal/app/accesscontrol"
)

type Focus int

const (
	FocusNone Focus = iota
	FocusContent
)

type Page int

const (
	PageList Page = iota
	PageUsers
	PageCreate
)

type Model struct {
	Active bool

	Focus Focus
	Page  Page

	Service *accessControlApp.Service

	Permissions []entity.Permission
	Users       []entity.User

	Cursor     int
	UserCursor int

	SelectedPermission *entity.Permission

	CreatePrivilege entity.Privilege

	Error   error
	Message string
}

func NewModel(
	deps *app.Deps,
	ctx *app.Context,
) Model {
	var address string

	if ctx.SelectedSoT != nil {
		address = ctx.SelectedSoT.Address
	}

	userRepository := user.NewRepository(deps.SqliteDB)
	userService := user.NewService(userRepository)
	permissionRepository := permission.NewRepository(deps.SqliteDB)
	permissionService := permission.NewService(permissionRepository)
	accesscontrolService := accesscontrol.NewService(userService, permissionService)

	return Model{
		Active: true,
		Focus:  FocusContent,
		Page:   PageList,

		Service: accessControlApp.NewService(
			deps.Client,
			address,
			accesscontrolService,
			ctx.IsRemote,
		),

		CreatePrivilege: entity.PrivilegeRead,
	}
}
