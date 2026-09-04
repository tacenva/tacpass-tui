package detail

import "github.com/tacenva/tacpass-core/entity"

type Focus int

const (
	FocusSidebar Focus = iota
	FocusContent
)

type Model struct {
	Width  int
	Height int

	SidebarCursor int

	Vaults []entity.Vault
	Cursor int

	Focus  Focus
	Active bool

	AddingVault bool
	VaultName   string
}

func New() Model {
	return Model{
		Active:        true,
		SidebarCursor: 0,
		Cursor:        0,

		Vaults: []entity.Vault{
			{
				ID:   "vault-001",
				Name: "Personal",
			},
			{
				ID:   "vault-002",
				Name: "Work",
			},
			{
				ID:   "vault-003",
				Name: "Infrastructure",
			},
		},
	}
}
