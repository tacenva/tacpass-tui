package form

import (
	"net"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultPort = "9443"

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return CancelMsg{}
			}

		case "ctrl+s":
			return m, m.submit()

		case "tab":
			m.nextField()

		case "shift+tab":
			m.previousField()

		case "up":
			m.previousField()

		case "down":
			m.nextField()

		case "backspace":
			m.deleteLast()
		}

		if len(msg.Runes) > 0 {
			m.insert(string(msg.Runes))
		}
	}

	return m, nil
}

func (m *Model) nextField() {
	m.EditField++

	if m.EditField > FieldPrivateKey {
		m.EditField = FieldHostname
	}
}

func (m *Model) previousField() {
	if m.EditField == FieldHostname {
		m.EditField = FieldPrivateKey
		return
	}

	m.EditField--
}

func (m *Model) insert(value string) {
	switch m.EditField {
	case FieldHostname:
		m.Hostname += value

	case FieldAddress:
		m.Address += value

	case FieldPublicKey:
		m.PublicKey += value

	case FieldPrivateKey:
		m.PrivateKey += value
	}
}

func (m *Model) deleteLast() {
	switch m.EditField {
	case FieldHostname:
		m.Hostname = deleteLastRune(m.Hostname)

	case FieldAddress:
		m.Address = deleteLastRune(m.Address)

	case FieldPublicKey:
		m.PublicKey = deleteLastRune(m.PublicKey)

	case FieldPrivateKey:
		m.PrivateKey = deleteLastRune(m.PrivateKey)
	}
}

func deleteLastRune(value string) string {
	runes := []rune(value)

	if len(runes) == 0 {
		return value
	}

	return string(runes[:len(runes)-1])
}

func normalizeAddress(address string) string {
	address = strings.TrimSpace(address)

	if address == "" {
		return ""
	}

	if !strings.HasPrefix(address, "https://") &&
		!strings.HasPrefix(address, "http://") {
		address = "https://" + address
	}

	// Kalau sudah punya port, jangan ditambah lagi.
	host := strings.TrimPrefix(address, "https://")
	host = strings.TrimPrefix(host, "http://")

	if _, _, err := net.SplitHostPort(host); err != nil {
		address += ":" + defaultPort
	}

	return address
}

func (m Model) submit() tea.Cmd {
	return func() tea.Msg {
		return SubmitMsg{
			Hostname:   m.Hostname,
			Address:    normalizeAddress(m.Address),
			PublicKey:  m.PublicKey,
			PrivateKey: m.PrivateKey,
		}
	}
}
