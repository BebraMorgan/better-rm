package browser

import (
	"brm/localization"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mattn/go-runewidth"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(str string) string {
	return ansiRegexp.ReplaceAllString(str, "")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func readDirSorted(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		iIsDir := entries[i].IsDir()
		jIsDir := entries[j].IsDir()
		if iIsDir && !jIsDir {
			return true
		}
		if !iIsDir && jIsDir {
			return false
		}
		return entries[i].Name() < entries[j].Name()
	})
	return entries, nil
}

func confirmDeletePrompt(count int) (bool, error) {
	var label string
	if count == 1 {
		label = localization.GetMessage("confirm_delete_file", "")
	} else {
		label = localization.GetMessage("confirm_delete_files", count)
	}
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false, fmt.Errorf("%s", localization.GetMessage("deletion_cancelled_by_user"))
	}
	answer := strings.ToLower(strings.TrimSpace(result))
	return answer == "y" || answer == "", nil
}

func (m *Model) isVisualSelected(i int) bool {
	if !m.visualMode {
		return false
	}
	start, end := m.visualStart, m.cursor
	if start > end {
		start, end = end, start
	}
	return i >= start && i <= end
}

func (m Model) renderHeader() string {
	headerContent := fmt.Sprintf(" %s", m.path)
	headerContentWidth := runewidth.StringWidth(headerContent)
	padding := max(0, m.width-headerContentWidth)
	paddedHeader := headerContent + strings.Repeat(" ", padding)
	return fmt.Sprintf("%s%s%s\n", FgWhiteBright, paddedHeader, Reset)
}

func (m Model) renderFooter() string {
	footerContent := "↑/↓ or j/k — move, Enter/l — open dir, Backspace/h — up, v — visual mode, T — trash, d/delete — delete, q — quit"
	if m.isInTrash() {
		footerContent += " | R — restore"
	}
	footerContentWidth := runewidth.StringWidth(footerContent)
	padding := max(0, m.width-footerContentWidth)
	paddedFooter := footerContent + strings.Repeat(" ", padding)
	return fmt.Sprintf("\n%s%s%s\n", "\033[30;42m", paddedFooter, Reset)
}

func (m Model) renderEntries() string {
	var s strings.Builder
	const fixedHeight = 30
	maxVisibleEntries := fixedHeight
	startIdx := 0
	if m.cursor >= maxVisibleEntries && len(m.entries) > maxVisibleEntries {
		startIdx = m.cursor - maxVisibleEntries + 1
	}
	if len(m.entries) < maxVisibleEntries {
		maxVisibleEntries = len(m.entries)
		startIdx = 0
	}
	const highlightStart = "\033[30;43m"
	const highlightEnd = "\033[0m"
	for i := startIdx; i < startIdx+maxVisibleEntries && i < len(m.entries); i++ {
		entry := m.entries[i]
		fullPath := filepath.Join(m.path, entry.Name())
		cursor := "  "
		selectionIndicator := " "
		lineContent := entry.Name()
		if entry.IsDir() {
			lineContent = fmt.Sprintf("%s%s%s", FgCyan+Bold, lineContent+"/", Reset)
		} else {
			lineContent = fmt.Sprintf("%s%s%s", FgWhite, lineContent, Reset)
		}
		if _, ok := m.selected[fullPath]; ok && !m.isVisualSelected(i) && !(m.cursor == i && !m.visualMode) {
			selectionIndicator = FgYellow + "*" + Reset
			lineContent = fmt.Sprintf("%s%s%s", FgGreen, entry.Name(), Reset)
			if entry.IsDir() {
				lineContent = fmt.Sprintf("%s%s%s", FgGreen+Bold, entry.Name()+"/", Reset)
			}
		}
		line := fmt.Sprintf("%s%s %s", cursor, selectionIndicator, lineContent)
		visibleLen := runewidth.StringWidth(stripAnsi(line))
		spacesToAdd := m.width - visibleLen
		if spacesToAdd < 0 {
			spacesToAdd = 0
		}
		lineWithPadding := line + strings.Repeat(" ", spacesToAdd)
		if m.isVisualSelected(i) || (m.cursor == i && !m.visualMode) {
			lineWithPadding = highlightStart + stripAnsi(lineWithPadding) + highlightEnd
		}
		s.WriteString(lineWithPadding + "\n")
	}
	emptyLinesCount := fixedHeight - (len(m.entries) - startIdx)
	width := m.width
	if width <= 0 {
		width = 80
	}
	for range emptyLinesCount {
		s.WriteString(strings.Repeat(" ", width) + "\n")
	}
	return s.String()
}

func (m Model) renderSelected() string {
	if len(m.selected) == 0 {
		return ""
	}
	selectedPaths := make([]string, 0, len(m.selected))
	for p := range m.selected {
		selectedPaths = append(selectedPaths, filepath.Base(p))
	}
	sort.Strings(selectedPaths)
	displaySelected := strings.Join(selectedPaths, ", ")
	maxWidth := m.width - len("Selected: ")
	if runewidth.StringWidth(displaySelected) > maxWidth {
		displaySelected = runewidth.Truncate(displaySelected, maxWidth-3, "...")
	}
	return fmt.Sprintf("Selected: %s%s%s\n", FgYellow, displaySelected, Reset)
}

func (m Model) renderError() string {
	if m.err == nil {
		return ""
	}
	return fmt.Sprintf("\n%sError: %v%s\n", FgRed+Bold, m.err, Reset)
}

func (m Model) View() string {
	var s strings.Builder
	s.WriteString(m.renderHeader())
	s.WriteString(m.renderEntries())
	s.WriteString(m.renderFooter())
	s.WriteString(m.renderSelected())
	s.WriteString(m.renderError())
	return s.String()
}

