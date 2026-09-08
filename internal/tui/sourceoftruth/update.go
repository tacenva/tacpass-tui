package sourceoftruth

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tacenva/database"
	"github.com/tacenva/tacpass-core/util/keyring"

	"github.com/tacenva/tacpass-tui/internal/api"
	"github.com/tacenva/tacpass-tui/internal/app"
	"github.com/tacenva/tacpass-tui/internal/config"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/detail"
	"github.com/tacenva/tacpass-tui/internal/tui/sourceoftruth/form"
)

func InitCmd() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// =========================
	// FORM
	// =========================

	if m.Form.Active {
		switch msg := msg.(type) {
		case form.SubmitMsg:
			keypair := keyring.KeyPair{
				PublicKey:  msg.PublicKey,
				PrivateKey: msg.PrivateKey,
			}

			var err error

			if m.Form.Editing {
				if len(m.SoTList) == 0 {
					return m, nil
				}

				selectedSoT := m.SoTList[m.Cursor]

				selectedSoT.Hostname = msg.Hostname
				selectedSoT.Address = msg.Address
				selectedSoT.KeyPair = keypair

				err = m.SoTService.Update(&selectedSoT)
			} else {
				_, err = m.SoTService.Create(
					msg.Address,
					keypair,
				)
			}

			if err != nil {
				return m, nil
			}

			m.Form = form.Model{}

			if err := m.Load(m.masterKey); err != nil {
				return m, nil
			}

			return m, nil

		case form.CancelMsg:
			m.Form = form.Model{}

			return m, nil
		}

		updated, cmd := m.Form.Update(msg)
		m.Form = updated

		return m, cmd
	}

	// =========================
	// DETAIL
	// =========================

	if m.Detail.Active {
		updated, cmd := m.Detail.Update(msg)
		m.Detail = updated

		return m, cmd
	}

	// =========================
	// LIST
	// =========================

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit

		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down":
			maxCursor := len(m.SoTList)

			if m.Cursor < maxCursor {
				m.Cursor++
			}

		case "e":
			if len(m.SoTList) == 0 {
				return m, nil
			}

			selectedSoT := m.SoTList[m.Cursor]

			m.Form = form.NewEdit(
				selectedSoT.Hostname,
				selectedSoT.Address,
				selectedSoT.KeyPair.PublicKey,
				selectedSoT.KeyPair.PrivateKey,
			)

		case "enter":
			if m.Cursor == len(m.SoTList) {
				m.Form = form.New()
				return m, nil
			}

			if len(m.SoTList) == 0 {
				return m, nil
			}

			selectedSoT := m.SoTList[m.Cursor]

			isRemote := selectedSoT.Address != "localhost"

			if isRemote {
				m.appDeps.Client.ConfigureTLS(
					selectedSoT.Address,
					api.TLSConfig{
						Fingerprint: selectedSoT.TLSFingerprint,

						OnFirstTrust: func(
							fingerprint string,
						) error {
							selectedSoT.TLSFingerprint = fingerprint
							return m.SoTService.Update(&selectedSoT)
						},
					},
				)
			}

			nodeDBDir := m.appDeps.Config.Path(
				config.NodeDirName,
				selectedSoT.ID,
				"vault",
			)

			m.appDeps.Client.SetToken(
				selectedSoT.AuthToken,
			)

			m.Detail = detail.New(
				m.appDeps,
				&app.Context{
					SelectedSoT: &selectedSoT,
					NodeDB:      database.New(nodeDBDir),
					IsRemote:    isRemote,
				},
				m.masterKey,
				m.coreServices.Vault,
				m.coreServices.Auth,
				m.coreServices.AccessControl,
			)
		}
	}

	return m, nil
}
