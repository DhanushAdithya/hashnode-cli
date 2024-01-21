package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type postModel struct {
	Heading   string
	Brief     string
	Published string
	Author    string
	URL       string
	Content   string
}

func (i postModel) Title() string       { return i.Heading }
func (i postModel) Description() string { return i.Brief }
func (i postModel) FilterValue() string { return i.Heading }

func (m postModel) Init() tea.Cmd {
	return nil
}

func (m postModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m postModel) View() string {
	md, _ := glamour.Render(m.Content, "dark")
	return md
}
