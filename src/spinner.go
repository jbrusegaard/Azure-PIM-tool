package src

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func StartSpinner(message string, spinnerType spinner.Spinner) *tea.Program {
	p := tea.NewProgram(initialModel(message, spinnerType))
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	return p
}

type errMsg error

type spinnerModel struct {
	spinner  spinner.Model
	quitting bool
	err      error
	message  string
}

func initialModel(message string, spinnerType spinner.Spinner) spinnerModel {
	s := spinner.New()
	s.Spinner = spinnerType
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#BB15D4"))
	return spinnerModel{spinner: s, message: message}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

type UpdateMessageMsg struct {
	NewMessage string
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	case UpdateMessageMsg: // handle message update
		m.message = msg.NewMessage
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n%s %s\n\n", m.spinner.View(), m.message)
	if m.quitting {
		return "\n\n\n\n\n"
		// return str + "\n"
	}
	return str
}
