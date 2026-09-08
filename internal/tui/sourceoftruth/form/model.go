package form

import "github.com/tacenva/tacpass-tui/internal/tui/state"

type Field int

const (
	FieldHostname Field = iota
	FieldAddress
	FieldPublicKey
	FieldPrivateKey
)

type Model struct {
	Active bool

	Width  int
	Height int

	Hostname   string
	Address    string
	PublicKey  string
	PrivateKey string

	EditField Field
	Editing   bool

	ActionState *state.Async
}

func New(actionState *state.Async) Model {
	return Model{
		Active:      true,
		EditField:   FieldHostname,
		ActionState: actionState,
	}
}

func NewEdit(
	hostname string,
	address string,
	publicKey string,
	privateKey string,
	actionState *state.Async,
) Model {
	return Model{
		Active:      true,
		Hostname:    hostname,
		Address:     address,
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
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
}

type CancelMsg struct{}
