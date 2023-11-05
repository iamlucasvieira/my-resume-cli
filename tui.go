package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// model is the Bubble Tea model
type model struct {
	list     list.Model
	chosen   bool
	choice   int
	subLists map[int]list.Model
}

// Set the default style for the list
var docStyle = lipgloss.NewStyle().Margin(1, 2)

// item is a list item
type item struct {
	title, desc string
}

// Title returns the title of the item
func (i item) Title() string { return i.title }

// Description returns the description of the item
func (i item) Description() string { return i.desc }

// FilterValue returns the value to filter
func (i item) FilterValue() string { return i.title }

// initialModel returns the initial model for the Bubble Tea program
func initialModel() model {

	// Get data from api
	a := api{}

	cv, err := a.requestAll()
	if err != nil {
		fmt.Printf("Error getting data from API: %s\n", err)
	}

	items := []list.Item{
		item{"Info", ""},
		item{"Education", ""},
		item{"Experience", ""},
	}

	m := model{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		subLists: make(map[int]list.Model),
	}

	m.subLists[0] = list.New(infoToItems(cv.Info), list.NewDefaultDelegate(), 0, 0)
	m.subLists[1] = list.New(experienceToItems(cv.Education), list.NewDefaultDelegate(), 0, 0)
	m.subLists[2] = list.New(experienceToItems(cv.Experience), list.NewDefaultDelegate(), 0, 0)

	// Set the title of each list as the item title
	for i, v := range items {
		subList := m.subLists[i] // Get the value from the map
		if item, ok := v.(item); ok {
			subList.Title = item.Title()

			keysExtra := func() []key.Binding {
				return []key.Binding{
					key.NewBinding(
						key.WithKeys("backspace"),
						key.WithHelp("backspace", "back"),
					),
				}
			}
			subList.AdditionalShortHelpKeys = keysExtra
			subList.AdditionalFullHelpKeys = keysExtra
			m.subLists[i] = subList
		}
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	if m.chosen {
		return updateChoice(msg, m)
	}
	return updateMenu(msg, m)
}

func updateMenu(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.chosen = true
			m.choice = m.list.Index()
			newList := m.subLists[m.choice]
			h := m.list.Width()
			v := m.list.Height()
			newList.SetSize(h, v)
			m.subLists[m.choice] = newList
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateChoice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "backspace":
			m.chosen = false
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		newList := m.subLists[m.choice]
		newList.SetSize(msg.Width-h, msg.Height-v)
		m.subLists[m.choice] = newList
	}
	newList, newCmd := m.subLists[m.choice].Update(msg)
	m.subLists[m.choice] = newList
	cmd = newCmd
	return m, cmd
}

func (m model) View() string {
	if m.chosen {
		return docStyle.Render(m.subLists[m.choice].View())
	} else {
		return docStyle.Render(m.list.View())
	}
}

func infoToItems(info []Info) []list.Item {
	items := make([]list.Item, len(info))
	for i, v := range info {
		items[i] = item{v.Name, v.Value}
	}
	return items
}

func experienceToItems(experience []Experience) []list.Item {
	items := make([]list.Item, len(experience))
	for i, v := range experience {
		// Title includes name, institution and dates
		title := fmt.Sprintf("%s - %s (%s - %s)", v.Name, v.Institution, v.Start, v.End)
		// Description is a bullet list of descriptions
		desc := ""
		for _, d := range v.Description {
			desc += fmt.Sprintf("• %s\n", d)
		}
		items[i] = item{title, desc}
	}
	return items
}
