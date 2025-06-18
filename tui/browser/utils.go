package browser

import (
	"brm/actions"
	"brm/trash"
	"fmt"
	"os"
	"path/filepath"
)

func (m *Model) isInTrash() bool {
	path, err := actions.GetTrashPath()
	m.err = fmt.Errorf("cannot open trash: %v", err)
	return m.path == path
}

func (m *Model) openTrash() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		m.err = fmt.Errorf("cannot determine home directory: %v", err)
		return
	}
	trashPath := filepath.Join(homeDir, ".trash")
	entries, err := readDirSorted(trashPath)
	if err != nil {
		m.err = fmt.Errorf("cannot open trash: %v", err)
		return
	}
	m.path = trashPath
	m.entries = entries
	m.cursor = 0
	m.err = nil
}

func removeTrashInfo(list []trash.TrashInfo, entry trash.TrashInfo) []trash.TrashInfo {
	for i, v := range list {
		if v.TrashName == entry.TrashName {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}
