package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type userStep int

const (
	stepFirstName userStep = iota
	stepLastName
	stepDisplayName
	stepEmail
	stepPassword
	userSummary
)

type UserCreateModel struct {
	step        userStep
	answers     map[string]string
	firstName   textinput.Model
	lastName    textinput.Model
	displayName textinput.Model
	email       textinput.Model
	password    textinput.Model
}

func NewUserCreateModel() *UserCreateModel {
	ti := func() textinput.Model {
		t := textinput.New()
		t.Placeholder = ""
		t.CharLimit = 256
		return t
	}

	return &UserCreateModel{
		step:        stepFirstName,
		answers:     make(map[string]string),
		firstName:   ti(),
		lastName:    ti(),
		displayName: ti(),
		email:       ti(),
		password:    ti(),
	}
}

func (m *UserCreateModel) Init() tea.Cmd {
	m.firstName.Focus()
	return textinput.Blink
}

func (m *UserCreateModel) focusCurrentInput() {
	switch m.step {
	case stepFirstName:
		m.firstName.Focus()
	case stepLastName:
		m.lastName.Focus()
	case stepDisplayName:
		m.displayName.Focus()
	case stepEmail:
		m.email.Focus()
	case stepPassword:
		m.password.Focus()
	}
}

func (m *UserCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepFirstName:
		var cmd tea.Cmd
		m.firstName, cmd = m.firstName.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["FirstName"] = m.firstName.Value()
				m.step = stepLastName
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepLastName:
		var cmd tea.Cmd
		m.lastName, cmd = m.lastName.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["LastName"] = m.lastName.Value()
				m.step = stepDisplayName
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepFirstName
				m.firstName.Focus()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepDisplayName:
		var cmd tea.Cmd
		m.displayName, cmd = m.displayName.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["DisplayName"] = m.displayName.Value()
				m.step = stepEmail
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepLastName
				m.lastName.Focus()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepEmail:
		var cmd tea.Cmd
		m.email, cmd = m.email.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Email"] = m.email.Value()
				m.step = stepPassword
				m.focusCurrentInput()
				return m, textinput.Blink
			} else if key.Type == tea.KeyLeft {
				m.step = stepDisplayName
				m.displayName.Focus()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case stepPassword:
		var cmd tea.Cmd
		m.password, cmd = m.password.Update(msg)
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.Type == tea.KeyEnter {
				m.answers["Password"] = m.password.Value()
				m.step = userSummary
				return m, nil
			} else if key.Type == tea.KeyLeft {
				m.step = stepEmail
				m.email.Focus()
				return m, textinput.Blink
			}
		}
		return m, cmd

	case userSummary:
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "enter":
				return m, tea.Quit
			case "left":
				m.step = stepPassword
				m.focusCurrentInput()
				return m, textinput.Blink
			}
		}
	}
	return m, nil
}

func (m *UserCreateModel) View() string {
	var b strings.Builder

	switch m.step {
	case stepFirstName:
		b.WriteString("Enter First Name:\n")
		b.WriteString(m.firstName.View())

	case stepLastName:
		b.WriteString("Enter Last Name:\n")
		b.WriteString(m.lastName.View())

	case stepDisplayName:
		b.WriteString("Enter Display Name:\n")
		b.WriteString(m.displayName.View())

	case stepEmail:
		b.WriteString("Enter Email:\n")
		b.WriteString(m.email.View())

	case stepPassword:
		b.WriteString("Enter Password:\n")
		b.WriteString(m.password.View())

	case userSummary:
		b.WriteString("\nSummary:\n")
		for _, key := range []string{"FirstName", "LastName", "DisplayName", "Email", "Password"} {
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

func (m *UserCreateModel) Answers() map[string]string {
	return m.answers
}

func RunCreateUserWizard() (map[string]string, error) {
	final, err := tea.NewProgram(NewUserCreateModel()).Run()
	if err != nil {
		return nil, err
	}
	return final.(*UserCreateModel).Answers(), nil
}
