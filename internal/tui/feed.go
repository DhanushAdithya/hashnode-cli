package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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
}

type updateSelection struct {
	post postModel
}

func (m feedModel) Init() tea.Cmd {
	return nil
}

func (m feedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Posts.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.FocusedPost = m.Posts.SelectedItem().(postModel)
			return m, func() tea.Msg {
				return updateSelection{
					post: m.FocusedPost,
				}
			}
		}
	}

	var cmd tea.Cmd
	m.Posts, cmd = m.Posts.Update(msg)
	return m, cmd
}

func (m feedModel) View() string {
	return m.Posts.View()
}
