package ui

import (
	"time"

	"github.com/chann44/goTele/internals"
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type appState int

type clearErrorMsg struct{}

type inputType int

const (
	inputTypeText inputType = iota
	inputTypeFile
	inputTypeUrl
)

const (
	appSelectedSource appState = iota
	appRunning
)

type model struct {
	app_state  appState
	filepicker filepicker.Model
	lines      []string
	offset     int // scroll offset (which line is at top)
	viewport   int // height of viewport
	autoScroll bool
	quitting   bool
	err        error
	cursor     int
}

func ClearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func InitialModel() model {

	return model{
		app_state:  appSelectedSource,
		offset:     0,
		viewport:   20,
		autoScroll: true,
	}
}

func (m model) Init() tea.Cmd {
	return internals.Tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			m.autoScroll = !m.autoScroll
		case "up", "k":
			if m.offset > 0 {
				m.offset--
			}
		case "down", "j":
			if m.offset < len(m.lines)-1 {
				m.offset++
			}
		case "r":
			m.offset = 0
		}

	case tea.WindowSizeMsg:
		m.viewport = msg.Height - 4

	case internals.TickMsg:
		if m.autoScroll && m.offset < len(m.lines)-1 {
			m.offset++
		}
		return m, internals.Tick()
	}

	return m, nil
}
