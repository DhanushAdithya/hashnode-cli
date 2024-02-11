package tui

import (
	"strings"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
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

func (m feedModel) InitPosts() ([]list.Item, string) {
	r := make(chan struct{})
	if newItems, cursor := fetchPosts(r, m.FeedType, m.MinRead, m.MaxRead, m.EndCursor); len(newItems) > 0 {
		var it []list.Item
		for _, post := range newItems {
			it = append(it, post)
		}
		return it, cursor
	}
	return []list.Item{}, ""
}

func (m feedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Posts.SetSize(msg.Width, msg.Height-2)
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			idx := utils.FindIndex(fetch.FeedTypes, m.FeedType)
			if idx == len(fetch.FeedTypes)-1 {
				m.FeedType = fetch.FeedTypes[0]
				m.EndCursor = ""
				items, cursor := m.InitPosts()
				m.Posts.SetItems(items)
				m.EndCursor = cursor
			} else {
				m.FeedType = fetch.FeedTypes[idx+1]
				m.EndCursor = ""
				items, cursor := m.InitPosts()
				m.Posts.SetItems(items)
				m.EndCursor = cursor
			}
		case "shift+tab":
			idx := utils.FindIndex(fetch.FeedTypes, m.FeedType)
			if idx == 0 {
				m.FeedType = fetch.FeedTypes[len(fetch.FeedTypes)-1]
				m.EndCursor = ""
				items, cursor := m.InitPosts()
				m.Posts.SetItems(items)
				m.EndCursor = cursor
			} else {
				m.FeedType = fetch.FeedTypes[idx-1]
				m.EndCursor = ""
				items, cursor := m.InitPosts()
				m.Posts.SetItems(items)
				m.EndCursor = cursor
			}
		case "esc":
			if m.Posts.FilterState() != 1 {
				return m, nil
			}
		case "enter":
			return m.UpdateSelectedPost()
		case "j", "down":
			if m.Posts.Index() == len(m.Posts.Items())-1 {
				items := m.Posts.Items()
				it, cursor := m.InitPosts()
				items = append(items, it...)
				m.Posts.SetItems(items)
				m.EndCursor = cursor
			}
		}
	}

	m.Posts, cmd = m.Posts.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m feedModel) View() string {
	var styledTypes []string
	for _, typ := range fetch.FeedTypes {
		if strings.ToUpper(typ) == strings.ToUpper(m.FeedType) {
			styledTypes = append(styledTypes, utils.ActiveTabStyle.Render(typ))
		} else {
			styledTypes = append(styledTypes, utils.InactiveTabStyle.Render(typ))
		}
	}
	types := lipgloss.JoinHorizontal(lipgloss.Left, styledTypes...)
	return lipgloss.JoinVertical(lipgloss.Top, types, m.Posts.View())
}
