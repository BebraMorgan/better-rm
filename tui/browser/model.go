package browser

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	entries     []os.DirEntry
	cursor      int
	path        string
	err         error
	width       int
	height      int
	selected    map[string]struct{}
	visualMode  bool
	visualStart int
}

func NewModel(startPath string) Model {
	if startPath == "" {
		startPath, _ = os.Getwd()
	}
	entries, err := readDirSorted(startPath)
	return Model{
		entries:  entries,
		cursor:   0,
		path:     startPath,
		err:      err,
		selected: make(map[string]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}
