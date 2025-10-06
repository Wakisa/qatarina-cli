package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type step int

const (
	stepTitle step = iota
	stepKind
	stepProjectID
	stepDescription
	stepCode
	stepFeature
	stepIsDraft
	stepTags
	stepSummary
)

var kindOptions = []string{
	"general", "adhoc", "triangle", "integration", "user_acceptance",
	"regression", "security", "user_interface", "scenario",
}

type listItem struct{ value string }

func (i listItem) Title() string       { return fmt.Sprintf("  %s", i.value) }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.value }

type CreateModel struct {
	step        step
	answers     map[string]string
	title       textinput.Model
	projectID   textinput.Model
	description textinput.Model
	code        textinput.Model
	feature     textinput.Model
	isDraft     textinput.Model
	tags        textinput.Model
	kindList    list.Model
}

func NewCreateModel() *CreateModel {
	ti := func() textinput.Model {
		t := textinput.New()
		t.Placeholder = ""
		t.CharLimit = 256
		return t
	}

	items := make([]list.Item, len(kindOptions))
	for i, k := range kindOptions {
		items[i] = listItem{k}
	}

	kindList := list.New(items, list.NewDefaultDelegate(), 30, 9)
	kindList.Title = "Choose a test case kind"
	kindList.SetShowHelp(false)
	kindList.SetFilteringEnabled(false)
	kindList.DisableQuitKeybindings()

	return &CreateModel{
		step:        stepTitle,
		answers:     make(map[string]string),
		title:       ti(),
		projectID:   ti(),
		description: ti(),
		code:        ti(),
		feature:     ti(),
		isDraft:     ti(),
		tags:        ti(),
		kindList:    kindList,
	}
}

func (m *CreateModel) Init() tea.Cmd {
	m.title.Focus()
	return textinput.Blink
}

func (m *CreateModel) focusCurrentInput() {
	switch m.step {
	case stepProjectID:
		m.projectID.Focus()
	case stepDescription:
		m.description.Focus()
	case stepCode:
		m.code.Focus()
	case stepFeature:
		m.feature.Focus()
	case stepIsDraft:
		m.isDraft.Focus()
	case stepTags:
		m.tags.Focus()
	}
}

func (m *CreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepTitle:
		var cmd tea.Cmd
		m.title, cmd = m.title.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Title"] = m.title.Value()
				m.step = stepKind
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepKind:
		var cmd tea.Cmd
		m.kindList, cmd = m.kindList.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Kind"] = m.kindList.SelectedItem().(listItem).value
				m.step = stepProjectID
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepTitle
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepProjectID:
		var cmd tea.Cmd
		m.projectID, cmd = m.projectID.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				val := m.projectID.Value()
				if _, err := strconv.ParseInt(val, 10, 64); err == nil && val != "0" {
					m.answers["Project ID"] = val
					m.step = stepDescription
					m.focusCurrentInput()
					return m, textinput.Blink
				} else {
					m.projectID.SetValue("")
				}
			} else if key.Type == tea.KeyLeft {
				m.step = stepKind
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepDescription:
		var cmd tea.Cmd
		m.description, cmd = m.description.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Description"] = m.description.Value()
				m.step = stepCode
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepProjectID
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepCode:
		var cmd tea.Cmd
		m.code, cmd = m.code.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Code"] = m.code.Value()
				m.step = stepFeature
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepDescription
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepFeature:
		var cmd tea.Cmd
		m.feature, cmd = m.feature.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Feature/Module"] = m.feature.Value()
				m.step = stepIsDraft
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepCode
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepIsDraft:
		var cmd tea.Cmd
		m.isDraft, cmd = m.isDraft.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				val := m.isDraft.Value()
				if val == "true" || val == "false" {
					m.answers["Is Draft"] = val
					m.step = stepTags
					m.focusCurrentInput()
					return m, textinput.Blink
				} else {
					m.isDraft.SetValue("")
				}
			} else if key.Type == tea.KeyLeft {
				m.step = stepFeature
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepTags:
		var cmd tea.Cmd
		m.tags, cmd = m.tags.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Tags"] = m.tags.Value()
				m.step = stepSummary
				return m, nil
			} else if key.Type == tea.KeyLeft {
				m.step = stepIsDraft
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepSummary:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				return m, tea.Quit
			case "left":
				m.step = stepTags
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
	}

	return m, nil
}

func (m *CreateModel) View() string {
	var b strings.Builder

	switch m.step {
	case stepTitle:
		b.WriteString("Enter Title:\n")
		b.WriteString(m.title.View())

	case stepKind:
		b.WriteString("Select Kind (↑/↓ to navigate, Enter to choose):\n")
		b.WriteString(m.kindList.View())

	case stepProjectID:
		b.WriteString("Enter Project ID:\n")
		b.WriteString(m.projectID.View())

	case stepDescription:
		b.WriteString("Enter Description:\n")
		b.WriteString(m.description.View())

	case stepCode:
		b.WriteString("Enter Code:\n")
		b.WriteString(m.code.View())

	case stepFeature:
		b.WriteString("Enter Feature/Module:\n")
		b.WriteString(m.feature.View())

	case stepIsDraft:
		b.WriteString("Is Draft (true|false):\n")
		b.WriteString(m.isDraft.View())

	case stepTags:
		b.WriteString("Enter Tags (comma-separated):\n")
		b.WriteString(m.tags.View())

	case stepSummary:
		b.WriteString("\nSummary:\n")
		for _, key := range []string{
			"Title", "Kind", "Project ID", "Description",
			"Code", "Feature/Module", "Is Draft", "Tags",
		} {
			val := strings.TrimSpace(m.answers[key])
			if val == "" {
				val = "[missing]"
			}
			b.WriteString(fmt.Sprintf("• %s: %s\n", key, val))
		}
		b.WriteString("\nPress Enter to submit or ← to go back.")
	}

	return b.String()
}

func (m *CreateModel) Answers() []string {
	return []string{
		m.answers["Title"],
		m.answers["Kind"],
		m.answers["Project ID"],
		m.answers["Description"],
		m.answers["Code"],
		m.answers["Feature/Module"],
		m.answers["Is Draft"],
		m.answers["Tags"],
	}
}

func RunCreateTestCase() ([]string, error) {
	final, err := tea.NewProgram(NewCreateModel()).Run()
	if err != nil {
		return nil, err
	}
	return final.(*CreateModel).Answers(), nil
}
