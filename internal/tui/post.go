package tui

import (
	"fmt"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
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
	Viewport  viewport.Model
	ready     bool
	spinner   spinner.Model
	width     int
	height    int
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

func (m postModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if !m.ready {
		m.Viewport = viewport.New(m.width, m.height)
		md, _ := glamour.Render(m.Content, "dark")
		m.Viewport.SetContent(md)
		m.Viewport.Height = m.height - 3
		m.Viewport.YPosition = 3
		m.ready = true
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
	return m.Markdown()
}
