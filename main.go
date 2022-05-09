package main

import (
	"fmt"
	"math/rand"
	"os"
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
	ok         = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	fieldStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
)

const (
	hidden = iota
	revealed
	flagged
)

func initialModel() model {
	rand.Seed(time.Now().UnixNano())
	mineField := newMineField(newField(9, 9, 10))
	return model{
		field:          mineField,
		cursor:         point{0, 0},
		gameState:      playing,
		tilesRemaining: 71,
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

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if m.gameState == lost || m.gameState == won {
		return updateGameOver(msg, m)
	}
	return updateGameLoop(msg, m)
}

func updateGameLoop(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var tilesRevealed int
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor.y > 0 {
				m.cursor.y--
			}
		case "down", "j":
			if m.cursor.y < len(m.field)-1 {
				m.cursor.y++
			}
		case "left", "h":
			if m.cursor.x > 0 {
				m.cursor.x--
			}
		case "right", "l":
			if m.cursor.x < len(m.field[0])-1 {
				m.cursor.x++
			}
		case "enter":
			if tilesRevealed = m.field.revealTile(m.cursor.x, m.cursor.y); tilesRevealed == -1 {
				// Player activated a mine and lost
				m.gameState = lost
			}
			m.tilesRemaining -= tilesRevealed
			if m.tilesRemaining == 0 {
				m.gameState = won
			}
		case " ":
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
			m = initialModel()
		}

	}
	return m, nil
}

func (m model) View() string {
	//	s := "Minesweeper! \n"
	s := banner
	s += "\nControls:\n"
	s += "- arrow keys, 'wasd' or 'hjkl' to move cursor\n"
	s += "- enter to reveal, space to flag\n"
	s += "- q to quit\n"
	var field string
	for y, row := range m.field {
		for x, col := range row {
			c := ""
			switch col.state {
			case hidden:
				c = fmt.Sprintf(" %s ", unknownRune)
			case revealed:
				c = fmt.Sprintf(" %s ", col.char)
				c = ok.Render(c)
			case flagged:
				c = fmt.Sprintf(" %s ", flagRune)
				c = flag.Render(c)
			}
			if x == m.cursor.x && y == m.cursor.y {
				c = selected.Render(c)
			}
			field += c
		}

		if y != len(m.field)-1 {
			field += "\n"
		}
	}
	s += fieldStyle.Render(field)
	s += fmt.Sprintf("\n\nUnmined tiles remaining: %d\n", m.tilesRemaining)
	if m.gameState == won {
		s += "You won! Play again? (r to retry, q to quit)"
	} else if m.gameState == lost {
		s += "You exploded. Play again? (r to retry, q to quit)"
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
