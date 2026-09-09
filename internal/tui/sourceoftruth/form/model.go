package form

import (
	"github.com/tacenva/tacpass-tui/internal/entity"
	"github.com/tacenva/tacpass-tui/internal/tui/state"
)

type Field int

const (
	FieldHostname Field = iota
	FieldAddress
	FieldPublicKey
	FieldPrivateKey
	FieldSyncMode
)

type Model struct {
	Active bool

	Width  int
	Height int

	Hostname   string
	Address    string
	PublicKey  string
	PrivateKey string
	SyncMode   entity.SyncMode

	EditField Field
	Editing   bool

	ActionState *state.Async
}

func New(actionState *state.Async) Model {
	return Model{
		Active:      true,
		EditField:   FieldHostname,
		SyncMode:    entity.SyncModeAuto,
		ActionState: actionState,
	}
}

func NewEdit(
	hostname string,
	address string,
	publicKey string,
	privateKey string,
	syncMode entity.SyncMode,
	actionState *state.Async,
) Model {
	return Model{
		Active:      true,
		Hostname:    hostname,
		Address:     address,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		SyncMode:    syncMode,
		EditField:   FieldHostname,
		Editing:     true,
		ActionState: actionState,
	}
}

type SubmitMsg struct {
	Hostname   string
	Address    string
	PublicKey  string
	PrivateKey string
	SyncMode   entity.SyncMode
}

type CancelMsg struct{}
