package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var bolded = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("3")).Bold(true)
var flag = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
var ok = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))  // Green

const (
	hidden = iota
	revealed
	flagged
)

func initialModel() model {
	rand.Seed(time.Now().UnixNano())
	mineField := newMineField(newField(9, 9, 10))
	return model{
		field:  mineField,
		cursor: point{0, 0},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor.y > 0 {
				m.cursor.y--
			}
		case "down", "j":
			if m.cursor.y < len(m.field) {
				m.cursor.y++
			}
		case "left", "h":
			if m.cursor.x > 0 {
				m.cursor.x--
			}
		case "right", "l":
			if m.cursor.x < len(m.field[0]) {
				m.cursor.x++
			}
		case "enter":
			if m.field.revealTile(m.cursor.x, m.cursor.y) {
				// Player activated a mine and lost
				return m, tea.Quit
			}
		case " ":
			m.field.flagTile(m.cursor.x, m.cursor.y)
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Minesweeper! \n"
	s += "Use HJKL to move cursor, enter to reveal, space to flag\n"
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
				c = bolded.Render(c)
			}
			s += c

		}
		s += "\n"
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
