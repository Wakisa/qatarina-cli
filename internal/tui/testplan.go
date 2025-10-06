package tui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wakisa/qatarina-cli/internal/schema"
)

type testCaseItem schema.TestCaseResponse

func (i testCaseItem) DisplayTitle() string { return i.Title }
func (i testCaseItem) DisplayDescription() string {
	return fmt.Sprintf("Code: %s | Kind: %s", i.Code, i.Kind)
}
func (i testCaseItem) FilterValue() string { return i.Title }
func (i testCaseItem) GetID() string       { return i.ID }

type AssignModel struct {
	list     list.Model
	selected map[string]schema.TestCaseAssignment
	project  int64
	plan     int64
	quitting bool
	done     bool
}

func RunAssignUI(projectID, planID int64, testCases []schema.TestCaseResponse) (*AssignModel, error) {
	items := make([]list.Item, len(testCases))
	for i, tc := range testCases {
		items[i] = testCaseItem(tc)
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select test cases (↑/↓ to navigate, space to toggle, enter to submit)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	m := &AssignModel{
		list:     l,
		selected: make(map[string]schema.TestCaseAssignment),
		project:  projectID,
		plan:     planID,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	return finalModel.(*AssignModel), nil
}

func (m *AssignModel) Init() tea.Cmd {
	return nil
}

func (m *AssignModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			m.list.CursorUp()
		case "down", "j":
			m.list.CursorDown()
		case " ":
			tc := m.list.SelectedItem().(testCaseItem)
			if _, ok := m.selected[tc.ID]; ok {
				delete(m.selected, tc.ID)
			} else {
				m.selected[tc.ID] = schema.TestCaseAssignment{TestCaseID: tc.ID}
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *AssignModel) View() string {
	var b strings.Builder
	for i, item := range m.list.Items() {
		tc := item.(testCaseItem)
		prefix := "[ ]"
		if _, ok := m.selected[tc.ID]; ok {
			prefix = "[x]"
		}
		cursor := " "
		if i == m.list.Index() {
			cursor = "=>"
		}
		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, prefix, tc.DisplayTitle()))
	}
	b.WriteString("\nSelected:\n")
	for _, item := range m.list.Items() {
		tc := item.(testCaseItem)
		if _, ok := m.selected[tc.ID]; ok {
			b.WriteString(fmt.Sprintf("• %s\n", tc.DisplayTitle()))
		}
	}
	if m.quitting {
		b.WriteString("\nExiting...\n")
	}
	return b.String()
}

func (m *AssignModel) CollectedAssignments() []schema.TestCaseAssignment {
	assignments := make([]schema.TestCaseAssignment, 0, len(m.selected))
	reader := bufio.NewReader(os.Stdin)

	for id := range m.selected {
		var title string
		for _, item := range m.list.Items() {
			tc := item.(testCaseItem)
			if tc.ID == id {
				title = tc.DisplayTitle()
				break
			}
		}

		fmt.Printf("Enter user IDs for test case \"%s\" (comma-separated): ", title)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Split(input, ",")
		var userIDs []int64
		for _, p := range parts {
			p = strings.TrimSpace(p)
			uid, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				fmt.Printf("Invalid user ID: %s\n", p)
				continue
			}
			userIDs = append(userIDs, uid)
		}
		m.selected[id] = schema.TestCaseAssignment{
			TestCaseID: id,
			UserIDs:    userIDs,
		}
		assignments = append(assignments, m.selected[id])
	}
	return assignments
}
