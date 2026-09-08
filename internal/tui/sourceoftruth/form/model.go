package form

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
}

func New() Model {
	return Model{
		Active:    true,
		EditField: FieldHostname,
	}
}

func NewEdit(
	hostname string,
	address string,
	publicKey string,
	privateKey string,
) Model {
	return Model{
		Active:     true,
		Hostname:   hostname,
		Address:    address,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		EditField:  FieldHostname,
		Editing:    true,
	}
}

type SubmitMsg struct {
	Hostname   string
	Address    string
	PublicKey  string
	PrivateKey string
}

type CancelMsg struct{}
