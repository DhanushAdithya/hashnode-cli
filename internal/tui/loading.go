package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Loading struct {
	spinner  spinner.Model
	loading  bool
	response chan struct{}
}

func newLoading(response chan struct{}) Loading {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Loading{
		spinner:  s,
		loading:  true,
		response: response,
	}
}

func (l Loading) Init() tea.Cmd {
	return l.spinner.Tick
}

func (l Loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.spinner, _ = l.spinner.Update(msg)
	select {
	case <-l.response:
		l.loading = false
		return l, tea.Quit
	default:
		return l, l.spinner.Tick
	}
}

func (l Loading) View() string {
	if l.loading {
		return fmt.Sprintf(" %s Loading ... ", l.spinner.View())
	} else {
		return ""
	}
}

func RenderLoad(response chan struct{}) {
	p := tea.NewProgram(newLoading(response))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
