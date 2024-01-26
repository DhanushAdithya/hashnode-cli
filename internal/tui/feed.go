package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
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

type loadMore struct{}

func (m feedModel) Init() tea.Cmd {
	return nil
}

func (m feedModel) UpdateSelectedPost() (tea.Model, tea.Cmd) {
	var ok bool
	selectedItem := m.Posts.SelectedItem()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	t := textarea.New()
	t.Placeholder = "Write a thoughtful comment"
	t.ShowLineNumbers = false

	if m.FocusedPost, ok = selectedItem.(postModel); ok {
		m.FocusedPost.spinner = s
		m.FocusedPost.ready = false
		m.FocusedPost.width = m.Posts.Width()
		m.FocusedPost.height = m.Posts.Height()
		m.FocusedPost.Comment = t
		m.FocusedPost.Keys = keys
		m.FocusedPost.Help = help.New()
		return m, func() tea.Msg {
			return updateSelection{post: m.FocusedPost}
		}
	}
	return m, nil
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
			return m.UpdateSelectedPost()
		case "j", "down":
			if m.Posts.Index() == len(m.Posts.Items())-1 {
				items := m.Posts.Items()
				r := make(chan struct{})
				if newItems := fetchPosts(r); len(newItems) > 0 {
					var it []list.Item
					for _, post := range newItems {
						it = append(it, post)
					}
					items = append(items, it...)
					m.Posts.SetItems(items)
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
