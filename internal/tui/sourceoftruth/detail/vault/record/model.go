package vaultrecord

import (
	"time"

	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-core/util/keyring"
	vaultCore "github.com/tacenva/tacpass-core/vault"
	"github.com/tacenva/tacpass-tui/internal/app"
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

	VaultService *vaultCore.Service

	context *app.Context
}

func New(context *app.Context, VaultService *vaultCore.Service) Model {
	return Model{
		Records: []entity.VaultRecord{},
		Cursor:  0,
		Focus:   FocusContent,
		Active:  false,

		VaultService: VaultService,
		context:      context,
	}
}

func (m *Model) Load(selectedVault *entity.Vault) error {
	vaultRecords, err := m.VaultService.DecryptedRecordList(
		m.context.AuthUser,
		selectedVault.ID,
		(*keyring.KeyPair)(&m.context.SelectedSoT.KeyPair),
	)
	if err != nil {
		return err
	}

	m.SelectedVault = selectedVault
	m.Records = vaultRecords

	return nil
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
	// m.EditExpiredAt = record.ExpiredAt.Format("2006-01-02 15:04")

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

	// if expiredAt, err := parseExpiredAt(m.EditExpiredAt); err == nil {
	// 	record.ExpiredAt = expiredAt
	// }

	_, err := m.VaultService.UpdateRecord(m.context.AuthUser, m.SelectedVault.ID, record, &m.context.SelectedSoT.KeyPair)
	if err != nil {
		panic(err)
	}

	m.clearForm()
	m.Focus = FocusContent
}

func (m *Model) saveNew() {
	if m.EditName == "" {
		return
	}

	// expiredAt, err := parseExpiredAt(m.EditExpiredAt)
	// if err != nil {
	// 	return
	// }
	if m.VaultService == nil {
		panic("VaultService is nil")
	}

	if m.context == nil {
		panic("context is nil")
	}

	if m.context.AuthUser == nil {
		panic("AuthUser is nil")
	}

	if m.SelectedVault == nil {
		panic("SelectedVault is nil")
	}
	record := entity.VaultRecord{
		Name:     m.EditName,
		Endpoint: m.EditEndpoint,
		Password: m.EditPassword,
		// ExpiredAt: expiredAt,
	}

	newRecord, err := m.VaultService.AppendRecord(
		m.context.AuthUser,
		m.SelectedVault.ID,
		&record,
		&m.context.SelectedSoT.KeyPair,
	)
	if err != nil {
		panic(err)
	}

	m.Records = append(m.Records, *newRecord)
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
