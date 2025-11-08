package ui

import (
	"strings"

	"github.com/chann44/goTele/internals"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	lines      []string
	offset     int // scroll offset (which line is at top)
	viewport   int // height of viewport
	autoScroll bool
}

func InitialModel() model {
	// Sample teleprompter text
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

	lines := strings.Split(text, "\n")

	return model{
		lines:      lines,
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
