package tui

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/benhsm/minesweeper/game"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

type point struct {
	x int
	y int
}

type gameKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Reveal key.Binding
	Flag   key.Binding
	Menu   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k gameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Menu}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k gameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Reveal, k.Flag},              // Second column
		{k.Help, k.Quit, k.Menu},        // Third column
	}
}

var gameKeys = gameKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h", "a"),
		key.WithHelp("←/h/a", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "d"),
		key.WithHelp("→/l/d", "move left"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Reveal: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "reveal"),
	),
	Flag: key.NewBinding(
		key.WithKeys("f", ";"),
		key.WithHelp("f/;", "flag"),
	),
	Menu: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "menu"),
	),
}

type gameOverKeyMap struct {
	Quit  key.Binding
	Retry key.Binding
	Menu  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k gameOverKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Retry, k.Menu}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k gameOverKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Retry, k.Menu},
	}
}

var gameOverKeys = gameOverKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Retry: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "retry"),
	),
	Menu: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "menu"),
	),
}

type gameModel struct {
	endKeys gameOverKeyMap
	keys    gameKeyMap
	help    help.Model
	field   game.MineField
	cursor  point
	inGame  bool
}

func newGameModel(height, width, mines int) gameModel {
	rand.Seed(time.Now().UnixNano())

	mineField := game.NewMineField(height, width, mines)
	mineField.GameState = playing

	return gameModel{
		endKeys: gameOverKeys,
		keys:    gameKeys,
		help:    help.New(),
		field:   mineField,
		cursor:  point{0, 0},
	}
}

func (m gameModel) update(msg tea.Msg) (gameModel, tea.Cmd) {

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.field.GameState {
	case lost, won:
		return updateGameOver(msg, m)
	default:
		return updateGameLoop(msg, m)
	}
}

func updateGameLoop(msg tea.Msg, m gameModel) (gameModel, tea.Cmd) {
	var tilesRevealed int
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor.y > 0 {
				m.cursor.y--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor.y < len(m.field.Tiles)-1 {
				m.cursor.y++
			}
		case key.Matches(msg, m.keys.Left):
			if m.cursor.x > 0 {
				m.cursor.x--
			}
		case key.Matches(msg, m.keys.Right):
			if m.cursor.x < len(m.field.Tiles[0])-1 {
				m.cursor.x++
			}
		case key.Matches(msg, m.keys.Reveal):
			if tilesRevealed = m.field.RevealTile(m.cursor.x, m.cursor.y); tilesRevealed == -1 {
				// Player activated a mine and lost
				m.field.GameState = lost
			}
			m.field.TilesRemaining -= tilesRevealed
			if m.field.TilesRemaining == 0 {
				m.field.GameState = won
			}
		case key.Matches(msg, m.keys.Flag):
			m.field.FlagTile(m.cursor.x, m.cursor.y)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Menu):
			m.inGame = false
		}
	}
	return m, nil
}

func updateGameOver(msg tea.Msg, m gameModel) (gameModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.endKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.endKeys.Retry):
			m = newGameModel(len(m.field.Tiles), len(m.field.Tiles[0]), m.field.Mines)
			m.inGame = true
		case key.Matches(msg, m.endKeys.Menu):
			m.inGame = false
		}

	}
	return m, nil
}

func (m gameModel) view() string {
	//	s := "Minesweeper! \n"
	var s string
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
		s += "You won! Play again?"
	} else if m.field.GameState == lost {
		s += "You lost. Play again?"
	}

	if m.field.GameState == won || m.field.GameState == lost {
		s += "\n" + m.help.View(m.endKeys)
	} else {
		s += "\n" + m.help.View(m.keys)
	}
	return s
}
