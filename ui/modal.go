package ui

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/chann44/goTele/internals"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"

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
	addSource
	appRunning
)

type model struct {
	selectedFile string
	textinput    textinput.Model
	app_state    appState
	filepicker   filepicker.Model
	lines        []string
	offset       int // scroll offset (which line is at top)
	viewport     int // height of viewport
	width        int // terminal width
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
	case addSource:
		return m.updateSourceSelection(msg)
	}

	return m, nil
}

func (m model) updateSourceSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle quit keys first
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.viewport = msg.Height - 4
		m.width = msg.Width
	}

	// Handle filepicker updates first (before intercepting keys)
	if m.selectedType == inputTypeFile {
		m.filepicker, cmd = m.filepicker.Update(msg)

		// Did the user select a file? (This only returns true for files, not directories)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			// Automatically load the file when selected
			fileContent := internals.ReadFile(m.selectedFile)
			wrappedLines := wrapText(string(fileContent), m.width)
			m.lines = append(m.lines, wrappedLines...)
			m.app_state = appRunning
			return m, internals.Tick()
		}

		// Did the user select a disabled file?
		if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
			m.err = errors.New(path + " is not valid.")
			m.selectedFile = ""
			return m, tea.Batch(cmd, ClearErrorAfter(2*time.Second))
		}

		// Let filepicker handle all other messages (including Enter for directory navigation)
		// The filepicker will automatically navigate into directories when Enter is pressed
		return m, cmd
	}

	// Handle text input mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			text := m.textinput.Value()
			wrappedLines := wrapText(text, m.width)
			m.lines = append(m.lines, wrappedLines...)
			m.app_state = appRunning
			return m, internals.Tick()
		}
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) updateInputSelection(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
			switch m.selectedType {
			case inputTypeText:
				ti := textinput.New()
				ti.Placeholder = "Enter your text..."
				ti.Focus()
				ti.CharLimit = 0 // 0 means no limit
				if m.width > 0 {
					ti.Width = m.width - 4 // Leave some margin
				} else {
					ti.Width = 50
				}
				m.textinput = ti
				cmd = textinput.Blink
			case inputTypeFile:
				fp := filepicker.New()
				fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
				fp.CurrentDirectory, _ = os.UserHomeDir()
				fp.DirAllowed = true  // Allow selecting directories for navigation
				fp.FileAllowed = true // Allow selecting files
				if m.viewport > 0 {
					fp.Height = m.viewport - 6 // Leave space for header/footer
				} else {
					fp.Height = 10 // Default height
				}
				m.filepicker = fp
				cmd = m.filepicker.Init()

			}

			m.app_state = addSource
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.viewport = msg.Height - 4
		m.width = msg.Width
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
		m.width = msg.Width

	case internals.TickMsg:
		if m.autoScroll && m.offset < len(m.lines)-1 {
			m.offset++
		}
		return m, internals.Tick()
	}

	return m, nil
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		width = 80
	}

	effectiveWidth := width - 4

	if effectiveWidth <= 0 {
		effectiveWidth = 40 // Minimum reasonable width
	}

	var lines []string
	words := strings.Fields(text)

	if len(words) == 0 {
		return []string{""}
	}

	currentLine := words[0]

	if len(currentLine) > effectiveWidth {
		for len(currentLine) > effectiveWidth {
			lines = append(lines, currentLine[:effectiveWidth])
			currentLine = currentLine[effectiveWidth:]
		}
	}

	for i := 1; i < len(words); i++ {
		word := words[i]

		if len(word) > effectiveWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}
			for len(word) > effectiveWidth {
				lines = append(lines, word[:effectiveWidth])
				word = word[effectiveWidth:]
			}
			currentLine = word
			continue
		}

		testLine := currentLine + " " + word
		if len(testLine) <= effectiveWidth {
			currentLine = testLine
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
