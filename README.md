# 🎬 goTele

```
  ██████╗  ██████╗ ████████╗███████╗██╗     ███████╗
  ██╔════╝ ██╔═══██╗╚══██╔══╝██╔════╝██║     ██╔════╝
  ██║  ███╗██║   ██║   ██║   █████╗  ██║     █████╗  
  ██║   ██║██║   ██║   ██║   ██╔══╝  ██║     ██╔══╝  
  ╚██████╔╝╚██████╔╝   ██║   ███████╗███████╗███████╗
   ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝╚══════╝╚══════╝
```

> A beautiful terminal-based teleprompter built with Go and Bubble Tea

## ✨ Features

- 🎯 **Auto-scrolling text** - Smooth automatic scrolling with customizable speed
- 🎨 **Beautiful TUI** - Elegant terminal interface with highlighted center line
- ⌨️ **Keyboard controls** - Intuitive navigation and playback controls
- 📜 **Scroll management** - Manual and automatic scrolling modes
- 🎬 **Visual focus** - Center line highlighting for easy reading
- ⚡ **Lightweight** - Fast and efficient Go implementation

## 🚀 Quick Start

### Prerequisites

- Go 1.25.3 or later
- A terminal with ANSI color support

### Installation

```bash
# Clone the repository
git clone https://github.com/chann44/goTele.git
cd goTele

# Install dependencies
go mod download

# Run the application
go run cmd/main.go
```

## 🎮 Usage

### Controls

| Key | Action |
|-----|--------|
| `Space` | Toggle auto-scroll (pause/play) |
| `↑` / `k` | Scroll up one line |
| `↓` / `j` | Scroll down one line |
| `r` | Reset to beginning |
| `q` / `Ctrl+C` | Quit application |

### Features

- **Auto-scroll mode**: Text automatically scrolls at a steady pace
- **Manual mode**: Use arrow keys to navigate at your own pace
- **Center highlighting**: The current line is highlighted in yellow for easy focus
- **Visual feedback**: Status indicator shows whether auto-scroll is playing or paused

## 📁 Project Structure

```
goTele/
├── cmd/
│   └── main.go          # Application entry point
├── internals/
│   └── tick.go          # Timer tick implementation
├── ui/
│   ├── modal.go         # UI model definition
│   └── view.go          # UI rendering logic
├── go.mod               # Go module dependencies
└── README.md           # This file
```

## 🛠️ Development

### Building

```bash
# Build the binary
go build -o goTele cmd/main.go

# Run the binary
./goTele
```

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - Bubble Tea components

## 🎨 Customization

To customize the teleprompter text, edit the `text` variable in `ui/modal.go`:

```go
text := `Your custom text here.
Each line will be displayed separately.
Add as many lines as you need.`
```

To adjust the scroll speed, modify the tick interval in `internals/tick.go`:

```go
return tea.Tick(200*time.Millisecond, ...) // Adjust duration
```

## 📝 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 👤 Author

**chann44**

- GitHub: [@chann44](https://github.com/chann44)

---

Made with ❤️ using Go and Bubble Tea

