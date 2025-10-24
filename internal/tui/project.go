package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	stepName step = iota
	stepProjectDescription
	stepVersion
	stepWebsiteURL
	stepGitHubURL
	stepProjectSummary
)

type CreateProjectModel struct {
	step        step
	answers     map[string]string
	name        textinput.Model
	description textinput.Model
	version     textinput.Model
	websiteURL  textinput.Model
	githubURL   textinput.Model
}

func NewCreateProjectModel() *CreateProjectModel {
	ti := func() textinput.Model {
		t := textinput.New()
		t.Placeholder = ""
		t.CharLimit = 256
		return t
	}

	return &CreateProjectModel{
		step:        stepName,
		answers:     make(map[string]string),
		name:        ti(),
		description: ti(),
		version:     ti(),
		websiteURL:  ti(),
		githubURL:   ti(),
	}
}

func (m *CreateProjectModel) Init() tea.Cmd {
	m.name.Focus()
	return textinput.Blink
}

func (m *CreateProjectModel) focusCurrentInput() {
	switch m.step {
	case stepProjectDescription:
		m.description.Focus()
	case stepVersion:
		m.version.Focus()
	case stepWebsiteURL:
		m.websiteURL.Focus()
	case stepGitHubURL:
		m.githubURL.Focus()
	}
}

func (m *CreateProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepName:
		var cmd tea.Cmd
		m.name, cmd = m.name.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Name"] = m.name.Value()
				m.step = stepProjectDescription
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepProjectDescription:
		var cmd tea.Cmd
		m.description, cmd = m.description.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Description"] = m.description.Value()
				m.step = stepVersion
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepName
				m.name.Focus()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepVersion:
		var cmd tea.Cmd
		m.version, cmd = m.version.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Version"] = m.version.Value()
				m.step = stepWebsiteURL
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepProjectDescription
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepWebsiteURL:
		var cmd tea.Cmd
		m.websiteURL, cmd = m.websiteURL.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Website URL"] = m.websiteURL.Value()
				m.step = stepGitHubURL
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepVersion
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepGitHubURL:
		var cmd tea.Cmd
		m.githubURL, cmd = m.githubURL.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["GitHub URL"] = m.githubURL.Value()
				m.step = stepProjectSummary
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepWebsiteURL
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepProjectSummary:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				return m, tea.Quit
			case "left":
				m.step = stepGitHubURL
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
	}

	return m, nil
}

func (m *CreateProjectModel) View() string {
	var b strings.Builder

	switch m.step {
	case stepName:
		b.WriteString("Enter Project Name:\n")
		b.WriteString(m.name.View())

	case stepProjectDescription:
		b.WriteString("Enter Description:\n")
		b.WriteString(m.description.View())

	case stepVersion:
		b.WriteString("Enter Version:\n")
		b.WriteString(m.version.View())

	case stepWebsiteURL:
		b.WriteString("Enter Website URL:\n")
		b.WriteString(m.websiteURL.View())

	case stepGitHubURL:
		b.WriteString("Enter GitHub URL (optional):\n")
		b.WriteString(m.githubURL.View())

	case stepProjectSummary:
		b.WriteString("\nSummary:\n")
		for _, key := range []string{
			"Name", "Description", "Version", "Website URL", "GitHub URL",
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

func (m *CreateProjectModel) Answers() map[string]string {
	return m.answers
}

func RunCreateProject() (map[string]string, error) {
	final, err := tea.NewProgram(NewCreateProjectModel()).Run()
	if err != nil {
		return nil, err
	}
	return final.(*CreateProjectModel).Answers(), nil
}
