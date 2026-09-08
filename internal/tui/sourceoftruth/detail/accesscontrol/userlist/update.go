package userlist

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func debugLog(format string, args ...any) {
	f, err := os.OpenFile(
		"debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return
	}
	defer f.Close()

	log.New(
		f,
		"[DEBUG] ",
		log.LstdFlags,
	).Printf(format, args...)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case UsersLoadedMsg:
		if m.ScreenState.Loading {
			if msg.Err != nil {
				m.ScreenState.Fail(msg.Err)
				return m, nil
			}

			m.Users = msg.Users
			m.normalizeCursor()
			m.ScreenState.Success()

			return m, nil
		}

		if m.ActionState.Loading {
			if msg.Err != nil {
				m.ActionState.Fail(msg.Err)
				return m, nil
			}

			m.Users = msg.Users
			m.normalizeCursor()
			m.ActionState.Success()

			return m, nil
		}

	case UserUpdatedMsg:
		if msg.Err != nil {
			m.ActionState.Fail(msg.Err)
			return m, nil
		}

		return m, m.load()

	case tea.KeyMsg:
		if m.ScreenState.Loading ||
			m.ActionState.Loading {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Users)-1 {
				m.Cursor++
			}

		case "a":
			if len(m.Users) == 0 {
				debugLog("approve: no users")
				return m, nil
			}

			debugLog(
				"approve key: cursor=%d userID=%s permissionID=%s",
				m.Cursor,
				m.Users[m.Cursor].ID,
				m.PermissionID,
			)

			m.ActionState.Start()

			return m, m.approveUser()

		case "r":
			if len(m.Users) == 0 {
				return m, nil
			}

			m.ActionState.Start()

			return m, m.revokeUser()

		case "esc":
			m.Focus = FocusNone
			m.Active = false
		}
	}

	return m, nil
}

func (m *Model) normalizeCursor() {
	if len(m.Users) == 0 {
		m.Cursor = 0
		return
	}

	if m.Cursor >= len(m.Users) {
		m.Cursor = len(m.Users) - 1
	}
}

func (m Model) approveUser() tea.Cmd {
	if len(m.Users) == 0 {
		debugLog("approveUser: no users")
		return nil
	}

	if m.Cursor < 0 || m.Cursor >= len(m.Users) {
		debugLog(
			"approveUser: invalid cursor=%d users=%d",
			m.Cursor,
			len(m.Users),
		)
		return nil
	}

	userID := m.Users[m.Cursor].ID

	debugLog(
		"approveUser: userID=%s permissionID=%s",
		userID,
		m.PermissionID,
	)

	return func() tea.Msg {
		debugLog(
			"approveUser cmd: calling service userID=%s",
			userID,
		)

		user, err := m.service.ApproveUser(userID)

		if err != nil {
			debugLog(
				"approveUser cmd: error=%v",
				err,
			)
		} else {
			debugLog(
				"approveUser cmd: success userID=%s",
				user.ID,
			)
		}

		return UserUpdatedMsg{
			User: user,
			Err:  err,
		}
	}
}

func (m Model) revokeUser() tea.Cmd {
	if len(m.Users) == 0 {
		return nil
	}

	if m.Cursor < 0 || m.Cursor >= len(m.Users) {
		return nil
	}

	userID := m.Users[m.Cursor].ID

	return func() tea.Msg {
		user, err := m.service.RevokeUser(userID)

		return UserUpdatedMsg{
			User: user,
			Err:  err,
		}
	}
}
