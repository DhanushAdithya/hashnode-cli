package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type feedModel struct {
	FeedType    string
	Tags        []string
	MinRead     int
	MaxRead     int
	Page        int
	EndCursor   string
	HasNext     bool
	Posts       list.Model
	FocusedPost postModel
	Help        help.Model
}

type updateSelection struct {
	post postModel
}

func (m feedModel) Init() tea.Cmd {
	return nil
}

func (m feedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Posts.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, nil
		case "enter":
			var ok bool
			selectedItem := m.Posts.SelectedItem()
			if m.FocusedPost, ok = selectedItem.(postModel); ok {
				m.FocusedPost.spinner = spinner.New()
				m.FocusedPost.spinner.Spinner = spinner.Dot
				m.FocusedPost.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
				m.FocusedPost.ready = false
				m.FocusedPost.width = m.Posts.Width()
				m.FocusedPost.height = m.Posts.Height()
				m.FocusedPost.Keys = keys
				m.FocusedPost.Help = help.New()
				return m, func() tea.Msg {
					return updateSelection{
						post: m.FocusedPost,
					}
				}
			}
		}
	}

	m.Posts, cmd = m.Posts.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m feedModel) View() string {
	return m.Posts.View()
}
