package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wakisa/qatarina-cli/internal/client"
	"github.com/wakisa/qatarina-cli/internal/schema"
)

type testCaseItem schema.TestCaseResponse

func (i testCaseItem) DisplayTitle() string {
	return i.Title
}

func (i testCaseItem) DescriptionInfo() string {
	return fmt.Sprintf("Code: %s | Kind: %s", i.Code, i.Kind)
}

func (i testCaseItem) FilterValue() string {
	return i.Title
}

type model struct {
	list     list.Model
	selected map[string]schema.TestCaseAssignment
	project  int64
	plan     int64
	quitting bool
	done     bool
}

func RunAssignUI(projectID, planID int64) error {
	testCases, err := fetchTestCases(projectID)
	if err != nil {
		return err
	}

	items := make([]list.Item, len(testCases))
	for i, tc := range testCases {
		items[i] = testCaseItem(tc)
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select test cases (↑/↓ to navigate, press space to toggle, enter to submit)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	m := model{
		list:     l,
		selected: make(map[string]schema.TestCaseAssignment),
		project:  projectID,
		plan:     planID,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	if finalModel.(model).done {
		return finalModel.(model).submitAssignments()
	}

	return nil

}

func fetchTestCases(projectID int64) ([]schema.TestCaseResponse, error) {
	path := fmt.Sprintf("v1/projects/%d/test-cases", projectID)
	resp, err := client.Default().Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error: %s", string(bodyBytes))
	}

	var wrapper struct {
		TestCases []schema.TestCaseResponse `json:"test_cases"`
	}
	if err := json.Unmarshal(bodyBytes, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.TestCases, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			i := m.list.SelectedItem().(testCaseItem)
			if _, ok := m.selected[i.ID]; ok {
				delete(m.selected, i.ID)
			} else {
				m.selected[i.ID] = schema.TestCaseAssignment{TestCaseID: i.ID}
			}

		case "enter":
			m.done = true
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m model) View() string {
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
		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, prefix, tc.Title))
	}
	b.WriteString("\nSelected:\n")
	for _, a := range m.selected {
		b.WriteString(fmt.Sprintf("• %s\n", a.TestCaseID))
	}
	if m.quitting {
		b.WriteString("\nExiting...\n")
	}
	return b.String()
}

func (m model) submitAssignments() error {
	if len(m.selected) == 0 {
		fmt.Println("No test cases selected.")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)
	for id := range m.selected {
		fmt.Printf("Enter user IDs for test case %s (enter one ID and hit enter or multiple IDs comma-separated): ", id)
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
	}

	payload := schema.AssignTestToPlanRequest{
		ProjectID:    m.project,
		PlanID:       m.plan,
		PlannedTests: make([]schema.TestCaseAssignment, 0, len(m.selected)),
	}
	for _, a := range m.selected {
		payload.PlannedTests = append(payload.PlannedTests, a)
	}

	body, _ := json.Marshal(payload)
	path := fmt.Sprintf("v1/test-plans/%d/test-cases", m.plan)
	resp, err := client.Default().Post(path, body)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("API error:", string(bodyBytes))
		return nil
	}

	var message schema.MessageResponse
	if err := json.Unmarshal(bodyBytes, &message); err != nil {
		fmt.Println("Failed to decode response:", err)
		return err
	}
	fmt.Println(message.Message)
	return nil
}
