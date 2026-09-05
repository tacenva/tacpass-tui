package vaultrecord

import (
	"time"

	"github.com/tacenva/tacpass-core/entity"
)

type Focus int

const (
	FocusContent Focus = iota
	FocusEdit
	FocusNew
	Unfocus
)

type EditField int

const (
	FieldName EditField = iota
	FieldEndpoint
	FieldPassword
	FieldExpiredAt
)

type Model struct {
	Active bool

	Width  int
	Height int

	SelectedVault *entity.Vault
	ShowPassword  bool

	Records []entity.VaultRecord
	Cursor  int

	Focus Focus

	// Form state
	EditField     EditField
	EditName      string
	EditEndpoint  string
	EditPassword  string
	EditExpiredAt string
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
		Focus:  FocusContent,
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

func (m *Model) startEdit() {
	record := m.Selected()
	if record == nil {
		return
	}

	m.EditField = FieldName

	m.EditName = record.Name
	m.EditEndpoint = record.Endpoint
	m.EditPassword = record.Password
	m.EditExpiredAt = record.ExpiredAt.Format("2006-01-02 15:04")

	m.ShowPassword = false
	m.Focus = FocusEdit
}

func (m *Model) startNew() {
	m.EditField = FieldName

	m.EditName = ""
	m.EditEndpoint = ""
	m.EditPassword = ""
	m.EditExpiredAt = ""

	m.ShowPassword = false
	m.Focus = FocusNew
}

func (m *Model) cancelForm() {
	m.EditName = ""
	m.EditEndpoint = ""
	m.EditPassword = ""
	m.EditExpiredAt = ""

	m.ShowPassword = false
	m.Focus = FocusContent
}

func (m *Model) saveEdit() {
	record := m.Selected()
	if record == nil {
		return
	}

	record.Name = m.EditName
	record.Endpoint = m.EditEndpoint
	record.Password = m.EditPassword

	if expiredAt, err := parseExpiredAt(m.EditExpiredAt); err == nil {
		record.ExpiredAt = expiredAt
	}

	m.clearForm()
	m.Focus = FocusContent
}

func (m *Model) saveNew() {
	if m.EditName == "" {
		return
	}

	expiredAt, err := parseExpiredAt(m.EditExpiredAt)
	if err != nil {
		return
	}

	record := entity.VaultRecord{
		ID:        newID(),
		Name:      m.EditName,
		Endpoint:  m.EditEndpoint,
		Password:  m.EditPassword,
		ExpiredAt: expiredAt,
	}

	m.Records = append(m.Records, record)
	m.Cursor = len(m.Records) - 1

	m.clearForm()
	m.Focus = FocusContent
}

func (m *Model) clearForm() {
	m.EditName = ""
	m.EditEndpoint = ""
	m.EditPassword = ""
	m.EditExpiredAt = ""
	m.ShowPassword = false
}

func parseExpiredAt(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	return time.Parse("2006-01-02 15:04", value)
}

func newID() string {
	return time.Now().Format("20060102150405.000000000")
}
