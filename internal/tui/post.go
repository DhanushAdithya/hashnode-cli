package tui

import (
	"fmt"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type postModel struct {
	Heading   string
	Brief     string
	Published string
	Author    string
	URL       string
	Content   string
	Viewport  viewport.Model
	ready     bool
	spinner   spinner.Model
	width     int
	height    int
	Help      help.Model
	Keys      keyMap
}

type keyMap struct {
	Back     key.Binding
	Up       key.Binding
	Down     key.Binding
	Open     key.Binding
	Help     key.Binding
	Quit     key.Binding
	Like     key.Binding
	Comment  key.Binding
	Bookmark key.Binding
}

var keys = keyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Back to Feed"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Scroll down"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "Open post in browser"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Show help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	),
	Like: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Like post"),
	),
	Comment: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Comment on post"),
	),
	Bookmark: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "Bookmark post"),
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
	return []key.Binding{k.Help, k.Back}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Open, k.Like, k.Comment, k.Bookmark},
		{k.Help, k.Back, k.Quit},
	}
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
		m.Viewport.Height = m.height - 5
		m.Viewport.YPosition = 3
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			return m, nil
		case key.Matches(msg, keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
			if m.Help.ShowAll {
				m.Viewport.Height = m.height - 9
			} else {
				m.Viewport.Height = m.height - 5
			}
		case key.Matches(msg, keys.Open):
			utils.OpenBrowser(m.URL)
		case key.Matches(msg, keys.Like):
			return m, nil
		case key.Matches(msg, keys.Comment):
			return m, nil
		case key.Matches(msg, keys.Bookmark):
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
	var help string
	help = m.Help.View(m.Keys)
	return m.Markdown() + "\n\n" + lipgloss.NewStyle().MarginLeft(2).Render(help)
}
