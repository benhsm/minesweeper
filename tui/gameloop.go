package tui

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/benhsm/minesweeper/game"
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

type gameModel struct {
	field  game.MineField
	cursor point
	height int
	width  int
}

const (
	playing = iota
	lost
	won
)

type point struct {
	x int
	y int
}

func newGameModel() gameModel {
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
	mineField := game.NewMineField(height, width, mines)
	mineField.GameState = playing

	return gameModel{
		field:  mineField,
		cursor: point{0, 0},
	}
}

func updateGameLoop(msg tea.Msg, m gameModel) (tea.Model, tea.Cmd) {
	var tilesRevealed int
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "w":
			if m.cursor.y > 0 {
				m.cursor.y--
			}
		case "down", "j", "s":
			if m.cursor.y < len(m.field.Tiles)-1 {
				m.cursor.y++
			}
		case "left", "h", "a":
			if m.cursor.x > 0 {
				m.cursor.x--
			}
		case "right", "l", "d":
			if m.cursor.x < len(m.field.Tiles[0])-1 {
				m.cursor.x++
			}
		case " ":
			if tilesRevealed = m.field.RevealTile(m.cursor.x, m.cursor.y); tilesRevealed == -1 {
				// Player activated a mine and lost
				m.field.GameState = lost
			}
			m.field.TilesRemaining -= tilesRevealed
			if m.field.TilesRemaining == 0 {
				m.field.GameState = won
			}
		case "f":
			m.field.FlagTile(m.cursor.x, m.cursor.y)
		}
	}
	return m, nil
}

func updateGameOver(msg tea.Msg, m gameModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			height := m.height
			width := m.width
			m = newGameModel()
			m.width = width
			m.height = height
		}

	}
	return m, nil
}

func (m gameModel) View() string {
	//	s := "Minesweeper! \n"
	s := banner
	controls := "\nControls:\n"
	controls += "- arrow keys, 'wasd' or 'hjkl' to move cursor\n"
	controls += "- spacebar to reveal, 'f' to flag\n"
	controls += "- 'q' to quit\n"
	var field string
	for y, row := range m.field.Tiles {
		for x, col := range row {
			c := ""
			switch col.State {
			case game.Hidden:
				c = fmt.Sprintf(" %s ", unknownRune)
			case game.Revealed:
				c = strconv.Itoa(col.Val)
				style := noColor
				switch col.Val {
				case game.Mine:
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
			case game.Flagged:
				c = fmt.Sprintf(" %s ", flagRune)
				c = flag.Render(c)
			}
			if x == m.cursor.x && y == m.cursor.y {
				c = selected.Render(c)
			}
			field += c
		}

		if y != len(m.field.Tiles)-1 {
			field += "\n"
		}
	}
	field = fieldStyle.Render(field)
	s = lipgloss.JoinVertical(lipgloss.Center, s, field)
	s += fmt.Sprintf("\n\nUnmined tiles remaining: %d\n", m.field.TilesRemaining)
	if m.field.GameState == won {
		s += "You won! Play again?\n('r' to retry, 'q' to quit)"
	} else if m.field.GameState == lost {
		s += "You lost. Play again?\n('r' to retry, 'q' to quit)"
	}
	s = lipgloss.JoinHorizontal(lipgloss.Center, s, controls)
	s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		lipgloss.PlaceVertical(m.height, lipgloss.Center, s))
	return s
}
