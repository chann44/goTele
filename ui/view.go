package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	// Calculate center line of viewport

	var b strings.Builder

	// Header
	status := "PLAYING"
	if !m.autoScroll {
		status = "PAUSED"
	}
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render(fmt.Sprintf("🎬 Teleprompter [%s] | Line %d/%d", status, m.offset+1, len(m.lines)))

	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("─", 60) + "\n")

	// Render visible lines
	for i := 0; i < m.viewport; i++ {
		// Center line index in the viewport
		centerLine := m.viewport / 2

		// Calculate which actual line this viewport position should show
		// If offset is 0 and we're at centerLine, we show lines[0]
		lineIdx := m.offset + (i - centerLine)

		// Render empty lines if we're before the start or after the end
		if lineIdx < 0 || lineIdx >= len(m.lines) {
			b.WriteString("\n")
			continue
		}

		line := m.lines[lineIdx]

		// Calculate distance from center
		distFromCenter := i - centerLine
		if distFromCenter < 0 {
			distFromCenter = -distFromCenter
		}

		var styledLine string

		// Center line (highlighted and large)
		if i == centerLine {
			styledLine = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFF00")).
				Background(lipgloss.Color("#333333")).
				Padding(0, 2).
				Render("▶ " + line + " ◀")

			// Lines near center (medium size, slightly dimmed)
		} else if distFromCenter <= 2 {
			styledLine = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Render("  " + line)

			// Lines further away (smaller, more dimmed)
		} else if distFromCenter <= 5 {
			styledLine = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render("  " + line)

			// Lines far from center (very dim)
		} else {
			styledLine = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#444444")).
				Render("  " + line)
		}

		b.WriteString(styledLine + "\n")
	}

	// Footer
	b.WriteString(strings.Repeat("─", 60) + "\n")
	help := lipgloss.NewStyle().
		Faint(true).
		Render("space: pause/play • ↑/↓: manual scroll • r: reset • q: quit")
	b.WriteString(help)

	return b.String()
}
