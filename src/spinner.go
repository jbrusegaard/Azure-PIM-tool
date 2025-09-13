package src

import (
	"fmt"
	"os"
	"time"

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
	Quitting   bool
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
		m.quitting = true
		time.Sleep(100 * time.Millisecond) // allow time for terminal to reset
		return m, nil

	case UpdateMessageMsg: // handle message update
		if msg.Quitting {
			m.quitting = true
			time.Sleep(100 * time.Millisecond) // allow time for terminal to reset
			return m, tea.Quit
		}
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
		return fmt.Sprintf("\n ❌ %s: %s\n\n", m.message, m.err.Error())
	}
	if m.quitting {
		return "\n\n\n"
	}
	str := fmt.Sprintf("\n%s %s\n", m.spinner.View(), m.message)
	return str
}
