package ui

import (
	"strings"
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
	app_state    appState
	filepicker   filepicker.Model
	lines        []string
	offset       int // scroll offset (which line is at top)
	viewport     int // height of viewport
	autoScroll   bool
	quitting     bool
	err          error
	cursor       int
	selectedType inputType
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
	switch m.app_state {
	case appSelectedSource:
		return m.updateInputSelection(msg)
	case appRunning:
		return m.UpdateTelePrompter(msg)
	}
	return m, nil
}

func (m model) updateInputSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			}
		case "enter":
			m.selectedType = inputType(m.cursor)
			m.app_state = appRunning

			// Initialize with sample text for now
			if m.selectedType == inputTypeText {
				text := `Welcome to the presentation.
Today we will discuss the future of technology.
Artificial intelligence is transforming our world.
Machine learning enables computers to learn from data.
Deep learning uses neural networks with many layers.
Natural language processing helps computers understand text.
Computer vision allows machines to interpret images.
Robotics combines AI with physical automation.
The Internet of Things connects billions of devices.
Cloud computing provides scalable infrastructure.
Blockchain enables decentralized applications.
Quantum computing promises exponential speedups.
Biotechnology merges biology with technology.
Renewable energy powers a sustainable future.
Space exploration opens new frontiers.
Virtual reality creates immersive experiences.
Augmented reality overlays digital information.
5G networks enable faster connectivity.
Edge computing brings processing closer to data.
Thank you for your attention.`
				m.lines = strings.Split(text, "\n")
			}

			return m, internals.Tick()
		}
	case tea.WindowSizeMsg:
		m.viewport = msg.Height - 4
	}
	return m, nil
}

func (m model) UpdateTelePrompter(msg tea.Msg) (tea.Model, tea.Cmd) {
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
