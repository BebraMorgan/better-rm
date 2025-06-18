package browser

import (
	"brm/actions"
	"brm/trash"
	"fmt"
	"os"
	"path/filepath"

	"brm/localization"
)

func (m *Model) restoreSelected() {
	if !m.isInTrash() {
		m.err = fmt.Errorf("%s", localization.GetMessage("restoration_only_in_trash"))
		return
	}
	if len(m.selected) == 0 {
		m.err = fmt.Errorf("%s", localization.GetMessage("no_files_selected"))
		return
	}
	err := actions.Restore()
	if err != nil {
		m.err = fmt.Errorf("%s", localization.GetMessage("error_restoring_file", "files", err))
		return
	}
	m.selected = make(map[string]struct{})
	entries, err := readDirSorted(m.path)
	if err != nil {
		m.err = fmt.Errorf("%s", localization.GetMessage("failed_to_refresh_directory", err))
		return
	}
	m.entries = entries
	m.err = nil
}

func (m *Model) restoreVisualSelected() {
	if !m.isInTrash() {
		m.err = fmt.Errorf("%s", localization.GetMessage("restoration_only_in_trash"))
		return
	}
	if !m.visualMode {
		return
	}
	start, end := m.visualStart, m.cursor
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end && i < len(m.entries); i++ {
		entry := m.entries[i]
		fullPath := filepath.Join(m.path, entry.Name())
		trashInfoPath, err := trash.GetTrashInfoPath()
		if err != nil {
			m.err = fmt.Errorf("%s", localization.GetMessage("unable_to_get_trash_info_path", err))
			return
		}
		entries, err := trash.LoadTrashInfo(trashInfoPath)
		if err != nil {
			m.err = fmt.Errorf("%s", localization.GetMessage("unable_to_load_trash_info", err))
			return
		}
		found := false
		for _, info := range entries {
			if info.TrashName == entry.Name() {
				restoreErr := os.Rename(fullPath, info.OriginalPath)
				if restoreErr != nil {
					m.err = fmt.Errorf("%s", localization.GetMessage("error_restoring_file", fullPath, restoreErr))
					return
				}
				err := trash.SaveTrashInfo(trashInfoPath, removeTrashInfo(entries, info))
				if err != nil {
					m.err = fmt.Errorf("%s", localization.GetMessage("unable_to_update_trash_info", err))
					return
				}
				found = true
				break
			}
		}
		if !found {
			m.err = fmt.Errorf("%s", localization.GetMessage("could_not_find_original_path_for", entry.Name()))
			return
		}
	}
	entries, err := readDirSorted(m.path)
	if err != nil {
		m.err = fmt.Errorf("%s", localization.GetMessage("failed_to_refresh_directory", err))
		return
	}
	m.entries = entries
	m.visualMode = false
	m.err = nil
}

func (m *Model) deleteSelected() {
	if len(m.selected) == 0 && !m.visualMode && m.cursor < len(m.entries) {
		entry := m.entries[m.cursor]
		fullPath := filepath.Join(m.path, entry.Name())
		confirmed, err := confirmDeletePrompt(1)
		if err != nil || !confirmed {
			m.err = fmt.Errorf("%s", localization.GetMessage("deletion_cancelled_by_user"))
			return
		}
		if m.isInTrash() {
			if entry.IsDir() {
				err = os.RemoveAll(fullPath)
			} else {
				err = os.Remove(fullPath)
			}
			if err != nil {
				m.err = fmt.Errorf("%s", localization.GetMessage("error_deleting_file", fullPath, err))
				return
			}
		} else {
			err = actions.SaveDelete(fullPath)
			if err != nil {
				m.err = fmt.Errorf("%s", localization.GetMessage("error_deleting_file", fullPath, err))
				return
			}
		}
	} else if len(m.selected) > 0 {
		confirmed, err := confirmDeletePrompt(len(m.selected))
		if err != nil || !confirmed {
			m.err = fmt.Errorf("%s", localization.GetMessage("deletion_cancelled_by_user"))
			return
		}
		for path := range m.selected {
			if m.isInTrash() {
				err = os.RemoveAll(path)
			} else {
				err = actions.SaveDelete(path)
			}
			if err != nil {
				m.err = fmt.Errorf("%s", localization.GetMessage("error_deleting_file", path, err))
				return
			}
			delete(m.selected, path)
		}
	} else if m.visualMode {
		m.deleteVisualSelected()
		return
	}
	entries, err := readDirSorted(m.path)
	if err != nil {
		m.err = fmt.Errorf("%s", localization.GetMessage("error_refreshing_directory", err))
		return
	}
	if m.cursor > 0 && len(entries) <= m.cursor {
		m.cursor--
	}
	m.entries = entries
	m.err = nil
}

func (m *Model) deleteVisualSelected() {
	if !m.visualMode {
		return
	}
	start, end := m.visualStart, m.cursor
	if start > end {
		start, end = end, start
	}
	var pathsToDelete []string
	for i := start; i <= end && i < len(m.entries); i++ {
		fullPath := filepath.Join(m.path, m.entries[i].Name())
		pathsToDelete = append(pathsToDelete, fullPath)
	}
	confirmed, err := confirmDeletePrompt(len(pathsToDelete))
	if err != nil || !confirmed {
		m.err = fmt.Errorf("%s", localization.GetMessage("deletion_cancelled_by_user"))
		return
	}
	inTrash := m.isInTrash()
	for _, path := range pathsToDelete {
		if inTrash {
			err = os.RemoveAll(path)
			if err != nil {
				m.err = fmt.Errorf("%s", localization.GetMessage("error_deleting_file", path, err))
				return
			}
		} else {
			err = actions.SaveDelete(path)
			if err != nil {
				m.err = fmt.Errorf("%s", localization.GetMessage("error_deleting_file", path, err))
				return
			}
		}
	}
	entries, err := readDirSorted(m.path)
	if err != nil {
		m.err = fmt.Errorf("%s", localization.GetMessage("error_refreshing_directory", err))
		return
	}
	m.entries = entries
	m.cursor = 0
	m.err = nil
}
