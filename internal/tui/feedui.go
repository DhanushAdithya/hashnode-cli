package tui

import (
	"strings"
	"sync"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

var wg sync.WaitGroup

type (
	focusedModel int
	initSearch   struct{}
)

const (
	feed focusedModel = iota
	post
	search
)

type MainModel struct {
	state  focusedModel
	feed   tea.Model
	post   tea.Model
	search tea.Model
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.search, cmd = m.search.Update(msg)
		cmds = append(cmds, cmd)
		m.post, cmd = m.post.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.state != search {
				return m, tea.Quit
			}
		case "ctrl+c":
			return m, tea.Quit
		case "s":
			if m.state != search {
				m.state = search
				m.search, _ = m.search.Update(initSearch{})
				return m, nil
			}
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
	case search:
		searchuh, newCmd := m.search.Update(msg)
		m.search = searchuh
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
	case search:
		return m.search.View()
	}
	return ""
}

func fetchPosts(
	r chan struct{},
	feedType string,
	minRead int,
	maxRead int,
	after string,
) ([]postModel, string) {
	feedResponse := fetch.FeedResponse(r, feedType, minRead, maxRead, after)
	if len(feedResponse.Errors) > 0 {
		close(r)
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
			ReadTime:  post.Node.ReadTimeInMinutes,
		})
	}
	close(r)
	wg.Wait()
	return posts, feedResponse.Data.Feed.PageInfo.EndCursor
}

func GetFeed() {
	r := make(chan struct{})
	F := feedModel{
		FeedType:    strings.ToUpper(viper.GetString("type")),
		Tags:        viper.GetStringSlice("tags"),
		MinRead:     viper.GetInt("min-read"),
		MaxRead:     viper.GetInt("max-read"),
		FocusedPost: postModel{},
	}
	P := postModel{}
	S := searchModel{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		RenderLoad(r)
	}()
	posts, cursor := fetchPosts(r, F.FeedType, F.MinRead, F.MaxRead, "")
	if len(posts) > 0 {
		F.EndCursor = cursor
		P = posts[0]
		var items []list.Item
		for _, post := range posts {
			items = append(items, post)
		}
		F.Posts = list.New(items, list.NewDefaultDelegate(), 0, 0)
		F.Posts.SetShowTitle(false)
	}
	main := MainModel{
		state:  feed,
		feed:   F,
		post:   P,
		search: S.initSearchModel(),
	}
	p := tea.NewProgram(main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		utils.Exit(err)
	}
}
