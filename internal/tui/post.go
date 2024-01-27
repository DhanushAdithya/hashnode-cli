package tui

import (
	"fmt"

	"github.com/DhanushAdithya/hashnode-cli/internal/fetch"
	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type postModel struct {
	Id        string
	Heading   string
	Brief     string
	Published string
	Author    string
	URL       string
	Content   string
	ReadTime  int
	Viewport  viewport.Model
	ready     bool
	spinner   spinner.Model
	width     int
	height    int
	Help      help.Model
	Keys      keyMap
	Comment   textarea.Model
}

type keyMap struct {
	Back    key.Binding
	Up      key.Binding
	Down    key.Binding
	Open    key.Binding
	Help    key.Binding
	Quit    key.Binding
	Like    key.Binding
	Comment key.Binding
}

type backMsg struct{}

var keys = keyMap{
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc/bs", "back"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open in browser"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	Like: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "like"),
	),
	Comment: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "comment"),
	),
}

func (i postModel) Title() string       { return i.Heading }
func (i postModel) Description() string { return i.Brief }
func (i postModel) FilterValue() string { return i.Heading }

func (m postModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m postModel) Markdown() string {
	return fmt.Sprintf("%s\n%s", utils.RenderTitle(m.Heading, m.width), m.Viewport.View())
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Back, k.Quit, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Open, k.Like, k.Comment},
		{k.Help, k.Back, k.Quit},
	}
}

func (m postModel) StatusMsg() string {
	return utils.RenderStatusBar(
		m.width,
		m.Author,
		m.ReadTime,
		m.Published,
		m.Viewport.ScrollPercent(),
	)
}

func (m postModel) RenderComment() string {
	return ""
}

func (m postModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if !m.ready {
		m.Viewport = viewport.New(m.width, m.height)
		md, _ := glamour.Render(m.Content, "dark")
		m.Viewport.SetContent(md)
		m.Viewport.Height = m.height - 5
		m.Viewport.YPosition = 3
		m.ready = true
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Viewport.Width = m.width
		if m.Help.ShowAll {
			m.Viewport.Height = m.height - 8
		} else {
			m.Viewport.Height = m.height - 5
		}
		m.Viewport.YPosition = 3
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.Comment.Focused() {
				m.Comment.Blur()
				return m, nil
			}
			return m, func() tea.Msg {
				return backMsg{}
			}
		case key.Matches(msg, keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
			if m.Help.ShowAll {
				m.Keys.Help.SetHelp("?", "close help")
				m.Viewport.Height = m.height - 10
			} else {
				m.Keys.Help.SetHelp("?", "more")
				m.Viewport.Height = m.height - 7
			}
		case key.Matches(msg, keys.Open):
			utils.OpenBrowser(m.URL)
		case key.Matches(msg, keys.Like):
			fetch.LikeResponse(m.Id)
		case key.Matches(msg, keys.Comment):
			// m.Comment.Focus()
			fetch.CommentResponse(m.Id, "test")
			return m, nil
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m postModel) View() string {
	if !m.ready {
		return fmt.Sprintf("%s Loading ...", m.spinner.View())
	}
	help := m.Help.View(m.Keys)
	help = "\n" + lipgloss.NewStyle().MarginLeft(2).Render(help)
	return lipgloss.JoinVertical(lipgloss.Left, m.Markdown(), m.StatusMsg(), help)
}
