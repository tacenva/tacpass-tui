package vaultrecord

import (
	"time"

	"github.com/tacenva/tacpass-core/entity"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/app/vault"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
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

	SelectedVaultAccess *entity.VaultAccess
	ShowPassword        bool

	Records []entity.VaultRecord
	Cursor  int

	Focus Focus

	// Form state
	EditField     EditField
	EditName      string
	EditEndpoint  string
	EditPassword  string
	EditExpiredAt string

	VaultServiceTUI *vault.Service

	context *app.Context

	ScreenState *state.Async
	ActionState *state.Async
}

func New(
	context *app.Context,
	vaultServiceTUI *vault.Service,
	ScreenState *state.Async,
	ActionState *state.Async,
) Model {
	return Model{
		Records: []entity.VaultRecord{},

		Cursor: 0,
		Focus:  FocusContent,
		Active: false,

		VaultServiceTUI: vaultServiceTUI,
		context:         context,

		ScreenState: ScreenState,
		ActionState: ActionState,
	}
}

func (m *Model) Load(
	selectedVaultAccess *entity.VaultAccess,
) error {
	if selectedVaultAccess == nil {
		return nil
	}

	m.SelectedVaultAccess = selectedVaultAccess
	m.ScreenState.Start()

	vaultRecords, err := m.VaultServiceTUI.ListRecords(
		selectedVaultAccess,
	)
	if err != nil {
		m.ScreenState.Fail(err)
		return err
	}

	m.Records = vaultRecords
	m.Cursor = 0
	m.ScreenState.Success()

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
	m.EditExpiredAt = ""

	if !record.ExpiredAt.IsZero() {
		m.EditExpiredAt = record.ExpiredAt.Format(
			"2006-01-02",
		)
	}

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

	expiredAt, err := parseExpiredAt(
		m.EditExpiredAt,
	)
	if err != nil {
		m.ActionState.Fail(err)
		return
	}

	record.Name = m.EditName
	record.Endpoint = m.EditEndpoint
	record.Password = m.EditPassword
	record.ExpiredAt = expiredAt

	m.ActionState.Start()

	_, err = m.VaultServiceTUI.UpdateRecord(
		m.SelectedVaultAccess,
		record,
	)
	if err != nil {
		m.ActionState.Fail(err)
		return
	}

	m.Records[m.Cursor] = *record

	m.clearForm()
	m.Focus = FocusContent
	m.ActionState.Success()
}

func (m *Model) saveNew() {
	if m.EditName == "" {
		return
	}

	expiredAt, err := parseExpiredAt(
		m.EditExpiredAt,
	)
	if err != nil {
		m.ActionState.Fail(err)
		return
	}

	record := entity.VaultRecord{
		Name:      m.EditName,
		Endpoint:  m.EditEndpoint,
		Password:  m.EditPassword,
		ExpiredAt: expiredAt,
	}

	m.ActionState.Start()

	newRecord, err := m.VaultServiceTUI.AppendRecord(
		m.SelectedVaultAccess,
		&record,
	)
	if err != nil {
		m.ActionState.Fail(err)
		return
	}

	m.Records = append(
		m.Records,
		*newRecord,
	)

	m.Cursor = len(m.Records) - 1

	m.clearForm()
	m.Focus = FocusContent
	m.ActionState.Success()
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

	return time.Parse(
		"2006-01-02",
		value,
	)
}

func (m *Model) deleteSelected() {
	record := m.Selected()
	if record == nil {
		return
	}

	m.ActionState.Start()

	err := m.VaultServiceTUI.DeleteRecord(
		m.SelectedVaultAccess,
		record,
	)
	if err != nil {
		m.ActionState.Fail(err)
		return
	}

	m.Records = append(
		m.Records[:m.Cursor],
		m.Records[m.Cursor+1:]...,
	)

	if m.Cursor >= len(m.Records) {
		m.Cursor = len(m.Records) - 1
	}

	if m.Cursor < 0 {
		m.Cursor = 0
	}

	m.ShowPassword = false
	m.ActionState.Success()
}
