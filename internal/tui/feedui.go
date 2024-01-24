package tui

import (
	"fmt"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type focusedModel int

const (
	feed focusedModel = iota
	post
)

type MainModel struct {
	state focusedModel
	feed  tea.Model
	post  tea.Model
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case updateSelection:
		m.state = post
		m.post = msg.post
	case backMsg:
		m.state = feed
	}
	switch m.state {
	case feed:
		feeduh, newCmd := m.feed.Update(msg)
		m.feed = feeduh
		cmd = newCmd
	case post:
		postuh, newCmd := m.post.Update(msg)
		m.post = postuh
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case feed:
		return m.feed.View()
	case post:
		return m.post.View()
	}
	return ""
}

func fetchPosts() []postModel {
	feedResponse := fetch.FeedResponse()
	if len(feedResponse.Errors) > 0 {
		utils.RenderAPIErrors(feedResponse.Errors)
	}
	var posts []postModel
	for _, post := range feedResponse.Data.Feed.Edges {
		posts = append(posts, postModel{
			Id:        post.Node.Id,
			Heading:   post.Node.Title,
			Brief:     post.Node.Brief,
			Published: post.Node.PublishedAt,
			Author:    post.Node.Author.Name,
			URL:       post.Node.URL,
			Content:   post.Node.Content.Markdown,
		})
	}
	return posts
}

func GetFeed() {
	F := feedModel{
		FeedType:    viper.GetString("type"),
		Tags:        viper.GetStringSlice("tags"),
		MinRead:     viper.GetInt("min-read"),
		MaxRead:     viper.GetInt("max-read"),
		FocusedPost: postModel{},
	}
	P := postModel{}
	posts := fetchPosts()
	if len(posts) > 0 {
		P = posts[0]
		var items []list.Item
		for _, post := range posts {
			items = append(items, post)
		}
		F.Posts = list.New(items, list.NewDefaultDelegate(), 0, 0)
		F.Posts.Title = fmt.Sprintf("Feed (%d posts)", len(posts))
	}
	main := MainModel{
		state: feed,
		feed:  F,
		post:  P,
	}
	p := tea.NewProgram(main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		utils.Exit(err)
	}
}
