package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type ImportModel struct {
	list     list.Model
	FilePath string
	dir      string
	done     bool
	err      error
}

type fileItem struct {
	label string
	path  string
}

func (f fileItem) Title() string       { return f.label }
func (f fileItem) Description() string { return "" }
func (f fileItem) FilterValue() string { return f.label }

func RunImportExcelUI() (*ImportModel, error) {
	startDir, _ := os.Getwd()
	m := NewImportModel(startDir)
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}
	return finalModel.(*ImportModel), nil
}

func NewImportModel(dir string) *ImportModel {
	items := []list.Item{}

	// Add parent directory option
	if dir != "/" {
		items = append(items, fileItem{label: "..", path: filepath.Dir(dir)})
	}

	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(dir, name)
		if entry.IsDir() || filepath.Ext(name) == ".xlsx" {
			items = append(items, fileItem{label: name, path: fullPath})
		}
	}

	width, height, _ := term.GetSize((os.Stdout.Fd()))
	l := list.New(items, list.NewDefaultDelegate(), width, height-5)
	l.SetSize(width-10, height-10)
	l.Title = fmt.Sprintf("Current directory: %s", dir)

	return &ImportModel{
		list: l,
		dir:  dir,
	}
}

func (m *ImportModel) Init() tea.Cmd {
	return nil
}

func (m *ImportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.list.SelectedItem().(fileItem)
			info, err := os.Stat(selected.path)
			if err != nil {
				m.err = fmt.Errorf("failed to access: %v", err)
				return m, tea.Quit
			}
			if info.IsDir() {
				// Navigate into folder or up
				return NewImportModel(selected.path), nil
			}
			m.FilePath = selected.path
			m.done = true
			return m, tea.Quit
		case "q", "esc":
			m.err = fmt.Errorf("file selection cancelled")
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *ImportModel) View() string {
	if m.done {
		return fmt.Sprintf("Selected file: %s\n", m.FilePath)
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	return m.list.View()
}
