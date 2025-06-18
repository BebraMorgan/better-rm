package browser

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Main(startPath string) error {
	model := NewModel(startPath)
	program := tea.NewProgram(model)
	return program.Start()
}
