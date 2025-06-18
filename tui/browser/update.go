package browser

import (
	"brm/localization"
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) moveUp() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *Model) moveDown() {
	if m.cursor < len(m.entries)-1 {
		m.cursor++
	}
}

func (m *Model) openDir() {
	if len(m.entries) == 0 {
		return
	}
	selectedEntry := m.entries[m.cursor]
	if selectedEntry.IsDir() {
		newPath := filepath.Join(m.path, selectedEntry.Name())
		entries, err := readDirSorted(newPath)
		if err == nil {
			m.path = newPath
			m.entries = entries
			m.cursor = 0
			m.err = nil
		} else {
			m.err = err
		}
	}
}

func (m *Model) goBack() {
	parent := filepath.Dir(m.path)
	entries, err := readDirSorted(parent)
	if err == nil {
		m.path = parent
		m.entries = entries
		m.cursor = 0
		m.err = nil
	} else {
		m.err = err
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "v":
			if !m.visualMode {
				m.visualMode = true
				m.visualStart = m.cursor
			} else {
				m.visualMode = false
			}
		case "T":
			m.openTrash()
			m.visualMode = false
		case "up", "k":
			m.moveUp()
		case "down", "j":
			m.moveDown()
		case "enter", "right", "l":
			m.openDir()
		case "backspace", "h", "left":
			m.goBack()
		case "R":
			if m.isInTrash() {
				if m.visualMode {
					m.restoreVisualSelected()
					m.visualMode = false
				} else {
					m.restoreSelected()
				}
			} else {
				m.err = fmt.Errorf("%s", localization.GetMessage("restoration_only_in_trash"))
			}
		case "d", "delete":
			if m.visualMode {
				m.deleteVisualSelected()
				m.visualMode = false
			} else {
				m.deleteSelected()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

