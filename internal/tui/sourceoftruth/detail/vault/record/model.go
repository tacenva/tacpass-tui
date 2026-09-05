package vaultrecord

import (
	"time"

	"github.com/tacenva/tacpass-core/entity"
)

type Model struct {
	Active bool

	Width  int
	Height int

	SelectedVault *entity.Vault
	ShowPassword  bool

	Records []entity.VaultRecord
	Cursor  int
}

func New() Model {
	return Model{
		Records: []entity.VaultRecord{
			{
				ID:        "01K00000000000000000000001",
				Name:      "GitHub",
				Endpoint:  "https://github.com",
				Password:  "github-password",
				ExpiredAt: time.Now().AddDate(0, 1, 0),
			},
			{
				ID:        "01K00000000000000000000002",
				Name:      "GitLab",
				Endpoint:  "https://gitlab.com",
				Password:  "gitlab-password",
				ExpiredAt: time.Now().AddDate(0, 2, 0),
			},
			{
				ID:        "01K00000000000000000000003",
				Name:      "Production Server",
				Endpoint:  "https://prod.example.com",
				Password:  "production-password",
				ExpiredAt: time.Now().AddDate(0, 3, 0),
			},
			{
				ID:        "01K00000000000000000000004",
				Name:      "Database",
				Endpoint:  "postgres://db.example.com:5432",
				Password:  "database-password",
				ExpiredAt: time.Now().AddDate(0, 1, 15),
			},
		},
		Cursor: 0,
		Active: false,
	}
}

func (m Model) Selected() *entity.VaultRecord {
	if len(m.Records) == 0 {
		return nil
	}

	if m.Cursor < 0 || m.Cursor >= len(m.Records) {
		return nil
	}

	return &m.Records[m.Cursor]
}
