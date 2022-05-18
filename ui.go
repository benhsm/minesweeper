package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const banner = `▙▗▌▗                             
▌▘▌▄ ▛▀▖▞▀▖▞▀▘▌  ▌▞▀▖▞▀▖▛▀▖▞▀▖▙▀▖
▌ ▌▐ ▌ ▌▛▀ ▝▀▖▐▐▐ ▛▀ ▛▀ ▙▄▘▛▀ ▌  
▘ ▘▀▘▘ ▘▝▀▘▀▀  ▘▘ ▝▀▘▝▀▘▌  ▝▀▘▘  `

var (
	selected   = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("3")).Bold(true)
	flag       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // Red
	green      = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	blue       = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))  // Blue
	red        = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // Red
	noColor    = lipgloss.NewStyle().Foreground(lipgloss.NoColor{})
	fieldStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
)

// Constants representing characters used to render certain game elements
const (
	mineRune    = "☀"
	flagRune    = ""
	unknownRune = "⛶"
)

const (
	playing = iota
	lost
	won
)

type model struct {
	field  mineField
	cursor point
	height int
	width  int
}

func (m *model) setSize(w, h int) {
	m.width = w
	m.height = h
}

func newModel() model {
	rand.Seed(time.Now().UnixNano())

	// Easy
	height := 9
	width := 9
	mines := 10

	// Normal
	// height := 16
	// width := 16
	// mines := 40

	// Expert
	// height := 30
	// width := 16
	// mines := 99
	mineField := newMineField(height, width, mines)
	mineField.gameState = playing

	return model{
		field:  mineField,
		cursor: point{0, 0},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.field.gameState {
	case lost, won:
		return updateGameOver(msg, m)
	default:
		return updateGameLoop(msg, m)
	}
}

func updateGameLoop(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var tilesRevealed int
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "w":
			if m.cursor.y > 0 {
				m.cursor.y--
			}
		case "down", "j", "s":
			if m.cursor.y < len(m.field.tiles)-1 {
				m.cursor.y++
			}
		case "left", "h", "a":
			if m.cursor.x > 0 {
				m.cursor.x--
			}
		case "right", "l", "d":
			if m.cursor.x < len(m.field.tiles[0])-1 {
				m.cursor.x++
			}
		case " ":
			if tilesRevealed = m.field.revealTile(m.cursor.x, m.cursor.y); tilesRevealed == -1 {
				// Player activated a mine and lost
				m.field.gameState = lost
			}
			m.field.tilesRemaining -= tilesRevealed
			if m.field.tilesRemaining == 0 {
				m.field.gameState = won
			}
		case "f":
			m.field.flagTile(m.cursor.x, m.cursor.y)
		}
	}
	return m, nil
}

func updateGameOver(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			height := m.height
			width := m.width
			m = newModel()
			m.width = width
			m.height = height
		}

	}
	return m, nil
}

func (m model) View() string {
	//	s := "Minesweeper! \n"
	s := banner
	controls := "\nControls:\n"
	controls += "- arrow keys, 'wasd' or 'hjkl' to move cursor\n"
	controls += "- spacebar to reveal, 'f' to flag\n"
	controls += "- 'q' to quit\n"
	var field string
	for y, row := range m.field.tiles {
		for x, col := range row {
			c := ""
			switch col.state {
			case hidden:
				c = fmt.Sprintf(" %s ", unknownRune)
			case revealed:
				c = strconv.Itoa(col.val)
				style := noColor
				switch col.val {
				case mine:
					c = mineRune
				case 0:
				case 1:
					style = blue
				case 2:
					style = green
				case 3:
					style = red
				default:
					style = red
				}
				c = fmt.Sprintf(" %s ", c)
				c = style.Render(c)
			case flagged:
				c = fmt.Sprintf(" %s ", flagRune)
				c = flag.Render(c)
			}
			if x == m.cursor.x && y == m.cursor.y {
				c = selected.Render(c)
			}
			field += c
		}

		if y != len(m.field.tiles)-1 {
			field += "\n"
		}
	}
	field = fieldStyle.Render(field)
	s = lipgloss.JoinVertical(lipgloss.Center, s, field)
	s += fmt.Sprintf("\n\nUnmined tiles remaining: %d\n", m.field.tilesRemaining)
	if m.field.gameState == won {
		s += "You won! Play again?\n('r' to retry, 'q' to quit)"
	} else if m.field.gameState == lost {
		s += "You lost. Play again?\n('r' to retry, 'q' to quit)"
	}
	s = lipgloss.JoinHorizontal(lipgloss.Center, s, controls)
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		lipgloss.PlaceVertical(m.height, lipgloss.Center, s))
	return s
}
