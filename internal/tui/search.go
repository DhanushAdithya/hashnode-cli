package tui

import (
	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchModel struct {
	Input       textinput.Model
	Posts       list.Model
	FocusedPost postModel
	Help        help.Model
	Width       int
	Height      int
}

func (s searchModel) initSearchModel() searchModel {
	input := textinput.New()
	input.Placeholder = "Search for a post"
	input.Focus()

	posts := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	posts.SetShowTitle(false)

	return searchModel{
		Input:       input,
		Posts:       posts,
		FocusedPost: postModel{},
	}
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) UpdateSelectedPost() (tea.Model, tea.Cmd) {
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

		postResponse := fetch.PostResponse(m.FocusedPost.Id)
		if len(postResponse.Errors) > 0 {
			utils.RenderAPIErrors(postResponse.Errors)
		}
		m.FocusedPost.URL = postResponse.Data.Post.URL
		m.FocusedPost.Content = postResponse.Data.Post.Content.Markdown
		m.FocusedPost.ReadTime = postResponse.Data.Post.ReadTimeInMinutes

		return m, func() tea.Msg {
			return updateSelection{post: m.FocusedPost}
		}
	}
	return m, nil
}

func (m searchModel) updateSearchPosts() []list.Item {
	results := fetch.SearchResponse(make(chan struct{}), m.Input.Value())
	var items []list.Item
	for _, post := range results.Hits {
		items = append(items, postModel{
			Id:        post.Id,
			Heading:   post.Title,
			Brief:     post.Brief,
			Published: post.PublishedAt,
			Author:    post.Author,
		})
	}
	return items
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Input, cmd = m.Input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Posts.SetSize(msg.Width, msg.Height-4)
	case initSearch:
		m.Input.SetValue("")
	case tea.KeyMsg:
		m.Posts.SetItems(m.updateSearchPosts())
		switch msg.String() {
		case "tab":
			if m.Input.Focused() {
				m.Input.Blur()
			} else {
				m.Input.Focus()
			}
		case "enter":
			return m.UpdateSelectedPost()
		}
	}

	if !m.Input.Focused() {
		m.Posts, cmd = m.Posts.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m searchModel) View() string {
	searchHeader := utils.ActiveTabStyle.MarginBottom(1).Render("Search")
	return lipgloss.JoinVertical(lipgloss.Top, searchHeader, m.Input.View(), m.Posts.View())
}
